package wechatBot

import (
	"regexp"
	"time"
	"fmt"
	"strings"
	"net/http"
	"encoding/xml"
	"strconv"
	"math/rand"
)

type  Wechat struct {
	HttpClient HttpClient
	Uuid string
	RedirectUrl string
	Cookies []*http.Cookie
	DeviceID string
	Sid string
	Skey string
	Uin string
	PassTicket string
}

type Error struct {
	Ret string  `xml:"ret"`
	Message string  `xml:"message"`
	Skey string  `xml:"skey"`
	Wxsid string  `xml:"wxsid"`
	Wxuin string  `xml:"wxuin"`
	PassTicket string  `xml:"pass_ticket"`
	Isgrayscale string `xml:"isgrayscale"`
}
/**
 * 初始化
 */
func (this *Wechat) Init() {
	// 初始化httpHeader
	var httpHeader = HttpHeader{
		"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
		"gzip, deflate, sdch, br",
		"zh-CN,zh;q=0.8",
		"keep-alive",
		"",
		"",
		"",
		"login.wx2.qq.com",
		"https://wx.qq.com/",
		"1",
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
	}

	this.HttpClient = HttpClient{
		httpHeader,
		time.Second * 30,
	}
}

/**
 * 获取uuid
 */
func (this *Wechat) GetUuid() (string, error) {
	getUuidApiUrl := Config["getUuidApi"] + GetUnixTime()
	content, _, err := this.HttpClient.Get(getUuidApiUrl)
	if err != nil {
		return ``, err
	}

	reg, err := regexp.Compile(`window.QRLogin.code = 200; window.QRLogin.uuid = "(.+)"`)
	if err != nil {
		return ``, err
	}

	uuid := reg.FindSubmatch([]byte(content))
	if len(uuid) != 2 {
		return ``, nil
	}
	this.Uuid = string(uuid[1])
	return this.Uuid, nil
}

/**
 * 获取redirectUrl
 */
func (this *Wechat) GetRedirectUrl() (err error){
	var tip int = 1
	for  {

		redirectUrl, err := this.Polling(tip)
		if err != nil {
			fmt.Println(err)
		}

		if "" == redirectUrl {
			continue
		}

		if "201" == redirectUrl {
			fmt.Println("用户已扫码,等待确认中...")
			time.Sleep(time.Second * time.Duration(1))
			tip = 0
			continue
		}

		this.RedirectUrl = redirectUrl
		return nil
	}
}

/**
 * 轮询,直到用户在手机微信上确认登录，确认登录后会返回redirectUrl
 */
func (this *Wechat) Polling(tip int) (url string, err error){
	loginPollApi := strings.Replace(Config["login_poll_api"], "{uuid}", this.Uuid, 1)
	loginPollApi = strings.Replace(loginPollApi, "{tip}", strconv.Itoa(tip), 1)
	loginPollApi = strings.Replace(loginPollApi, "{time}", GetUnixTime(), 1)
	this.HttpClient.HttpHeader.Host = "login.weixin.qq.com"
	content, _, err := this.HttpClient.Get(loginPollApi)
	if err != nil {
		return ``, err
	}

	regRedirectUri, err := regexp.Compile(`window.redirect_uri="(.+)";`)
	if err != nil {
		return ``, err
	}

	redirectUri := regRedirectUri.FindSubmatch([]byte(content))
	if len(redirectUri) == 2 {
		return string(redirectUri[1]), nil
	}

	if content == "window.code=201;" {
		return "201", nil
	}

	return ``, nil
}

/**
 * 访问redirectUrl以获取登录必须的cookie
 */
func (this *Wechat) WaitForLogin() (string, error) {
	this.HttpClient.HttpHeader.Accept = "application/json, text/plain, */*"
	this.HttpClient.HttpHeader.Referer = "https://wx2.qq.com/?&lang=zh_CN"
	this.HttpClient.HttpHeader.Host = "wx2.qq.com"
	content, cookies, err := this.HttpClient.Get(this.RedirectUrl + "&fun=new&version=v2&lang=zh_CN")
	this.Cookies = cookies
	return content, err
}

/**
 * 解析登陆返回的xml
 */
func (this *Wechat) AnalysisLoginXml(xmlStr string) (Error, error) {
	var v Error
	err := xml.Unmarshal([]byte(xmlStr), &v)
	if err != nil {
		return v, err
	}
	var max int64 = 999999999999999
	var min int64 = 100000000000000
	this.DeviceID = "e" + strconv.Itoa(int(rand.Int63n(max-min)+min))
	this.Sid = v.Wxsid
	this.Uin = v.Wxuin
	this.Skey = v.Skey
	this.PassTicket = v.PassTicket
	return v, nil
}


func (this *Wechat) WxInit() {
	this.HttpClient.HttpHeader.ContentType = "application/json;charset=UTF-8"
	this.HttpClient.HttpHeader.ContentLength = "101"
	this.HttpClient.HttpHeader.Cookie = Cookies2String(this.Cookies)
	wxInitApi := strings.Replace(Config["wx_init_api"], "{r}", strconv.Itoa(int(time.Now().Unix())), 1)
	var postData string = `{BaseRequest: {Uin: "`+ this.Uin + `", Sid: "`+ this.Sid +`", Skey: "`+ this.Skey +`", DeviceID: "`+ this.DeviceID +`"}}`
	content, _, _ := this.HttpClient.PostStr(wxInitApi, postData)
	fmt.Println("content:" + content)
}

func (this *Wechat) StatusNotify() {
	this.HttpClient.HttpHeader.ContentType = "application/json;charset=UTF-8"
	this.HttpClient.HttpHeader.ContentLength = "348"
	this.HttpClient.HttpHeader.Cookie = Cookies2String(this.Cookies)
	wxInitApi := strings.Replace(Config["wx_statusnotify_api"], "{pass_ticket}", this.PassTicket, 1)
	var postData string = `{BaseRequest: {Uin: "`+ this.Uin + `", Sid: "`+ this.Sid +`", Skey: "`+ this.Skey +`", DeviceID: "`+ this.DeviceID +`"}}`
	content, _, _ := this.HttpClient.PostStr(wxInitApi, postData)
	fmt.Println("content:" + content)
}

func (this *Wechat) MakeContactList() {

}

func (this *Wechat) GetBatchGroupMembers() {

}

func (this *Wechat) Sync() {

}

func (this *Wechat) SyncCheck() {

}

func (this *Wechat) SendMsg() {

}