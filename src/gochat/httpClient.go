package gochat

import (
	"net/http"
	"net"
	"time"
	"errors"
	"io/ioutil"
	"bytes"
	"net/http/cookiejar"
	"net/url"
)

type HttpClient struct {
	HttpHeader HttpHeader
}

type HttpHeader struct {
	Accept string
	AcceptEncoding string
	AcceptLanguage string
	Connection string
	ContentType string
	ContentLength string
	Cookies []*http.Cookie
	Host string
	Referer string
	UpgradeInsecureRequests string
}

func (this *HttpClient) Get(urlStr string, timeout time.Duration) (string, []*http.Cookie, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return "", nil, err
	}
	urlObj, err := url.Parse(urlStr)
	if err != nil {
		return "", nil, err
	}
	jar.SetCookies(urlObj, this.HttpHeader.Cookies)
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(timeout)
				c, err := net.DialTimeout(netw, addr, timeout) //连接超时时间
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
		Jar: jar,
	}
	req, err := http.NewRequest("GET", urlStr, nil)
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

	if this.HttpHeader.Referer != "" {
		req.Header.Add("Referer", this.HttpHeader.Referer)
	}

	if this.HttpHeader.UpgradeInsecureRequests != "" {
		req.Header.Add("Upgrade-Insecure-Requests", this.HttpHeader.UpgradeInsecureRequests)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")

	resp, err := client.Do(req)
	if (err != nil && err.Error() != "Get /: Cannot Redirect") {
		return ``, nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return string(body), resp.Cookies(), nil
}

func (this *HttpClient)Post(urlStr string, data []byte, timeout time.Duration) (string, []*http.Cookie, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return "", nil, err
	}
	urlObj, err := url.Parse(urlStr)
	if err != nil {
		return "", nil, err
	}
	jar.SetCookies(urlObj, this.HttpHeader.Cookies)
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(timeout)
				c, err := net.DialTimeout(netw, addr, timeout) //连接超时时间
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
		Jar: jar,
	}
	req, err := http.NewRequest("POST", urlStr, bytes.NewReader(data))
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

	if this.HttpHeader.Host != "" {
		req.Header.Add("Host", this.HttpHeader.Host)
	}

	if this.HttpHeader.Referer != "" {
		req.Header.Add("Referer", this.HttpHeader.Referer)
	}

	if this.HttpHeader.UpgradeInsecureRequests != "" {
		req.Header.Add("Upgrade-Insecure-Requests", this.HttpHeader.UpgradeInsecureRequests)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if (err != nil && err.Error() != "Get /: Cannot Redirect") {
		return ``, nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), resp.Cookies(), nil
}