package gochat

import (
	"strconv"
	"time"
	"net/http"
	"net/url"
)

type Utils struct {

}

func (this *Utils) getUnixMsTime() string {
	return strconv.FormatInt(time.Now().UnixNano() / 1000000, 10)
}

func (this *Utils) getUnixTime() string {
	return strconv.FormatInt(time.Now().Unix() / 1000000, 10)
}

func (this *Utils) cookies2Map(cookies []*http.Cookie) map[string]string {
	result := map[string]string{}
	for _, v := range cookies {
		result[v.Name] = v.Value
	}

	return result
}

func (this *Utils) cookies2String(cookies []*http.Cookie) string{
	result := ""
	for _, v := range cookies {
		result += v.Name + "=" + v.Value + ";"
	}

	return result
}

func (this *Utils) userName2Id(userName string) string {
	r := ""
	for _,v := range []byte(userName) {
		r += strconv.Itoa(int(v))
	}

	return this.substr(r, 0, 9)
}

func (this *Utils) substr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < 0 || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}

func (this *Utils) getHostByUrl(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "wx2.qq.com"
	}

	return u.Host
}