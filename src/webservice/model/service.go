package model

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"service/errormap"
	"strings"
	"util"
)

var (
	Service_key []byte
)

func Service_login(E string, G string) (string, int) {
	e := util.DecryptString(E, Service_key)
	es := strings.Split(e, ":")
	if len(es) != 2 {
		return "", errormap.IllegalTS
	}
	session_key, _ := hex.DecodeString(es[1])
	g := util.DecryptString(G, session_key)
	gs := strings.Split(g, ":")
	if len(gs) != 2 {
		return "", errormap.IllegalTS
	}
	if !util.CheckTimestamp(gs[1], session_key) {
		return "", errormap.IllegalTS
	}
	return util.GenTimestamp(session_key), errormap.Success
}

func Service_register(name string) {
	addr := "http://127.0.0.1:8000/user/register_service"
	resp, err := http.PostForm(addr, url.Values{"name": {name}})
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	type Data struct {
		Key string `json:"key"`
	}
	data := Data{}
	json.Unmarshal([]byte(string(body)), &data)

	key_str := data.Key
	//key_str := server_register_service(name)

	Service_key, _ = hex.DecodeString(key_str)
}
