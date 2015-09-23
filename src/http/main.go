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

	urls := make(map[int64]string)
	urls[0] = "127.0.0.1:6379"
	urls[1] = "127.0.0.1:6380"
	urls[2] = "127.0.0.1:6381"
	urls[3] = "127.0.0.1:6382"
	urls[4] = "127.0.0.1:6383"
	myredis.InitCluster(urls)
}

func regRouter(router *gin.Engine) {
	user.Register(router.Group("/"))
}

func main() {
	initDb()
	router := gin.Default()
	regRouter(router)
	listenAddr := "192.168.1.114:9876"
	err := http.ListenAndServe(listenAddr, router)
	if err != nil {
		panic(err)
	}
}
