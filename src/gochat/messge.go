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
	"os"
	"crypto/md5"
	"mime/multipart"
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
}

type MediaType int
const (
	_ MediaType = iota
	MEDIA_PIC
	MEDIA_VIDEO
	MEDIA_DOC
)

var mediaIndex int64 = 0

func (weChat *WeChat) SendTextMsg(content string, to string) (bool, error) {
	sendMsgApi := strings.Replace(weChatApi["sendMsgApi"], "{pass_ticket}", weChat.passTicket, 1)
	sendMsgApi = strings.Replace(sendMsgApi, "{host}", weChat.host, 1)
	msgId := utils.GetUnixMsTime() + strconv.Itoa(rand.Intn(10000))
	msg := map[string]interface{} {
		"Content":		content,
		"ToUserName":	to,
		"FromUserName": weChat.me.UserName,
		"LocalID":		msgId,
		"ClientMsgId":	msgId,
		"Type":			"1",
	}

	buffer := new(bytes.Buffer)
	enc := json.NewEncoder(buffer)
	// enc.SetEscapeHTML(false)
	err := enc.Encode(map[string]interface{}{
		`BaseRequest`: weChat.baseRequest,
		`Msg`:         msg,
		`Scene`:       0,
	})

	if err != nil {
		return false, err
	}

	respContent, err := weChat.httpClient.post(sendMsgApi, []byte(buffer.String()), time.Second * 5, &httpHeader{
		Host: 				weChat.host,
		Referer: 			"https://"+ weChat.host +"/?&lang=zh_CN",
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


// 发送图片消息
func (weChat *WeChat) SendImgMsg(toUserFrom string, mediaId string) error {
	sendImgMsgApi := strings.Replace(weChatApi["sendImgMsgApi"], "{host}", weChat.host, 1)
	sendImgMsgApi = strings.Replace(sendImgMsgApi, "{pass_ticket}", weChat.passTicket, 1)
	msgId := utils.GetUnixMsTime() + strconv.Itoa(rand.Intn(10000))
	msg := map[string]interface{} {
		"Content":		"",
		"ToUserName":	toUserFrom,
		"FromUserName": weChat.me.UserName,
		"LocalID":		msgId,
		"MediaId":		mediaId,
		"Type":			"3",
	}

	buffer := new(bytes.Buffer)
	enc := json.NewEncoder(buffer)
	// enc.SetEscapeHTML(false)
	err := enc.Encode(map[string]interface{}{
		"BaseRequest": weChat.baseRequest,
		"Msg":         msg,
		"Scene":       0,
	})

	if err != nil {
		return err
	}

	respContent, err := weChat.httpClient.post(sendImgMsgApi, []byte(buffer.String()), time.Second * 5, &httpHeader{
		Host: 				weChat.host,
		Referer: 			"https://"+ weChat.host +"/?&lang=zh_CN",
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

// 发送文件消息
func (weChat *WeChat) SendAppMsg(toUserName string, mediaId string, filename string, fileSize int64, ext string) error {
	sendAppMsgApi := strings.Replace(weChatApi["sendAppMsgApi"], "{host}", weChat.host, 1)
	sendAppMsgApi = strings.Replace(sendAppMsgApi, "{pass_ticket}", weChat.passTicket, 1)
	msgId := utils.GetUnixMsTime() + strconv.Itoa(rand.Intn(10000))
	content := fmt.Sprintf("<appmsg appid='wxeb7ec651dd0aefa9' sdkver=''><title>%s</title><des></des><action></action><type>6</type><content></content><url></url><lowurl></lowurl><appattach><totallen>%d</totallen><attachid>%s</attachid><fileext>%s</fileext></appattach><extinfo></extinfo></appmsg>", filename, fileSize, mediaId, ext)

	msg := map[string]interface{} {
		"ClientMsgId":	msgId,
		"Content":		content,
		"FromUserName": weChat.me.UserName,
		"LocalID":		msgId,
		"ToUserName":	toUserName,
		"Type":			"6",
	}

	buffer := new(bytes.Buffer)
	enc := json.NewEncoder(buffer)
	// enc.SetEscapeHTML(false)
	err := enc.Encode(map[string]interface{}{
		`BaseRequest`: weChat.baseRequest,
		`Msg`:         msg,
		`Scene`:       0,
	})

	if err != nil {
		return err
	}

	respContent, err := weChat.httpClient.post(sendAppMsgApi, []byte(buffer.String()), time.Second * 5, &httpHeader{
		Accept: 			"application/json, text/plain, */*",
		AcceptEncoding: 	"gzip, deflate, br",
		AcceptLanguage: 	"zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3",
		Connection: 		"keep-alive",
		ContentType:		"application/json;charset=utf-8",
		Host: 				weChat.host,
		Referer: 			"https://"+ weChat.host +"/?&lang=zh_CN",
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

// 上传文件方法
func (weChat *WeChat) UploadMedia(buf []byte, mediaType MediaType, fileType string, fileInfo os.FileInfo, toUserName string) (string, error) {

	mediaTypeStr := "doc"
	switch mediaType {
	case MEDIA_PIC:
		mediaTypeStr = "pic"
	case MEDIA_VIDEO:
		mediaTypeStr = "video"
	}

	fields := map[string]string {
		"id":					"WU_FILE_" + string(mediaIndex),
		"name":					fileInfo.Name(),
		"type":					fileType,
		"lastModifiedDate": 	fileInfo.ModTime().UTC().String(),
		"size": 				string(fileInfo.Size()),
		"mediatype": 			mediaTypeStr,
		"pass_ticket": 			weChat.passTicket,
		"webwx_data_ticket": 	weChat.httpClient.getDataTicket(),
	}

	media, err := json.Marshal(&map[string]interface{}{
		"BaseRequest":		weChat.baseRequest,
		"ClientMediaId":	utils.GetUnixMsTime(),
		"TotalLen":			string(fileInfo.Size()),
		"StartPos": 		0,
		"DataLen":			string(fileInfo.Size()),
		"MediaType": 		4,
		"UploadType": 		2,
		"ToUserName": 		toUserName,
		"FromUserName": 	weChat.me.UserName,
		"FileMd5": 			string(md5.New().Sum(buf)),
	})

	if err != nil {
		return "", err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fw, err := writer.CreateFormFile("filename", fileInfo.Name())
	if err != nil {
		return "", err
	}
	fw.Write(buf)

	for k, v := range fields {
		writer.WriteField(k, v)
	}

	writer.WriteField("uploadmediarequest", string(media))
	writer.Close()

	uploadMediaApi := strings.Replace(weChatApi["uploadMediaApi"], "{host}", weChat.host, 1)

	prefixs := []string{"file", "file2"}
	for _, prefix := range prefixs {
		uploadMediaApiDo := strings.Replace(uploadMediaApi, "{prefix}", prefix, 1)
		respContent, err := weChat.httpClient.upload(uploadMediaApiDo, body, time.Second * 5, &httpHeader{
			ContentType:		writer.FormDataContentType(),
			Host: 				prefix + "." + weChat.host,
			Referer: 			"https://"+ weChat.host +"/?&lang=zh_CN",
		})

		var resp uploadMediaResponse
		err = json.Unmarshal([]byte(respContent), &resp)
		if err != nil {
			return "", err
		}

		if (resp.BaseResponse.Ret == 0) {
			return resp.MediaId, nil
		}
	}

	return "", errors.New("UploadMedia Error")
}

// 授权好友请求
func (weChat *WeChat) VerifyUser(userName string, ticket string, verifyUserContent string) error {
	verifyUserApi := strings.Replace(weChatApi["verifyUserApi"], "{pass_ticket}", weChat.passTicket, 1)
	verifyUserApi = strings.Replace(verifyUserApi, "{host}", weChat.host, 1)
	verifyUserApi = strings.Replace(verifyUserApi, "{r}", utils.GetUnixMsTime(), 1)

	buffer := new(bytes.Buffer)
	enc := json.NewEncoder(buffer)
	// enc.SetEscapeHTML(false)
	err := enc.Encode(map[string]interface{}{
		"BaseRequest": 			weChat.baseRequest,
		"Opcode":				3,
		"SceneList":			[]int{33},
		"SceneListCount":		1,
		"VerifyContent":		verifyUserContent,
		"VerifyUserList":		[]map[string]string{{"Value": userName,"VerifyUserTicket": ticket}},
		"VerifyUserListSize":	1,
		"skey":       			weChat.baseRequest.Skey,
	})

	if err != nil {
		return err
	}

	respContent, err := weChat.httpClient.post(verifyUserApi, []byte(buffer.String()), time.Second * 5, &httpHeader{
		Host: 				weChat.host,
		Referer: 			"https://"+ weChat.host +"/?&lang=zh_CN",
	})

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
