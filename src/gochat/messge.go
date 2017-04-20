package gochat

import (
	"strings"
	"time"
	"strconv"
	"math/rand"
	"bytes"
	"encoding/json"
	"errors"
)

type sendMsgResponse struct {
	Response
	MsgID   string
	LocalID string
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
		"Type":			1,
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
