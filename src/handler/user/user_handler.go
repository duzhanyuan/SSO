package user

import (
	"github.com/gin-gonic/gin"
	"model/user"
	"service"
	"util"
	"util/ginutil"
)

func Register(router *gin.RouterGroup) {
	group := router.Group("/user")
	group.POST("/register", register)
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

	_, code := user.Register(form.UserName, form.Password)
	if code != service.Success {
		ginutil.ResponseJSONFailed(c, ginutil.JSONError{Code: code, Msg: service.ErrorMsg(code)})
	} else {
		data := Data{
			Token: "klsj",
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
		Token string
	}

	var form FormData
	if c.Bind(&form) != nil {
		return
	}
	userID, token, code := user.Login(form.UserName, form.Password)
	if code != service.Success {
		ginutil.ResponseJSONFailed(c, ginutil.JSONError{Code: code, Msg: service.ErrorMsg(code)})
	} else {
		data := Data{
			Token: token,
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
	if code != service.Success {
		ginutil.ResponseJSONFailed(c, ginutil.JSONError{Code: code, Msg: service.ErrorMsg(code)})
	} else {
		data := Data{
			Token: "klsj",
		}
		ginutil.ResponseJSONSuccess(c, data)
	}
}
