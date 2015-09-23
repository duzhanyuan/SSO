package main

import (
	"github.com/gin-gonic/gin"
	"handler/user"
	"net/http"
	"service/mongodb"
	"service/myredis"
)

func initDb() {
	mongodb.Init("127.0.0.1:27017")
	//myredis.Init("127.0.0.1:6379", 0)

	/* urls := make(map[uint32]string)*/
	//urls[0] = "127.0.0.1:6379"
	//urls[1] = "127.0.0.1:6380"
	//urls[2] = "127.0.0.1:6381"
	//urls[3] = "127.0.0.1:6382"
	/*urls[4] = "127.0.0.1:6383"*/
	//myredis.InitCluster(urls)
	addrs := make(map[string]int64)
	addrs["127.0.0.1:6379"] = 0
	addrs["127.0.0.1:6380"] = 0
	addrs["127.0.0.1:6381"] = 0
	addrs["127.0.0.1:6382"] = 0
	addrs["127.0.0.1:6383"] = 0
	nodes := make(map[string]string)
	nodes["192.168.101.1"] = "127.0.0.1:6379"
	nodes["192.168.101.2"] = "127.0.0.1:6379"
	nodes["192.168.101.3"] = "127.0.0.1:6380"
	nodes["192.168.101.4"] = "127.0.0.1:6380"
	nodes["192.168.101.5"] = "127.0.0.1:6381"
	nodes["192.168.101.6"] = "127.0.0.1:6381"
	nodes["192.168.101.7"] = "127.0.0.1:6382"
	nodes["192.168.101.8"] = "127.0.0.1:6382"
	nodes["192.168.101.9"] = "127.0.0.1:6383"
	nodes["192.168.101.10"] = "127.0.0.1:6383"
	myredis.InitCluster(addrs, nodes)
}

func regRouter(router *gin.Engine) {
	user.Register(router.Group("/"))
}

func main() {
	initDb()
	router := gin.Default()
	regRouter(router)
	listenAddr := "127.0.0.1:8000"
	err := http.ListenAndServe(listenAddr, router)
	if err != nil {
		panic(err)
	}
}
