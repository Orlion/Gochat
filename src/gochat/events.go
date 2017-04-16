package gochat

import (
	"strings"
	"fmt"
	"time"
)

type EventType int

const (
	_ EventType = iota
	NewMessageEvent
)

type Event struct {
	Type EventType
	Time int64
	Data interface{}
}

type EventMsgData struct{
	IsGroupMsg       bool
	IsMediaMsg       bool
	IsSendByMySelf   bool
	MsgType          int64
	AtMe             bool
	MediaUrl         string
	Content          string
	FromUserName     string
	FromGGID         string
	SenderUserName   string
	SenderGGID       string
	ToUserName       string
	ToGGID           string
	OriginalMsg      map[string]interface{}
}

type Listener interface {
	handle(event Event) error
}

func (this *Wechat) handleSyncResponse(resp *syncMessageResponse) {
	fmt.Println(resp.AddMsgCount)
	//if resp.AddMsgCount > 0 {
	//	for _, v := range resp.AddMsgList {
	//		go this.emitNewMessageEvent(v)
	//	}
	//}
}

func (this *Wechat) emitNewMessageEvent(msg map[string]interface{}) {
	fromUserName := msg["FromUserName"].(string)
	toUserName := msg["ToUserName"].(string)
	senderUserName := fromUserName
	content := msg["Content"].(string)
	isSendByMySelf := fromUserName == this.me.UserName
	var groupUserName string
	if strings.HasPrefix(fromUserName, "@@") {
		groupUserName = fromUserName
	} else if strings.HasPrefix(toUserName, "@@") {
		groupUserName = toUserName
	}

	isGroupMsg := false
	if len(groupUserName) > 0 {
		isGroupMsg = true
	}
	msgType := msg["MsgType"].(float64)
	mid := msg["MsgId"].(string)

	isMediaMsg := false
	mediaUrl := ""
	path := ""
	switch msgType {
	case 3: {
		// 图片
		path = "webwxgetmsgimg"
	}
	case 47: {
		// pid 图片
		pid, _ := msg["HasProductId"].(float64)
		if pid == 0 {
			path = "webwxgetmsgimg"
		}
	}
	case 34: {
		// 语音
		path = "webwxgetvoice"
	}
	case 43: {
		// 视频
		path = "webwxgetvideo"
	}
	}
	if len(path) > 0 {
		isMediaMsg = true
		mediaUrl = fmt.Sprintf(`https://wx2.qq.com/%s?msgid=%v&%v`, path, mid, this.skeyKV())
	}
	isAtMe := false
	//if isGroupMsg && !isSendByMySelf {
	//	atme := "@"
	//	if len(this.Me.DisplayName) > 0 {
	//		atme += this.Me.DisplayName
	//	} else {
	//		atme += this.Me.NickName
	//	}
	//	isAtMe = strings.Contains(content, atme)
	//	infos := strings.Split(content, atme)
	//	if len(infos) != 2{
	//		return
	//	}
	//	contact, err := this.ContactByUserName(infos[0])
	//	if err != nil {
	//		this.ForceUpateGroup(groupUserName)
	//		return
	//	}
	//
	//	senderUserName = contact.UserName
	//	content = infos[1]
	//}

	data := EventMsgData {
		IsGroupMsg:		isGroupMsg,
		IsMediaMsg:		isMediaMsg,
		IsSendByMySelf:	isSendByMySelf,
		MsgType:		int64(msgType),
		AtMe:			isAtMe,
		MediaUrl:		mediaUrl,
		Content:		content,
		FromUserName:	fromUserName,
		FromGGID:		"",
		SenderUserName:	senderUserName,
		SenderGGID:		"",
		ToUserName:		toUserName,
		ToGGID:			"",
		OriginalMsg:	msg,
	}

	event := Event {
		Type:	NewMessageEvent,
		Time:	time.Now().Unix(),
		Data:	data,
	}

	this.listener.handle(event)
}