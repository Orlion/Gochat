package gochat

import (
	"time"
	"strconv"
	"regexp"
	"errors"
	"strings"
	"gochat/utils"
	"fmt"
	"math/rand"
)

func (weChat *WeChat) login() error {

	var err error
	weChat.Uuid, weChat.baseRequest, weChat.passTicket, weChat.httpClient.Cookies, weChat.host, err = weChat.storage.getData()
	if err != nil {
		weChat.Uuid, err = weChat.getUuid()
		if err != nil {
			return err
		}

		go weChat.triggerGenUuidEvent(weChat.Uuid)
		weChat.logger.Println("[Info] Uuid=" + weChat.Uuid)

		tip := 1
		redirectUrl := ""
		for  {
			status, result, err := weChat.isAuth(tip)
			if err != nil {
				weChat.logger.Println("[Error] GetRedirectUrl Error :" + err.Error())
				time.Sleep(time.Second * time.Duration(1))
				continue
			}

			if 200 == status {
				redirectUrl = result
				weChat.logger.Println("[Info] Redirect=" + redirectUrl)
				go weChat.triggerConfirmAuthEvent(redirectUrl)
				break
			}

			if 201 == status {
				tip = 0
				weChat.logger.Println("[Info] Scan Code")
				go weChat.triggerScanCodeEvent(result)
			}
		}

		weChat.host = utils.GetHostByUrl(redirectUrl)
		err = weChat.doLogin(redirectUrl)
		if err != nil {
			return err
		}

		weChat.storage.setData(weChat.Uuid, weChat.baseRequest, weChat.passTicket, weChat.httpClient.Cookies, weChat.host)
	}

	weChat.logger.Println("[Info] Login: " + fmt.Sprintf("Sid=[ %s ], Uin=[ %s ], Skey=[ %s ], PassTicket=[ %s ]", weChat.baseRequest.Sid, weChat.baseRequest.Uin, weChat.baseRequest.Skey, weChat.passTicket,))
	go weChat.triggerLoginEvent(weChat.baseRequest.DeviceID)

	return nil
}

// 获取Uuid
func (weChat *WeChat) getUuid() (string, error) {

	getUuidApiUrl := weChatApi["getUuidApi"] + utils.GetUnixMsTime()
	content, err := weChat.httpClient.get(getUuidApiUrl, time.Second * 5, &httpHeader{
		Host: 				"login.wx2.qq.com",
		Referer: 			"https://wx2.qq.com/?&lang=zh_CN",
	})
	if err != nil {
		return "", err
	}

	reg, err := regexp.Compile(`window.QRLogin.code = 200; window.QRLogin.uuid = "(.+)"`)
	if err != nil {
		return "", err
	}

	uuidArr := reg.FindSubmatch([]byte(content))
	if len(uuidArr) != 2 {
		return "", errors.New("Uuid get failed")
	}

	return string(uuidArr[1]), nil
}

// 判断是否已授权登陆,获取redirectUrl
func (weChat *WeChat) isAuth(tip int) (int, string, error) {

	loginPollApi := strings.Replace(weChatApi["loginApi"], "{uuid}", weChat.Uuid, 1)
	loginPollApi = strings.Replace(loginPollApi, "{tip}", strconv.Itoa(tip), 1)
	loginPollApi = strings.Replace(loginPollApi, "{time}", utils.GetUnixMsTime(), 1)

	content, err := weChat.httpClient.get(loginPollApi, time.Second * 30, &httpHeader{
		Host: 				"login.wx2.qq.com",
		Referer: 			"https://wx2.qq.com/?&lang=zh_CN",
	})
	if err != nil {
		return 0, "", err
	}

	regRedirectUri, err := regexp.Compile(`window.redirect_uri="(.+)";`)
	if err != nil {
		return 0, "", err
	}

	redirectUriArr := regRedirectUri.FindSubmatch([]byte(content))
	if len(redirectUriArr) == 2 {
		return 200, string(redirectUriArr[1]), nil
	}

	regScanCode, err := regexp.Compile(`window.code=201;window.userAvatar = '(.+)';`)
	if err != nil {
		return 0, "", err
	}

	userAvatarArr := regScanCode.FindSubmatch([]byte(content))
	if len(userAvatarArr) == 2 {
		return 201, string(userAvatarArr[1]), nil
	}

	return 0, "", nil
}

// 请求redirectUrl 登录
func (weChat *WeChat) doLogin(redirectUrl string) error {
	content, err := weChat.httpClient.get(redirectUrl + "&fun=new&version=v2&lang=zh_CN", time.Second * 5, &httpHeader{
		Host: 				weChat.host,
		Referer: 			"https://"+ weChat.host +"/?&lang=zh_CN",
	})
	if err != nil {
		return err
	}

	var max int64 = 999999999999999
	var min int64 = 100000000000000
	weChat.baseRequest.DeviceID = "e" + strconv.Itoa(int(rand.Int63n(max-min) + min))
	weChat.baseRequest.Sid, weChat.baseRequest.Uin, weChat.baseRequest.Skey, weChat.passTicket, err = utils.AnalysisLoginXml(content)

	return err
}
