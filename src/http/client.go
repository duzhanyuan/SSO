package main

import (
	"crypto/rand"
	"crypto/rc4"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}
func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

func genRandomBytes(len int) []byte {
	rb := make([]byte, len)
	rand.Read(rb)
	return rb
}

func genTimestamp(key []byte) string {
	timestamp := time.Now().Unix()
	timestamp_bytes := Int64ToBytes(timestamp)
	c, _ := rc4.NewCipher(key)
	c.XORKeyStream(timestamp_bytes, timestamp_bytes)
	timestamp_str := hex.EncodeToString(timestamp_bytes)
	return timestamp_str
}

func checkTimestamp(timestamp_str string, key []byte) bool {
	c, _ := rc4.NewCipher(key)
	timestamp_bytes, _ := hex.DecodeString(timestamp_str)
	c.XORKeyStream(timestamp_bytes, timestamp_bytes)
	if len(timestamp_bytes) != 8 {
		return false
	}
	timestamp := BytesToInt64(timestamp_bytes)
	current := time.Now().Unix()
	if current > timestamp+60*5 || current+60*5 < timestamp {
		return false
	}
	return true
}

func encrypt_string(str string, key []byte) string {
	bytes := []byte(str)
	c, _ := rc4.NewCipher(key)
	c.XORKeyStream(bytes, bytes)
	return hex.EncodeToString(bytes)
}

func decrypt_string(str string, key []byte) string {
	bytes, _ := hex.DecodeString(str)
	c, _ := rc4.NewCipher(key)
	c.XORKeyStream(bytes, bytes)
	return string(bytes)
}

func doPost(op string, params map[string]string) string {
	addr := "http://127.0.0.1:8000/user/"
	if op == "register" {
		addr = addr + "register"
	} else if op == "login" {
		addr = addr + "login"
	} else if op == "logout" {
		addr = addr + "logout"
	} else if op == "apply_service" {
		addr = addr + "apply_service"
	} else {
		fmt.Println("undefined")
	}
	vals := url.Values{}
	for key, val := range params {
		tmp := []string{}
		tmp = append(tmp, val)
		vals[key] = tmp

	}
	resp, err := http.PostForm(addr, vals)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println(string(body))
	return string(body)
}

func doPost2(op string, params map[string]string) string {
	addr := "http://127.0.0.1:8001/user/"
	if op == "register" {
		addr = addr + "register"
	} else if op == "login" {
		addr = addr + "login"
	} else if op == "logout" {
		addr = addr + "logout"
	} else if op == "apply_service" {
		addr = addr + "apply_service"
	} else {
		fmt.Println("undefined")
	}
	vals := url.Values{}
	for key, val := range params {
		tmp := []string{}
		tmp = append(tmp, val)
		vals[key] = tmp

	}
	resp, err := http.PostForm(addr, vals)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println(string(body))
	return string(body)
}

type Client struct {
	Key        []byte
	SessionKey []byte
	TGT        string
	Name       string
}

func (c *Client) Register(username string, passwd string) {
	hashed_passwd := sha1.Sum([]byte(passwd + ":" + username))
	hashed_passwd_str := hex.EncodeToString(hashed_passwd[:])
	var body map[string]string
	body = make(map[string]string)
	body["username"] = username
	body["password"] = hashed_passwd_str
	//server_register(username, hashed_passwd_str)
	for key, val := range body {
		fmt.Println(key, val)
	}

	doPost("register", body)
}

func (c *Client) Login(username string, passwd string) {
	type Data struct {
		A string `json:a`
		B string `json:b`
	}
	c.Name = username
	hashed_passwd := sha1.Sum([]byte(passwd + ":" + username))
	timestamp_str := genTimestamp(hashed_passwd[:])

	body := make(map[string]string)
	body["username"] = username
	body["password"] = timestamp_str

	data := Data{}
	ret := doPost("login", body)
	json.Unmarshal([]byte(ret), &data)

	A, B := data.A, data.B

	sessionKey, _ := hex.DecodeString(A)
	c1, _ := rc4.NewCipher(hashed_passwd[:])
	c1.XORKeyStream(c.SessionKey, c.SessionKey)

	c.Key = hashed_passwd[:]
	c.SessionKey = sessionKey
	c.TGT = B
}

func (c *Client) ApplyService(service string) {
	type Data struct {
		E string `json:e`
		F string `json:f`
	}

	type Data2 struct {
		H string `json:h`
	}

	timestamp := genTimestamp(c.SessionKey)
	d_str := c.Name + ":" + timestamp
	D := encrypt_string(d_str, c.SessionKey)

	body := make(map[string]string)
	body["service"] = service
	body["TGT"] = c.TGT
	body["D"] = D

	data := Data{}
	ret := doPost("apply_service", body)
	json.Unmarshal([]byte(ret), &data)

	E, F := data.E, data.F

	service_session_key := []byte(decrypt_string(F, c.SessionKey))
	timestamp = genTimestamp(service_session_key)
	G := c.Name + ":" + timestamp
	G = encrypt_string(G, service_session_key)
	body2 := make(map[string]string)
	body2["E"] = E
	body2["G"] = G

	data2 := Data2{}
	//here the url should point to service
	ret2 := doPost2("login", body2)
	json.Unmarshal([]byte(ret2), &data2)

	H := data2.H
	if checkTimestamp(H, service_session_key) {
		return
	} else {
		return
	}

}

func (c *Client) Logout() {
	timestamp := genTimestamp(c.Key)
	params := make(map[string]string)
	params["username"] = c.Name
	params["timestamp"] = timestamp
	doPost("logout", params)
	//res := server_logout("logout", params)
}
func main() {
	client := Client{}
	//client.Register("1@qq.com", "111")
	client.Login("1@qq.com", "111")
}
