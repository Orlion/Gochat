package gochat

import (
	"time"
	"errors"
	"strings"
	"gochat/utils"
	"encoding/json"
)

type initRequest struct {
	BaseRequest baseRequest
}

type initResp struct {
	Response
	User    	Contact
	Skey    	string
	SyncKey 	syncKey
}

// init
func (weChat *WeChat) init() error {
	wxInitApi := strings.Replace(weChatApi["initApi"], "{r}", utils.GetUnixTime(), 1)
	wxInitApi = strings.Replace(wxInitApi, "{host}", weChat.host, 1)
	wxInitApi = strings.Replace(wxInitApi, "{pass_ticket}", weChat.passTicket, 1)

	postData, err := json.Marshal(initRequest{
		BaseRequest: weChat.baseRequest,
	})
	if err != nil {
		return err
	}

	content, err := weChat.httpClient.post(wxInitApi, postData, time.Second * 5, &httpHeader{
		Accept:				"application/json, text/plain, */*",
		ContentType:		"application/json;charset=UTF-8",
		Origin:				"https://" + weChat.host,
		Host: 				weChat.host,
		Referer: 			"https://"+ weChat.host +"/?&lang=zh_CN",
	})
	if err != nil {
		return err
	}

	var initRes initResp
	err = json.Unmarshal([]byte(content), &initRes)
	if err != nil {
		return err
	}

	if initRes.Response.BaseResponse.Ret != 0 {
		weChat.logger.Println("[Error] Init Failed. Res.Ret=" + string(initRes.Response.BaseResponse.Ret))
		return errors.New("Init Failed")
	}

	weChat.me = initRes.User
	weChat.baseRequest.Skey = initRes.Skey
	weChat.syncKey = initRes.SyncKey

	return nil
}