package gochat

import (
	"strings"
	"time"
	"strconv"
	"math/rand"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gochat/utils"
)

type sendMsgResponse struct {
	Response
	MsgID   string
	LocalID string
}

type verifyUserResponse struct {
	Response
}

type uploadMediaResponse struct {
	Response
	MediaId string
	StartPos string
	CDNThumbImgHeight string
	CDNThumbImgWidth string
}

func (WeChat *WeChat) SendTextMsg(content string, to string) (bool, error) {
	sendMsgApi := strings.Replace(weChatApi["sendMsgApi"], "{pass_ticket}", WeChat.passTicket, 1)
	sendMsgApi = strings.Replace(sendMsgApi, "{host}", WeChat.host, 1)
	msgId := utils.GetUnixMsTime() + strconv.Itoa(rand.Intn(10000))
	msg := map[string]interface{} {
		"Content":		content,
		"ToUserName":	to,
		"FromUserName": WeChat.me.UserName,
		"LocalID":		msgId,
		"ClientMsgId":	msgId,
		"Type":			"1",
	}

	buffer := new(bytes.Buffer)
	enc := json.NewEncoder(buffer)
	enc.SetEscapeHTML(false)
	err := enc.Encode(map[string]interface{}{
		`BaseRequest`: WeChat.baseRequest,
		`Msg`:         msg,
		`Scene`:       0,
	})

	if err != nil {
		return false, err
	}

	respContent, err := WeChat.httpClient.post(sendMsgApi, []byte(buffer.String()), time.Second * 5, &httpHeader{
		Host: 				WeChat.host,
		Referer: 			"https://"+ WeChat.host +"/?&lang=zh_CN",
	})

	var resp sendMsgResponse
	err = json.Unmarshal([]byte(respContent), &resp)
	if err != nil {
		return false, err
	}

	if (resp.BaseResponse.Ret != 0) {
		return false, errors.New("Send Msg Error. [msgId]:" + msgId)
	}

	return true, nil
}

func (WeChat *WeChat) VerifyUser(userName string, ticket string, verifyUserContent string) error {
	verifyUserApi := strings.Replace(weChatApi["verifyUserApi"], "{pass_ticket}", WeChat.passTicket, 1)
	verifyUserApi = strings.Replace(verifyUserApi, "{host}", WeChat.host, 1)
	verifyUserApi = strings.Replace(verifyUserApi, "{r}", utils.GetUnixMsTime(), 1)

	buffer := new(bytes.Buffer)
	enc := json.NewEncoder(buffer)
	enc.SetEscapeHTML(false)
	err := enc.Encode(map[string]interface{}{
		"BaseRequest": 			WeChat.baseRequest,
		"Opcode":				3,
		"SceneList":			[]int{33},
		"SceneListCount":		1,
		"VerifyContent":		verifyUserContent,
		"VerifyUserList":		[]map[string]string{{"Value": userName,"VerifyUserTicket": ticket}},
		"VerifyUserListSize":	1,
		"skey":       			WeChat.baseRequest.Skey,
	})

	if err != nil {
		return err
	}

	respContent, err := WeChat.httpClient.post(verifyUserApi, []byte(buffer.String()), time.Second * 5, &httpHeader{
		Host: 				WeChat.host,
		Referer: 			"https://"+ WeChat.host +"/?&lang=zh_CN",
	})

	fmt.Println(respContent)
	var resp verifyUserResponse
	err = json.Unmarshal([]byte(respContent), &resp)
	if err != nil {
		return err
	}

	if (resp.BaseResponse.Ret != 0) {
		return errors.New("VerifyUser Error")
	}

	return nil
}

func (WeChat *WeChat) UploadMedia() error {
	prefix := ""
	verifyUserApi := strings.Replace(weChatApi["uploadMediaApi"], "{host}", WeChat.host, 1)
	verifyUserApi = strings.Replace(verifyUserApi, "{prefix}", prefix, 1)

	respContent, err := WeChat.httpClient.post(verifyUserApi, []byte(""), time.Second * 5, &httpHeader{
		Accept: 			"*/*",
		AcceptEncoding: 	"gzip, deflate, br",
		AcceptLanguage: 	"zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3",
		Connection: 		"keep-alive",
		ContentType:		"multipart/form-data; boundary=----WebKitFormBoundary5kzhsa5PtvvA8b49",
		Host: 				prefix + "." + WeChat.host,
		Referer: 			"https://"+ WeChat.host +"/?&lang=zh_CN",
	})

	fmt.Println(respContent)
	var resp verifyUserResponse
	err = json.Unmarshal([]byte(respContent), &resp)
	if err != nil {
		return err
	}

	if (resp.BaseResponse.Ret != 0) {
		return errors.New("UploadMedia Error")
	}

	return nil
}

