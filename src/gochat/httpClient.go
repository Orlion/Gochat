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
	"io"
)

type httpClient struct {
	Cookies []*http.Cookie
}

type httpHeader struct {
	Accept 					string
	AcceptEncoding			string
	AcceptLanguage 			string
	Connection 				string
	ContentType 			string
	ContentLength 			string
	Origin 					string
	Host 					string
	Referer 				string
	UpgradeInsecureRequests	string
}

// 发起get请求
func (httpClient *httpClient) get(urlStr string, timeout time.Duration, header *httpHeader) (string, error) {

	client, err := httpClient.getClient(urlStr, timeout)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", err
	}

	httpClient.handleHeader(req, header)

	resp, err := client.Do(req)
	if (err != nil && err.Error() != "Get /: Cannot Redirect") {
		return "", err
	}

	if len(resp.Cookies()) > 0 {
		cookies := resp.Cookies()
		for _, v := range cookies {
			httpClient.Cookies = append(httpClient.Cookies, v)
		}
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return string(body), nil
}

// 发起post请求
func (httpClient *httpClient) post(urlStr string, data []byte, timeout time.Duration, header *httpHeader) (string, error) {

	client, err := httpClient.getClient(urlStr, timeout)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", urlStr, bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	httpClient.handleHeader(req, header)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if len(resp.Cookies()) > 0 {
		cookies := resp.Cookies()
		for _, v := range cookies {
			httpClient.Cookies = append(httpClient.Cookies, v)
		}
	}
	if (err != nil && err.Error() != "Get /: Cannot Redirect") {
		return "", err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

// 上传文件请求
func (httpClient *httpClient) upload(urlStr string, data io.Reader, timeout time.Duration, header *httpHeader) (string, error) {

	req, err := http.NewRequest("POST", urlStr, data)
	if err != nil {
		return "", err
	}

	httpClient.handleHeader(req, header)

	client, err := httpClient.getClient(urlStr, timeout)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if len(resp.Cookies()) > 0 {
		cookies := resp.Cookies()
		for _, v := range cookies {
			httpClient.Cookies = append(httpClient.Cookies, v)
		}
	}
	if (err != nil && err.Error() != "Get /: Cannot Redirect") {
		return "", err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

// 生成http.Client对象
func (httpClient *httpClient) getClient(urlStr string, timeout time.Duration) (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	urlObj, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	jar.SetCookies(urlObj, httpClient.Cookies)

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

	return client, nil
}

// 处理header
func (httpClient *httpClient) handleHeader(req *http.Request, header *httpHeader) {

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

	if header.Origin != "" {
		req.Header.Add("Origin", header.Origin)
	}

	if header.Referer != "" {
		req.Header.Add("Referer", header.Referer)
	}

	if header.UpgradeInsecureRequests != "" {
		req.Header.Add("Upgrade-Insecure-Requests", header.UpgradeInsecureRequests)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")
}