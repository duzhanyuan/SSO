package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"model/user"
	"service/errormap"
	"service/monitor"
	"util/ginutil"
)

func Register(router *gin.RouterGroup) {
	group := router.Group("/user")
	group.POST("/register", register)
	group.POST("/register_service", registerService)
	group.POST("/login", login)
	group.POST("/logout", logout)
	group.GET("/monitor", performance)
	group.POST("/apply_service", applyService)
}

func performance(c *gin.Context) {
	type Data struct {
		Nodes []int64 `json:"nodes"`
	}
	nodes := monitor.GetAllData()
	data := Data{
		Nodes: nodes,
	}
	ginutil.ResponseJSONSuccess(c, data)
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

	fmt.Println("formname", form.Name)
	key, code := user.RegisterService(form.Name)
	fmt.Println("key:", key, code)
	if code != errormap.Success {
		ginutil.ResponseJSONFailed(c, ginutil.JSONError{Code: code, Msg: errormap.ErrorMsg(code)})
	} else {
		c.JSON(200, gin.H{"key": key})
		//ginutil.ResponseJSONSuccess(c, data)
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
		c.JSON(200, gin.H{"A": a, "B": b})
		//util.Authorsize(c, userID, token)
	}
}

func logout(c *gin.Context) {
	type FormData struct {
		TGT       string `form:"TGT"`
		TimeStamp string `form:"timestamp"`
	}
	var form FormData
	if c.Bind(&form) != nil {
		return
	}
	code := user.Logout(form.TGT, form.TimeStamp)
	if code != errormap.Success {
		ginutil.ResponseJSONFailed(c, ginutil.JSONError{Code: code, Msg: errormap.ErrorMsg(code)})
	} else {
		ginutil.ResponseJSONSuccess(c, nil)
	}
}

func applyService(c *gin.Context) {
	type FormData struct {
		Service string `json:"Service"`
		TGT     string `json:"TGT"`
		D       string `json:"D"`
	}
	var form FormData
	if c.Bind(&form) != nil {
		return
	}
	fmt.Println("before ef", form.Service, form.TGT, form.D)
	E, F, code := user.Apply(form.Service, form.TGT, form.D)
	if code != errormap.Success {
		ginutil.ResponseJSONFailed(c, ginutil.JSONError{Code: code, Msg: errormap.ErrorMsg(code)})
	} else {
		fmt.Println(E, F, "EF")
		c.JSON(200, gin.H{"E": E, "F": F})
	}

}
