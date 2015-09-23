package main

import (
	"github.com/gin-gonic/gin"
	"handler/user"
	"net/http"
	"service/mongodb"
	"service/monitor"
	"service/myredis"
	"util"
)

func initDb() {
	var addrs map[string]int64
	var nodes map[string]string
	if util.IsDevelopment() {
		mongodb.Init("127.0.0.1:27017")
		addrs = make(map[string]int64)
		addrs["127.0.0.1:6379"] = 0
		/*addrs["127.0.0.1:6380"] = 0*/
		//addrs["127.0.0.1:6381"] = 0
		//addrs["127.0.0.1:6382"] = 0
		/*addrs["127.0.0.1:6383"] = 0*/
		nodes = make(map[string]string)
		nodes["192.168.101.1"] = "127.0.0.1:6379"
		nodes["192.168.101.2"] = "127.0.0.1:6379"
		/*     nodes["192.168.101.3"] = "127.0.0.1:6380"*/
		//nodes["192.168.101.4"] = "127.0.0.1:6380"
		//nodes["192.168.101.5"] = "127.0.0.1:6381"
		//nodes["192.168.101.6"] = "127.0.0.1:6381"
		//nodes["192.168.101.7"] = "127.0.0.1:6382"
		//nodes["192.168.101.8"] = "127.0.0.1:6382"
		//nodes["192.168.101.9"] = "127.0.0.1:6383"
		/*nodes["192.168.101.10"] = "127.0.0.1:6383"*/
	} else {
		mongodb.Init("192.168.1.16:27017")
		addrs = make(map[string]int64)
		addrs["192.168.1.11:6379"] = 0
		addrs["192.168.1.12:6379"] = 0
		addrs["192.168.1.13:6379"] = 0
		addrs["192.168.1.14:6379"] = 0
		addrs["192.168.1.15:6379"] = 0

		nodes = make(map[string]string)
		nodes["192.168.101.1"] = "192.168.1.11:6379"
		nodes["192.168.101.2"] = "192.168.1.11:6379"
		nodes["192.168.101.3"] = "192.168.1.12:6379"
		nodes["192.168.101.4"] = "192.168.1.12:6379"
		nodes["192.168.101.5"] = "192.168.1.13:6379"
		nodes["192.168.101.6"] = "192.168.1.13:6379"
		nodes["192.168.101.7"] = "192.168.1.14:6379"
		nodes["192.168.101.8"] = "192.168.1.14:6379"
		nodes["192.168.101.9"] = "192.168.1.15:6379"
		nodes["192.168.101.10"] = "192.168.1.15:6379"
	}
	//myredis.Init("127.0.0.1:6379", 0)
	myredis.InitCluster(addrs, nodes)
}

func regRouter(router *gin.Engine) {
	user.Register(router.Group("/"))
}

func main() {
	initDb()
	go monitor.Work()
	router := gin.Default()
	regRouter(router)

	listenAddr := "127.0.0.1:8000"
	err := http.ListenAndServe(listenAddr, router)
	if err != nil {
		panic(err)
	}
}