func (WeChat *WeChat) sendAppMsg() {

}

func (WeChat *WeChat) SendImgMsg(to string, mediaId string) error {
	sendImgMsgApi := strings.Replace(weChatApi["sendImgMsgApi"], "{host}", WeChat.host, 1)
	msgId := utils.GetUnixMsTime() + strconv.Itoa(rand.Intn(10000))
	msg := map[string]interface{} {
		"Content":		"",
		"ToUserName":	to,
		"FromUserName": WeChat.me.UserName,
		"LocalID":		msgId,
		"MediaId":		mediaId,
		"Type":			"3",
	}

	buffer := new(bytes.Buffer)
	enc := json.NewEncoder(buffer)
	enc.SetEscapeHTML(false)
	err := enc.Encode(map[string]interface{}{
		`BaseRequest`: WeChat.baseRequest,
		`Msg`:         msg,
		`Scene`:       0,
	})

	if err != nil {
		return err
	}

	respContent, err := WeChat.httpClient.post(sendImgMsgApi, []byte(buffer.String()), time.Second * 5, &httpHeader{
		Accept: 			"application/json, text/plain, */*",
		AcceptEncoding: 	"gzip, deflate, br",
		AcceptLanguage: 	"zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3",
		Connection: 		"keep-alive",
		ContentType:		"application/json;charset=utf-8",
		Host: 				WeChat.host,
		Referer: 			"https://"+ WeChat.host +"/?&lang=zh_CN",
	})

	var resp sendMsgResponse
	err = json.Unmarshal([]byte(respContent), &resp)
	if err != nil {
		return err
	}

	if (resp.BaseResponse.Ret != 0) {
		return errors.New("Send Img Msg Error. [msgId]:" + msgId)
	}

	return nil
}

// 发送
func (WeChat *WeChat) SendAppMsg(to string, filename string, mediaId string, ext string) error {
	sendAppMsgApi := strings.Replace(weChatApi["sendAppMsgApi"], "{host}", WeChat.host, 1)
	msgId := utils.GetUnixMsTime() + strconv.Itoa(rand.Intn(10000))
	content := fmt.Sprintf("<appmsg appid='wxeb7ec651dd0aefa9' sdkver=''><title>%s</title><des></des><action></action><type>6</type><content></content><url></url><lowurl></lowurl><appattach><totallen>10</totallen><attachid>%s</attachid><fileext>%s</fileext></appattach><extinfo></extinfo></appmsg>", filename, mediaId, ext)

	msg := map[string]interface{} {
		"ClientMsgId":	msgId,
		"Content":		content,
		"FromUserName": WeChat.me.UserName,
		"LocalID":		msgId,
		"ToUserName":	to,
		"Type":			"6",
	}

	buffer := new(bytes.Buffer)
	enc := json.NewEncoder(buffer)
	enc.SetEscapeHTML(false)
	err := enc.Encode(map[string]interface{}{
		`BaseRequest`: WeChat.baseRequest,
		`Msg`:         msg,
		`Scene`:       0,
	})

	if err != nil {
		return err
	}

	respContent, err := WeChat.httpClient.post(sendAppMsgApi, []byte(buffer.String()), time.Second * 5, &httpHeader{
		Accept: 			"application/json, text/plain, */*",
		AcceptEncoding: 	"gzip, deflate, br",
		AcceptLanguage: 	"zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3",
		Connection: 		"keep-alive",
		ContentType:		"application/json;charset=utf-8",
		Host: 				WeChat.host,
		Referer: 			"https://"+ WeChat.host +"/?&lang=zh_CN",
	})

	var resp sendMsgResponse
	err = json.Unmarshal([]byte(respContent), &resp)
	if err != nil {
		return err
	}

	if (resp.BaseResponse.Ret != 0) {
		return errors.New("Send App Msg Error. [msgId]:" + msgId)
	}

	return nil
}