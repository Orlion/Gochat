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

func (this *Wechat) SendTextMsg(content string, to string) (bool, error) {
	sendMsgApi := strings.Replace(Config["sendmsg_api"], "{pass_ticket}", this.passTicket, 1)
	sendMsgApi = strings.Replace(sendMsgApi, "{host}", this.host, 1)
	msgId := this.utils.getUnixMsTime() + strconv.Itoa(rand.Intn(10000))
	msg := map[string]interface{} {
		"Content":		content,
		"ToUserName":	to,
		"FromUserName": this.me.UserName,
		"LocalID":		msgId,
		"ClientMsgId":	msgId,
		"Type":			"1",
	}

	buffer := new(bytes.Buffer)
	enc := json.NewEncoder(buffer)
	enc.SetEscapeHTML(false)
	err := enc.Encode(map[string]interface{}{
		`BaseRequest`: this.baseRequest,
		`Msg`:         msg,
		`Scene`:       0,
	})

	if err != nil {
		return false, err
	}

	respContent, err := this.httpClient.post(sendMsgApi, []byte(buffer.String()), time.Second * 5, &HttpHeader{
		Accept: 			"application/json, text/plain, */*",
		AcceptEncoding: 	"gzip, deflate, br",
		AcceptLanguage: 	"zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3",
		Connection: 		"keep-alive",
		ContentType:		"application/json;charset=utf-8",
		Host: 				this.host,
		Referer: 			"https://"+ this.host +"/?&lang=zh_CN",
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

func (this *Wechat) VerifyUser(userName string, ticket string) error {
	verifyUserApi := strings.Replace(Config["verify_user_api"], "{pass_ticket}", this.passTicket, 1)
	verifyUserApi = strings.Replace(verifyUserApi, "{host}", this.host, 1)
	verifyUserApi = strings.Replace(verifyUserApi, "{r}", this.utils.getUnixMsTime(), 1)

	buffer := new(bytes.Buffer)
	enc := json.NewEncoder(buffer)
	enc.SetEscapeHTML(false)
	err := enc.Encode(map[string]interface{}{
		`BaseRequest`: 			this.baseRequest,
		`Opcode`:				3,
		`SceneList`:			[]int{33},
		`SceneListCount`:		1,
		`VerifyContent`:		``,
		`VerifyUserList`:		[]map[string]string{{"Value": userName,"VerifyUserTicket": ticket}},
		`VerifyUserListSize`:	1,
		`skey`:       			this.baseRequest.Skey,
	})

	if err != nil {
		return err
	}

	respContent, err := this.httpClient.post(verifyUserApi, []byte(buffer.String()), time.Second * 5, &HttpHeader{
		Accept: 			"application/json, text/plain, */*",
		AcceptEncoding: 	"gzip, deflate, br",
		AcceptLanguage: 	"zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3",
		Connection: 		"keep-alive",
		ContentType:		"application/json;charset=utf-8",
		Host: 				this.host,
		Referer: 			"https://"+ this.host +"/?&lang=zh_CN",
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

func (this *Wechat) UploadMedia() {
	prefix := ""
	verifyUserApi := strings.Replace(Config["upload_media_api"], "{host}", this.host, 1)
	verifyUserApi = strings.Replace(verifyUserApi, "{prefix}", prefix, 1)

	respContent, err := this.httpClient.post(verifyUserApi, time.Second * 5, &HttpHeader{
		Accept: 			"*/*",
		AcceptEncoding: 	"gzip, deflate, br",
		AcceptLanguage: 	"zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3",
		Connection: 		"keep-alive",
		ContentType:		"multipart/form-data; boundary=----WebKitFormBoundary5kzhsa5PtvvA8b49",
		Host: 				prefix + "." + this.host,
		Referer: 			"https://"+ this.host +"/?&lang=zh_CN",
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

func (this *Wechat) sendAppMsg() {

}

func (this *Wechat) SendImgMsg(to string, mediaId string) error {
	sendImgMsgApi := strings.Replace(Config["sendimgmsg_api"], "{host}", this.host, 1)
	msgId := this.utils.getUnixMsTime() + strconv.Itoa(rand.Intn(10000))
	msg := map[string]interface{} {
		"Content":		"",
		"ToUserName":	to,
		"FromUserName": this.me.UserName,
		"LocalID":		msgId,
		"MediaId":		mediaId,
		"Type":			"3",
	}

	buffer := new(bytes.Buffer)
	enc := json.NewEncoder(buffer)
	enc.SetEscapeHTML(false)
	err := enc.Encode(map[string]interface{}{
		`BaseRequest`: this.baseRequest,
		`Msg`:         msg,
		`Scene`:       0,
	})

	if err != nil {
		return false, err
	}

	respContent, err := this.httpClient.post(sendImgMsgApi, []byte(buffer.String()), time.Second * 5, &HttpHeader{
		Accept: 			"application/json, text/plain, */*",
		AcceptEncoding: 	"gzip, deflate, br",
		AcceptLanguage: 	"zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3",
		Connection: 		"keep-alive",
		ContentType:		"application/json;charset=utf-8",
		Host: 				this.host,
		Referer: 			"https://"+ this.host +"/?&lang=zh_CN",
	})

	var resp sendMsgResponse
	err = json.Unmarshal([]byte(respContent), &resp)
	if err != nil {
		return false, err
	}

	if (resp.BaseResponse.Ret != 0) {
		return false, errors.New("Send Img Msg Error. [msgId]:" + msgId)
	}

	return true, nil
}

func (this *Wechat) SendAppMsg(to string, filename string, mediaId string, ext string) error {
	sendAppMsgApi := strings.Replace(Config["sendappmsg_api"], "{host}", this.host, 1)
	msgId := this.utils.getUnixMsTime() + strconv.Itoa(rand.Intn(10000))
	content := fmt.Sprintf("<appmsg appid='wxeb7ec651dd0aefa9' sdkver=''><title>%s</title><des></des><action></action><type>6</type><content></content><url></url><lowurl></lowurl><appattach><totallen>10</totallen><attachid>%s</attachid><fileext>%s</fileext></appattach><extinfo></extinfo></appmsg>", filename, mediaId, ext)

	msg := map[string]interface{} {
		"ClientMsgId":	msgId,
		"Content":		content,
		"FromUserName": this.me.UserName,
		"LocalID":		msgId,
		"ToUserName":	to,
		"Type":			"6",
	}

	buffer := new(bytes.Buffer)
	enc := json.NewEncoder(buffer)
	enc.SetEscapeHTML(false)
	err := enc.Encode(map[string]interface{}{
		`BaseRequest`: this.baseRequest,
		`Msg`:         msg,
		`Scene`:       0,
	})

	if err != nil {
		return err
	}

	respContent, err := this.httpClient.post(sendAppMsgApi, []byte(buffer.String()), time.Second * 5, &HttpHeader{
		Accept: 			"application/json, text/plain, */*",
		AcceptEncoding: 	"gzip, deflate, br",
		AcceptLanguage: 	"zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3",
		Connection: 		"keep-alive",
		ContentType:		"application/json;charset=utf-8",
		Host: 				this.host,
		Referer: 			"https://"+ this.host +"/?&lang=zh_CN",
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