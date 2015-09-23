package user

import (
	"crypto/sha1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"gopkg.in/redis.v3"
	//"fmt"
	"io"
	"math/rand"
	"service"
	"service/mongodb"
	"service/myredis"
	"time"
	//"util"
)

const (
	tokenLen  = 20
	userTable = "userTable"
	pwdSalt   = "ZG1sMmFXRnVJR2x6SUdFZ1oyOXZaQ0JuYVhKcw=="
)

type User struct {
	UserID       string `json:"userid" bson:"_id,omitempty"`
	UserName     string `json:"username" bson:"username"`
	PwdEncrypted string `json:"password" bson:"password"`
}

func newToken() string {
	token := ""
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < tokenLen; i++ {
		token += string(int(rune('A')) + r.Intn(26))
	}
	return token
}

func passwordEncrypt(password string) string {
	pass := password + pwdSalt
	sh := sha1.New()
	_, err := io.WriteString(sh, pass)
	if err != nil {
		panic(err)
	}
	result := sh.Sum(nil)
	return string(result)
}

func Register(userName, password string) (*User, int) {
	user := User{}
	exist := mongodb.Exec(userTable, func(c *mgo.Collection) error {
		return c.Find(bson.M{"username": userName}).One(&user)
	})
	if exist {
		return nil, service.Exist
	}
	user.UserID = bson.NewObjectId().Hex()
	user.UserName = userName
	user.PwdEncrypted = passwordEncrypt(password)
	mongodb.Exec(userTable, func(c *mgo.Collection) error {
		return c.Insert(user)
	})
	return &user, service.Success
}

func Login(userName string, password string) (string, string, int) {
	user := User{}
	exist := mongodb.Exec(userTable, func(c *mgo.Collection) error {
		return c.Find(bson.M{"username": userName}).One(&user)
	})
	if !exist {
		return "", "", service.NotExist
	}
	if passwordEncrypt(password) != user.PwdEncrypted {
		return "", "", service.PwdError
	}
	token := newToken()
	client := myredis.ClusterClient(token)

	//client := myredis.Client()
	/*client := redis.NewClient(&redis.Options{*/
	//Addr:     "127.0.0.1:6379",
	//DB:       0,
	//PoolSize: 100,
	/*})*/
	client.Set(token, user.UserID, 0)

	return user.UserID, token, service.Success
}

func Logout(userID string) int {
	user := User{}
	exist := mongodb.FindByID(userTable, userID, &user)
	if !exist {
		return service.NotExist
	}
	//TODO: update redis
	return service.Success
}
