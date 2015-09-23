package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webservice/handler"
)

func regRouter(router *gin.Engine) {
	handler.Register(router.Group("/"))
}

func main() {
	router := gin.Default()
	regRouter(router)
	listenAddr := "127.0.0.1:8001"

	err := http.ListenAndServe(listenAddr, router)
	if err != nil {
		panic(err)
	}
}
