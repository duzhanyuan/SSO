package util

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

const (
	cookieKey   = "Y2hpdHVBdXRo"
	expiresDays = 7
	cookieFmt   = "%s authorized %s"
)

func setCookie(c *gin.Context, cookieStr string) {
	expiresTime := time.Now().Add(time.Hour * 24 * expiresDays)
	cookie := http.Cookie{
		Name:    cookieKey,
		Value:   cookieStr,
		Expires: expiresTime,
	}
	http.SetCookie(c.Writer, &cookie)
}

func getCookie(c *gin.Context) (string, error) {
	cookie, err := c.Request.Cookie(cookieKey)
	if err != nil {
		return "", err
	}
	return cookie.Value, err
}

func Authorsize(c *gin.Context, userID string, token string) {
	cookieStr := fmt.Sprintf(cookieFmt, userID, token)
	cookieStr = base64.StdEncoding.EncodeToString([]byte(cookieStr))
	setCookie(c, cookieStr)
}

func CheckAuth(c *gin.Context) (bool, string) {
	cookie, err := getCookie(c)
	if err != nil {
		return false, ""
	}
	tmp, _ := base64.StdEncoding.DecodeString(cookie)
	tmpStr := strings.Split(string(tmp), " ")
	adminID := tmpStr[0]
	//token := tmpStr[2]
	return true, adminID
	/*if admin.CheckToken(adminID, token) {*/
	//return true, adminID
	/*}*/
	//return false, ""
}

func GetUserID(c *gin.Context) string {
	suc, id := CheckAuth(c)
	if !suc {
		return ""
	}
	return id
}
