package ginutil

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	protoContentType = "application/x-protobuf"
	htmlContentType  = "text/html"
	jsonContentType  = "application/json"
	jspContentType   = "application/javascript"
	textContentType  = "text/plain; charset=utf-8"
	PlatformAndroid  = "android"
	PlatformIOS      = "ios"
)

// ErrorResponse return 400 with a base.Error obj
/*func ErrorResponse(c *gin.Context, code base.ErrorCode, detail string) {*/
//log.Error.Println(spew.Sprintf("ErrorResponse: %v", detail))
//c.Data(http.StatusBadRequest, protoContentType, util.MustMarshal(&base.Error{Code: code.Enum(), Detail: proto.String(detail)}))

//userID := GetUserID(c)
//redisClient := myredis.Client(userID)
//key := util.GetRequestBannerKey(userID, "400")

//redisClient.Incr(key)
//redisClient.Expire(key, time.Duration(time.Hour*24*3))
//}

// Response return 200 with

func ResponseNotFound(c *gin.Context) {
	c.Data(http.StatusNotFound, protoContentType, []byte{})
}

func ResponseHTML(c *gin.Context, html string) {
	c.Data(http.StatusOK, htmlContentType, []byte(html))
}

func responseJSON(c *gin.Context, json []byte) {
	c.Data(http.StatusOK, jsonContentType, json)
}

func ResponseJSONSuccess(c *gin.Context, data interface{}) {
	responseJSON(c, jsonResponseSuccess(data))
}

func ResponseJSONFailed(c *gin.Context, err JSONError) {
	responseJSON(c, jsonResponseFailed(err))
}

func responseJSP(c *gin.Context, json []byte) {
	c.Data(http.StatusOK, jspContentType, json)
}

func ResponseJSPSuccess(c *gin.Context, data interface{}, callback string) {
	responseJSP(c, JspResponseSuccess(data, callback))
}

func ResponseJSPFailed(c *gin.Context, err JSONError, callback string) {
	responseJSP(c, JspResponseFailed(err, callback))
}

func ResponseText(c *gin.Context, text string) {
	c.Data(http.StatusOK, textContentType, []byte(text))
}
