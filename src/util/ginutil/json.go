package ginutil

import (
	"encoding/json"
	"fmt"
)

type JSONError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

const (
	jsonpFmt = "window.%s&&window.%s(%s)"
)

func generateJSONResponse(data map[string]interface{}) []byte {
	result, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return result
}

func jsonResponseSuccess(data interface{}) []byte {
	resultMap := make(map[string]interface{})
	resultMap["status"] = 1
	resultMap["result"] = data
	resultMap["error"] = nil
	return generateJSONResponse(resultMap)
}

func jsonResponseFailed(err JSONError) []byte {
	resultMap := make(map[string]interface{})
	resultMap["status"] = 0
	resultMap["result"] = nil
	resultMap["error"] = err
	return generateJSONResponse(resultMap)
}

func jspResponse(data []byte, callback string) []byte {
	t := string(data)
	if callback != "" {
		t = fmt.Sprintf(jsonpFmt, callback, callback, t)
	}
	return []byte(t)
}

func JspResponseSuccess(data interface{}, callback string) []byte {
	t := jsonResponseSuccess(data)
	return jspResponse(t, callback)
}

func JspResponseFailed(err JSONError, callback string) []byte {
	t := jsonResponseFailed(err)
	return jspResponse(t, callback)
}
