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
	Cookies []*http.Cookie
}

type HttpHeader struct {
	Accept string
	AcceptEncoding string
	AcceptLanguage string
	Connection string
	ContentType string
	ContentLength string
	Host string
	Referer string
	UpgradeInsecureRequests string
}

func (this *HttpClient) get(urlStr string, timeout time.Duration, header *HttpHeader) (string, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return "", err
	}
	urlObj, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	jar.SetCookies(urlObj, this.Cookies)
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
		return ``, err
	}

	if header.Host != "" {
		req.Host = header.Host
	}

	if header.Accept != "" {
		req.Header.Add("Accept", header.Accept)
	}

	if header.AcceptEncoding != "" {
		req.Header.Add("Accept-Encoding", header.AcceptEncoding)
	}
	if header.AcceptLanguage != "" {
		req.Header.Add("Accept-Language", header.AcceptLanguage)
	}

	if header.Connection != "" {
		req.Header.Add("Connection", header.Connection)
	}

	if header.ContentType != "" {
		req.Header.Add("Content-Type", header.ContentType)
	}

	if header.ContentLength != "" {
		req.Header.Add("Content-Length", header.ContentLength)
	}

	if header.Referer != "" {
		req.Header.Add("Referer", header.Referer)
	}

	if header.UpgradeInsecureRequests != "" {
		req.Header.Add("Upgrade-Insecure-Requests", header.UpgradeInsecureRequests)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")

	resp, err := client.Do(req)
	if (err != nil && err.Error() != "Get /: Cannot Redirect") {
		return ``, err
	}

	if len(resp.Cookies()) > 0 {
		cookies := resp.Cookies()
		for _, v := range cookies {
			this.Cookies = append(this.Cookies, v)
		}
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return string(body), nil
}

func (this *HttpClient) post(urlStr string, data []byte, timeout time.Duration, header *HttpHeader) (string, error) {

	jar, err := cookiejar.New(nil)
	if err != nil {
		return "", err
	}
	urlObj, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	jar.SetCookies(urlObj, this.Cookies)

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
		return ``, err
	}
	if header.Host != "" {
		req.Host = header.Host
	}

	if header.Accept != "" {
		req.Header.Add("Accept", header.Accept)
	}

	if header.AcceptEncoding != "" {
		req.Header.Add("Accept-Encoding", header.AcceptEncoding)
	}
	if header.AcceptLanguage != "" {
		req.Header.Add("Accept-Language", header.AcceptLanguage)
	}

	if header.Connection != "" {
		req.Header.Add("Connection", header.Connection)
	}

	if header.ContentType != "" {
		req.Header.Add("Content-Type", header.ContentType)
	}

	if header.ContentLength != "" {
		req.Header.Add("Content-Length", header.ContentLength)
	}

	if header.Referer != "" {
		req.Header.Add("Referer", header.Referer)
	}

	if header.UpgradeInsecureRequests != "" {
		req.Header.Add("Upgrade-Insecure-Requests", header.UpgradeInsecureRequests)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if len(resp.Cookies()) > 0 {
		cookies := resp.Cookies()
		for _, v := range cookies {
			this.Cookies = append(this.Cookies, v)
		}
	}
	if (err != nil && err.Error() != "Get /: Cannot Redirect") {
		return ``, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}