package user

import (
	"github.com/gin-gonic/gin"
	"model/user"
	"service/errormap"
	"util"
	"util/ginutil"
)

func Register(router *gin.RouterGroup) {
	group := router.Group("/user")
	group.POST("/register", register)
	group.POST("/register_service", registerService)
	group.POST("/login", login)
	group.POST("/logout", logout)
}

func register(c *gin.Context) {
	type FormData struct {
		UserName string `form:"username"`
		Password string `form:"password"`
	}
	type Data struct {
		Token string
	}
	var form FormData
	if c.Bind(&form) != nil {
		return
	}

	code := user.Register(form.UserName, form.Password)
	if code != errormap.Success {
		ginutil.ResponseJSONFailed(c, ginutil.JSONError{Code: code, Msg: errormap.ErrorMsg(code)})
	} else {
		ginutil.ResponseJSONSuccess(c, nil)
	}
}

func registerService(c *gin.Context) {
	type FormData struct {
		Name string `form:"name"`
	}
	type Data struct {
		Key string `form:"key"`
	}
	var form FormData
	if c.Bind(&form) != nil {
		return
	}

	key, code := user.Register(form.UserName, form.Password)
	if code != errormap.Success {
		ginutil.ResponseJSONFailed(c, ginutil.JSONError{Code: code, Msg: errormap.ErrorMsg(code)})
	} else {
		data := Data{
			Key: key,
		}
		ginutil.ResponseJSONSuccess(c, data)
	}
}

func login(c *gin.Context) {
	type FormData struct {
		UserName string `form:"username"`
		Password string `form:"password"`
	}
	type Data struct {
		A string `json:a`
		B string `json:b`
	}

	var form FormData
	if c.Bind(&form) != nil {
		return
	}
	a, b, code := user.Login(form.UserName, form.Password)
	if code != errormap.Success {
		ginutil.ResponseJSONFailed(c, ginutil.JSONError{Code: code, Msg: errormap.ErrorMsg(code)})
	} else {
		data := Data{
			A: a,
			B: b,
		}
		util.Authorsize(c, userID, token)
		ginutil.ResponseJSONSuccess(c, data)
	}
}

func logout(c *gin.Context) {
	type Data struct {
		Token string
	}
	userID := util.GetUserID(c)
	code := user.Logout(userID)
	if code != errormap.Success {
		ginutil.ResponseJSONFailed(c, ginutil.JSONError{Code: code, Msg: errormap.ErrorMsg(code)})
	} else {
		data := Data{
			Token: "klsj",
		}
		ginutil.ResponseJSONSuccess(c, data)
	}
}
