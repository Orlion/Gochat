package gochat

import (
	"regexp"
	"time"
	"fmt"
	"strings"
	"encoding/xml"
	"strconv"
	"math/rand"
	"encoding/json"
)

type  Wechat struct {
	BaseRequest BaseRequest
	PassTicket 	string
	syncKey   	map[string]interface{}
	syncHost	string
	Uuid 		string
	Me			Contact
	Utils		Utils
	HttpClient	*HttpClient
}

type BaseRequest struct {
	Sid      	string
	Skey       	string
	Uin      	string
	DeviceID	string
}

type Response struct {
	BaseResponse *BaseResponse
}

type BaseResponse struct {
	Ret    int
	ErrMsg string
}

/**
 * 初始化
 */
func New() *Wechat{
	return & Wechat{
		Utils: Utils{},
		HttpClient: &HttpClient{
			HttpHeader: HttpHeader{
				"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
				"gzip, deflate, sdch, br",
				"zh-CN,zh;q=0.8",
				"keep-alive",
				"",
				"",
				nil,
				"login.wx2.qq.com",
				"https://wx.qq.com/",
				"1",
			},
		},
	}
}

/**
 * 获取uuid
 */
func (this *Wechat) Begin() (string, error) {
	getUuidApiUrl := Config["getUuidApi"] + this.Utils.getUnixMsTime()
	content, _, err := this.HttpClient.get(getUuidApiUrl, time.Second * 5)
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

func (this *Wechat) Login() error {
	var tip int = 1
	for  {
		redirectUrl, err := this.polling(tip)
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second * time.Duration(1))
			continue
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

		err = this.doLogin(redirectUrl)
		if err != nil {
			return err
		}

		this.init()

		return nil
	}
}

/**
 * 轮询,直到用户在手机微信上确认登录，确认登录后会返回redirectUrl
 */
func (this *Wechat) polling(tip int) (string, error){
	loginPollApi := strings.Replace(Config["login_poll_api"], "{uuid}", this.Uuid, 1)
	loginPollApi = strings.Replace(loginPollApi, "{tip}", strconv.Itoa(tip), 1)
	loginPollApi = strings.Replace(loginPollApi, "{time}", this.Utils.getUnixMsTime(), 1)

	this.HttpClient.HttpHeader.Accept = "*/*"
	this.HttpClient.HttpHeader.Host = "login.wx2.qq.com"
	this.HttpClient.HttpHeader.Referer = "https://wx2.qq.com/?&lang=zh_CN"
	content, _, err := this.HttpClient.get(loginPollApi, time.Second * 30)
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
func (this *Wechat) doLogin(redirectUrl string) error {
	this.HttpClient.HttpHeader.Accept = "application/json, text/plain, */*"
	this.HttpClient.HttpHeader.Host = "wx2.qq.com"
	content, cookie, err := this.HttpClient.get(redirectUrl + "&fun=new&version=v2&lang=zh_CN", time.Second * 5)
	if err != nil {
		return err
	}
	this.HttpClient.HttpHeader.Cookies = cookie
	this.BaseRequest, this.PassTicket, err = this.analysisLoginXml(content)

	return err
}

/**
 * 解析登陆返回的xml
 */
func (this *Wechat) analysisLoginXml(xmlStr string) (BaseRequest, string, error) {
	type Error struct {
		Ret string  `xml:"ret"`
		Message string  `xml:"message"`
		Skey string  `xml:"skey"`
		Wxsid string  `xml:"wxsid"`
		Wxuin string  `xml:"wxuin"`
		PassTicket string  `xml:"pass_ticket"`
		Isgrayscale string `xml:"isgrayscale"`
	}

	var v Error
	err := xml.Unmarshal([]byte(xmlStr), &v)
	if err != nil {
		return BaseRequest{}, "", err
	}
	var max int64 = 999999999999999
	var min int64 = 100000000000000
	baseRequest := BaseRequest{
		DeviceID: "e" + strconv.Itoa(int(rand.Int63n(max-min)+min)),
		Sid: v.Wxsid,
		Uin: v.Wxuin,
		Skey: v.Skey,
	}

	return baseRequest, v.PassTicket, nil
}

func (this *Wechat) init() error {
	wxInitApi := strings.Replace(Config["wx_init_api"], "{r}", strconv.Itoa(int(time.Now().Unix())), 1)
	type initRequest struct {
		BaseRequest BaseRequest
	}
	postData, err := json.Marshal(initRequest{
		BaseRequest: this.BaseRequest,
	})
	if err != nil {
		return err
	}
	cookies := this.HttpClient.HttpHeader.Cookies
	this.HttpClient.HttpHeader = HttpHeader{
		"",
		"",
		"",
		"",
		"",
		"",
		nil,
		"",
		"",
		"",
	}
	this.HttpClient.HttpHeader.Cookies = cookies
	// WTF!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	content, _, err := this.HttpClient.Post(wxInitApi, postData, time.Second * 5)
	if err != nil {
		return err
	}
	type initResp struct {
		Response
		User    Contact
		Skey    string
		SyncKey map[string]interface{}
	}
	var initres initResp
	err = json.Unmarshal([]byte(content), &initres)
	this.Me = initres.User
	this.BaseRequest.Skey = initres.Skey
	this.syncKey = initres.SyncKey
	return nil
}
