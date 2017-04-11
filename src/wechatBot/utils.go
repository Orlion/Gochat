package wechatBot

import (
	"strconv"
	"time"
	"net/http"
)

func GetUnixTime() string{
	return strconv.FormatInt(time.Now().UnixNano() / 1000000, 10)
}

/**
 * 处理cookie转成字符串字典
 */
func Cookies2Map(cookies []*http.Cookie) map[string]string {
	result := map[string]string{}
	for _, v := range cookies {
		result[v.Name] = v.Value
	}

	return result
}

/**
 * cookie 转成字符串
 */
func Cookies2String(cookies []*http.Cookie) string{
	result := ""
	for _, v := range cookies {
		result += v.Name + "=" + v.Value + ";"
	}

	return result
}