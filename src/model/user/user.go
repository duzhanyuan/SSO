package user

import (
	"crypto/rc4"
	"encoding/hex"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"service/errormap"
	"service/mongodb"
	"service/monitor"
	"service/myredis"
	"time"

	"util"
)

const (
	tokenLen         = 20
	keyLen           = 16
	serverPrivateKey = "abcdef"
	userTable        = "userTable"
	serviceTable     = "serviceTable"
)

type User struct {
	UserID       string `json:"userid" bson:"_id,omitempty"`
	UserName     string `json:"username" bson:"username"`
	PwdEncrypted string `json:"password" bson:"password"`
}

type Service struct {
	ServiceID   string `json:"serviceid" bson:"_id,omitempty"`
	ServiceName string `json:"servicename" bson:"servicename"`
	ServiceKey  string `json:"servicekey" bson:"servicekey"`
}

func newToken() string {
	token := ""
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < tokenLen; i++ {
		token += string(int(rune('A')) + r.Intn(26))
	}
	return token
}

func Register(userName, pwdEncrypted string) int {
	user := User{}
	exist := mongodb.Exec(userTable, func(c *mgo.Collection) error {
		return c.Find(bson.M{"username": userName}).One(&user)
	})
	if exist {
		return errormap.Exist
	}
	user.UserID = bson.NewObjectId().Hex()
	user.UserName = userName
	user.PwdEncrypted = pwdEncrypted
	mongodb.Exec(userTable, func(c *mgo.Collection) error {
		return c.Insert(user)
	})
	return errormap.Success
}

func RegisterService(name string) (string, int) {
	service := Service{}
	exist := mongodb.Exec(serviceTable, func(c *mgo.Collection) error {
		return c.Find(bson.M{"servicename": name}).One(&service)
	})
	if exist {
		return "", errormap.Exist
	}
	key := util.GenRandomBytes(16)
	keyStr := hex.EncodeToString(key)
	service.ServiceID = bson.NewObjectId().Hex()
	service.ServiceName = name
	service.ServiceKey = keyStr
	mongodb.Exec(serviceTable, func(c *mgo.Collection) error {
		return c.Insert(service)
	})
	client := myredis.ClusterClient(name)
	client.Set(name, keyStr, 0)
	return keyStr, errormap.Success
}

func Login(userName string, timestamp string) (string, string, int) {
	user := User{}
	exist := mongodb.Exec(userTable, func(c *mgo.Collection) error {
		return c.Find(bson.M{"username": userName}).One(&user)
	})
	if !exist {
		return "", "", errormap.NotExist
	}
	pwdBytes, _ := hex.DecodeString(user.PwdEncrypted)
	if util.CheckTimestamp(timestamp, pwdBytes) == false {
		return "", "", errormap.IllegalTS
	}

	key := util.GenRandomBytes(keyLen)
	c2, _ := rc4.NewCipher(pwdBytes)
	encryptedKeys := make([]byte, keyLen)
	c2.XORKeyStream(encryptedKeys, key)
	A := hex.EncodeToString(encryptedKeys)

	keyStr := hex.EncodeToString(key)
	client := myredis.ClusterClient(keyStr)
	client.Set(keyStr, userName, 0)
	bStr := keyStr + ":" + userName
	B := util.EncryptString(bStr, []byte(serverPrivateKey))
	monitor.IncrCount()
	return A, B, errormap.Success
}

func Logout(userID string) int {
	user := User{}
	exist := mongodb.FindByID(userTable, userID, &user)
	if !exist {
		return errormap.NotExist
	}
	//TODO: update redis
	return errormap.Success
}
