package handler

import (
	"github.com/gin-gonic/gin"
	"service/errormap"
	"util/ginutil"
	"webservice/model"
)

func Register(router *gin.RouterGroup) {
	group := router.Group("/web_service")
	group.POST("/register", register)
	group.POST("/login", login)
}

func register(c *gin.Context) {
	type FormData struct {
		Name string `json:"name"`
	}
	var form FormData
	if c.Bind(&form) != nil {
		return
	}
	model.Service_register(form.Name)
	ginutil.ResponseJSONSuccess(c, nil)
}

func login(c *gin.Context) {
	type FormData struct {
		E string `json:"e"`
		G string `json:"g"`
	}
	type Data struct {
		Key string `json:"key"`
	}
	var form FormData
	if c.Bind(&form) != nil {
		return
	}
	key, code := model.Service_login(form.E, form.G)
	if code != errormap.Success {
		ginutil.ResponseJSONFailed(c, ginutil.JSONError{Code: code, Msg: errormap.ErrorMsg(code)})
	} else {
		data := Data{
			Key: key,
		}
		ginutil.ResponseJSONSuccess(c, data)
	}
}
