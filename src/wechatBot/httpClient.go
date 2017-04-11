package wechatBot

import (
	"net/http"
	"net"
	"time"
	"errors"
	"io/ioutil"
	"net/url"
	"strings"
	"fmt"
)

type HttpClient struct {
	HttpHeader HttpHeader
	Timeout time.Duration
}

type HttpHeader struct {
	Accept string
	AcceptEncoding string
	AcceptLanguage string
	Connection string
	ContentType string
	ContentLength string
	Cookie string
	Host string
	Referer string
	UpgradeInsecureRequests string
	UserAgent string
}

func (this *HttpClient) Get(url string) (string, []*http.Cookie, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(this.Timeout)
				c, err := net.DialTimeout(netw, addr, this.Timeout) //连接超时时间
				if err != nil {
					return nil, err
				}

				c.SetDeadline(deadline)
				return c, nil
			},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("Cannot Redirect")
		},
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ``, nil, err
	}

	if this.HttpHeader.Host != "" {
		req.Host = this.HttpHeader.Host
	}

	if this.HttpHeader.Accept != "" {
		req.Header.Add("Accept", this.HttpHeader.Accept)
	}

	if this.HttpHeader.AcceptEncoding != "" {
		req.Header.Add("Accept-Encoding", this.HttpHeader.AcceptEncoding)
	}
	if this.HttpHeader.AcceptLanguage != "" {
		req.Header.Add("Accept-Language", this.HttpHeader.AcceptLanguage)
	}

	if this.HttpHeader.Connection != "" {
		req.Header.Add("Connection", this.HttpHeader.Connection)
	}

	if this.HttpHeader.ContentType != "" {
		req.Header.Add("Content-Type", this.HttpHeader.ContentType)
	}

	if this.HttpHeader.ContentLength != "" {
		req.Header.Add("Content-Length", this.HttpHeader.ContentLength)
	}

	if this.HttpHeader.Cookie != "" {
		req.Header.Add("Cookie", this.HttpHeader.Cookie)
	}

	if this.HttpHeader.Referer != "" {
		req.Header.Add("Referer", this.HttpHeader.Referer)
	}

	if this.HttpHeader.UpgradeInsecureRequests != "" {
		req.Header.Add("Upgrade-Insecure-Requests", this.HttpHeader.UpgradeInsecureRequests)
	}

	if this.HttpHeader.UserAgent != "" {
		req.Header.Add("User-Agent", this.HttpHeader.UserAgent)
	}

	resp, err := client.Do(req)
	if (err != nil && err.Error() != "Get /: Cannot Redirect") {
		return ``, nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return string(body), resp.Cookies(), nil
}

func (this *HttpClient)Post(url string, data url.Values) (string, []*http.Cookie, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(this.Timeout)
				c, err := net.DialTimeout(netw, addr, this.Timeout) //连接超时时间
				if err != nil {
					return nil, err
				}

				c.SetDeadline(deadline)
				return c, nil
			},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("Cannot Redirect")
		},
	}
	postData := ioutil.NopCloser(strings.NewReader(data.Encode()))
	fmt.Println(postData);
	req, err := http.NewRequest("POST", url, postData)
	if err != nil {
		return ``, nil, err
	}
	if this.HttpHeader.Host != "" {
		req.Host = this.HttpHeader.Host
	}

	if this.HttpHeader.Accept != "" {
		req.Header.Add("Referer", this.HttpHeader.Accept)
	}

	if this.HttpHeader.AcceptEncoding != "" {
		req.Header.Add("Accept-Encoding", this.HttpHeader.AcceptEncoding)
	}
	if this.HttpHeader.AcceptLanguage != "" {
		req.Header.Add("Accept-Language", this.HttpHeader.AcceptLanguage)
	}

	if this.HttpHeader.Connection != "" {
		req.Header.Add("Connection", this.HttpHeader.Connection)
	}

	if this.HttpHeader.ContentType != "" {
		req.Header.Add("Content-Type", this.HttpHeader.ContentType)
	}

	if this.HttpHeader.Cookie != "" {
		req.Header.Add("Cookie", this.HttpHeader.Cookie)
	}

	if this.HttpHeader.Host != "" {
		req.Header.Add("Host", this.HttpHeader.Host)
	}

	if this.HttpHeader.Referer != "" {
		req.Header.Add("Referer", this.HttpHeader.Referer)
	}

	if this.HttpHeader.UpgradeInsecureRequests != "" {
		req.Header.Add("Upgrade-Insecure-Requests", this.HttpHeader.UpgradeInsecureRequests)
	}

	if this.HttpHeader.UserAgent != "" {
		req.Header.Add("User-Agent", this.HttpHeader.UserAgent)
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if (err != nil && err.Error() != "Get /: Cannot Redirect") {
		return ``, nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), resp.Cookies(), nil
}

func (this *HttpClient)PostStr(url string, data string) (string, []*http.Cookie, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(this.Timeout)
				c, err := net.DialTimeout(netw, addr, this.Timeout) //连接超时时间
				if err != nil {
					return nil, err
				}

				c.SetDeadline(deadline)
				return c, nil
			},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("Cannot Redirect")
		},
	}
	postData := ioutil.NopCloser(strings.NewReader(data))
	fmt.Println(postData)
	req, err := http.NewRequest("POST", url, postData)
	if err != nil {
		return ``, nil, err
	}
	if this.HttpHeader.Host != "" {
		req.Host = this.HttpHeader.Host
	}

	if this.HttpHeader.Accept != "" {
		req.Header.Add("Referer", this.HttpHeader.Accept)
	}

	if this.HttpHeader.AcceptEncoding != "" {
		req.Header.Add("Accept-Encoding", this.HttpHeader.AcceptEncoding)
	}
	if this.HttpHeader.AcceptLanguage != "" {
		req.Header.Add("Accept-Language", this.HttpHeader.AcceptLanguage)
	}

	if this.HttpHeader.Connection != "" {
		req.Header.Add("Connection", this.HttpHeader.Connection)
	}

	if this.HttpHeader.ContentType != "" {
		req.Header.Add("Content-Type", this.HttpHeader.ContentType)
	}

	if this.HttpHeader.Cookie != "" {
		req.Header.Add("Cookie", this.HttpHeader.Cookie)
	}

	if this.HttpHeader.Host != "" {
		req.Header.Add("Host", this.HttpHeader.Host)
	}

	if this.HttpHeader.Referer != "" {
		req.Header.Add("Referer", this.HttpHeader.Referer)
	}

	if this.HttpHeader.UpgradeInsecureRequests != "" {
		req.Header.Add("Upgrade-Insecure-Requests", this.HttpHeader.UpgradeInsecureRequests)
	}

	if this.HttpHeader.UserAgent != "" {
		req.Header.Add("User-Agent", this.HttpHeader.UserAgent)
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if (err != nil && err.Error() != "Get /: Cannot Redirect") {
		return ``, nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), resp.Cookies(), nil
}
