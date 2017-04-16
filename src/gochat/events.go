package gochat

import (
	"strings"
	"fmt"
	"time"
)

type EventType int

const (
	_ EventType = iota
	ScanCodeEvent		// 已扫码，未确认
	ConfirmAuthEvent	// 已确认授权登录
	InitEvent			// 初始化完成
	ContactsInitEvent	// 联系人初始化完
	ContactChangeEvent	// 联系人改变了
	MsgEvent
)

type Event struct {
	Type EventType
	Time int64
	Data interface{}
}

type MsgEventData struct{
	IsGroupMsg       bool
	IsMediaMsg       bool
	IsSendByMySelf   bool
	MsgType          int64
	AtMe             bool
	MediaUrl         string
	Content          string
	FromUserName     string
	FromUserInfo     Contact
	SenderUserName   string
	SenderUserInfo   SenderUserInfo
	ToUserName       string
	ToUserInfo       Contact
	OriginalMsg      map[string]interface{}
}

type ContactChangeType int
const (
	_ ContactChangeType = iota
	ContactModify
	ContactDelete
)

type ContactEventData struct {
	ChangeType	ContactChangeType
	UserName	string
}

type SenderUserInfo struct {
	UserName	string
	NickName	string
	ContactType ContactType	// 发送人类型(好友，群成员，好友兼群成员)
}

type Listener interface {
	handle(event Event) error
}

func (this *Wechat) handleSyncResponse(resp *syncMessageResponse) {

	if resp.ModContactCount > 0 {
		for _, v := range resp.ModContactList {
			go this.emitContactChangeEvent(ContactModify, v["UserName"].(string))
		}
	}

	if resp.DelContactCount > 0 {
		for _, v := range resp.DelContactList {
			go this.emitContactChangeEvent(ContactDelete, v["UserName"].(string))
		}
	}

	if resp.AddMsgCount > 0 {
		for _, v := range resp.AddMsgList {
			go this.emitNewMessageEvent(v)
		}
	}
}

func (this *Wechat) emitNewMessageEvent(msg map[string]interface{}) {

	fromUserName := msg["FromUserName"].(string)
	toUserName := msg["ToUserName"].(string)
	senderUserName := fromUserName
	content := msg["Content"].(string)
	isSendByMySelf := false
	senderUserInfo := SenderUserInfo{}

	var groupUserName string
	if strings.HasPrefix(fromUserName, "@@") {	// 消息来自于群
		groupUserName = fromUserName
	} else if strings.HasPrefix(toUserName, "@@") { // 消息来自于群
		groupUserName = toUserName
	}

	isGroupMsg := false // 标识是否是群消息
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
	if isGroupMsg {
		atme := "@"
		if len(this.me.DisplayName) > 0 {
			atme += this.me.DisplayName
		} else {
			atme += this.me.NickName
		}
		isAtMe = strings.Contains(content, atme) // 标识是否是@我

		infos := strings.Split(content, ":<br/>")
		if len(infos) != 2{
			return
		}

		content = infos[1]
		if isAtMe && infos[0] == this.me.UserName {
			isSendByMySelf = true
		}

		contact, found  := this.contacts[fromUserName].MemberMap[infos[0]] // 根据content中UserName(消息发布人)找到详细数据
		if !found {
			// 需要更新群组
			_,err := this.updateOrAddContact([]string{fromUserName})
			if err != nil {
				return
			}

			contact, found = this.contacts[fromUserName].MemberMap[infos[0]]
			if !found {
				return
			}
		}

		senderUserName = infos[0] // 实际发布人
		senderUserInfo = SenderUserInfo{
			UserName: infos[0],
			NickName: contact.NickName,
			ContactType: GroupMember,
		}

	} else {
		isSendByMySelf = fromUserName == toUserName
		if isSendByMySelf {
			senderUserInfo = SenderUserInfo{
				UserName: senderUserName,
				NickName: this.me.NickName,
				ContactType: 0,
			}
		} else {
			senderUserInfo = SenderUserInfo{
				UserName: senderUserName,
				NickName: this.contacts[senderUserName].NickName,
				ContactType: this.contacts[senderUserName].Type,
			}
		}
	}

	data := MsgEventData {
		IsGroupMsg:		isGroupMsg,
		IsMediaMsg:		isMediaMsg,
		IsSendByMySelf:	isSendByMySelf,
		MsgType:		int64(msgType),
		AtMe:			isAtMe,
		MediaUrl:		mediaUrl,
		Content:		content,
		FromUserName:	fromUserName,
		FromUserInfo:	*this.contacts[fromUserName],
		SenderUserName:	senderUserName,
		SenderUserInfo:	senderUserInfo,
		ToUserName:		toUserName,
		ToUserInfo:		*this.contacts[toUserName],
		OriginalMsg:	msg,
	}

	event := Event {
		Type:	MsgEvent,
		Time:	time.Now().Unix(),
		Data:	data,
	}

	this.listener.handle(event)
}

func (this *Wechat) emitContactChangeEvent(contactChangeType ContactChangeType, userName string) {
	data := ContactEventData {
		ChangeType: contactChangeType,
		UserName:	userName,
	}

	this.listener.handle(Event{
		Type: ContactChangeEvent,
		Time: time.Now().Unix(),
		Data: data,
	})
}