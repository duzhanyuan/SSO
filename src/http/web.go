package main

import (
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
	addr := "http://127.0.0.1:10000/user/register_service"
	resp, err := http.PostForm(addr, url.Values{"name": {"test"}})
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

	data := Data{}
	json.Unmarshal([]byte(string(body)), &data)
	model.Service_key = ([]byte)(data.Key)

	listenAddr := "0.0.0.0:9999"
	err = http.ListenAndServe(listenAddr, router)
	if err != nil {
		panic(err)
	}
}
