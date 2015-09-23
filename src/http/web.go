package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"webservice/handler"
	"webservice/model"
)

func regRouter(router *gin.Engine) {
	handler.Register(router.Group("/"))
}

func main() {
	router := gin.Default()
	regRouter(router)
	//listenAddr := "127.0.0.1:10000/user/register_service"
	addr := "http://127.0.0.1:8000/user/register_service"
	resp, err := http.PostForm(addr, url.Values{"name": {"test1"}})
	type Data struct {
		Key string `json:"key"`
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("body", string(body))

	data := Data{}
	json.Unmarshal([]byte(string(body)), &data)
	model.Service_key, _ = hex.DecodeString(data.Key)
	fmt.Println("service_key:", model.Service_key)

	listenAddr := "0.0.0.0:9999"
	err = http.ListenAndServe(listenAddr, router)
	if err != nil {
		panic(err)
	}
}
