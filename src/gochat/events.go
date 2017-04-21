package gochat

import (
	"strings"
	"fmt"
	"time"
)

type EventType int

const (
	_ EventType = iota
	GenUuidEvent		// 生成Uuid
	ScanCodeEvent		// 已扫码，未确认
	ConfirmAuthEvent	// 已确认授权登录
	InitEvent			// 初始化完成
	ContactsInitEvent	// 联系人初始化完
	ContactChangeEvent	// 联系人改变了
	MsgEvent			// 消息
	FriendReqEvent		// 好友申请
	LocationEvent		// 位置消息
)

type Event struct {
	Time int64
	Data interface{}
}

type MsgEventData struct{
	IsGroupMsg       bool
	IsMediaMsg       bool
	IsLocationMsg 	 bool
	IsSendByMySelf   bool
	MsgType          int64
	AtMe             bool
	MediaUrl         string
	Content          string
	FromUserName     string
	FromUserInfo     Contact
	SenderUserName   string
	SenderUserInfo   SenderUserInfo
	SenderUserId	 string					// 根据SendUserName生成ID
	ToUserName       string
	ToUserInfo       Contact
	OriginalMsg      map[string]interface{}
	LocationX		 string					// 位置，eg: 36.093239
	LocationY 		 string					// 位置，eg：123.376060
	LocationLabel  	 string					// 位置文本
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

// 有好友请求时会填充该结构体
type FriendReqEventData struct {
	Alias		string
	AttrStatus	string
	City		string
	Content		string
	NickName	string
	OpCode		string
	Province	string
	QQNum		string
	Scene		string
	Sex			string
	Signature	string
	Ticket		string
	UserName	string
	VerifyFlag	string
}


type SenderUserInfo struct {
	UserName	string
	NickName	string
	ContactType ContactType	// 发送人类型(好友，群成员，好友兼群成员)
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
	if "fmessage" == fromUserName { // 加好友消息
		friendReqEventData, err := msg["RecommendInfo"].(FriendReqEventData)
		fmt.Println(err)
		if err {
			event := Event {
				Time:	time.Now().Unix(),
				Data:	friendReqEventData,
			}
			handler, found := this.handlers[FriendReqEvent]
			if found {
				handler(event)
			}
		}

		return
	}
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
	case 37: {
		// 好友请求
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

		contact := &Memeber{}
		for {
			fromGroup, found := this.contacts[fromUserName]
			if found {
				contact, found  = fromGroup.MemberMap[infos[0]] // 根据content中UserName(消息发布人)找到详细数据
				if found {
					break
				}
			}

			_,err := this.updateOrAddContact([]string{fromUserName})
			if err != nil {
				return
			}

			contact, found = this.contacts[fromUserName].MemberMap[infos[0]]
			if !found {
				return
			}
		}

		if nil == contact {
			return
		}

		senderUserName = infos[0] // 实际发布人
		senderUserInfo = SenderUserInfo{
			UserName: infos[0],
			NickName: contact.NickName,
			ContactType: GroupMember,
		}

	} else {

		isSendByMySelf = fromUserName == this.me.UserName
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

	fromUserInfo := this.me
	if !isSendByMySelf {
		fromUserInfo = *this.contacts[fromUserName]
	}

	toUserInfo := this.me
	if toUserName != this.me.UserName {
		toUserInfo = *this.contacts[toUserName]
	}

	senderUserId := this.utils.userName2Id(senderUserName)
	data := MsgEventData {
		IsGroupMsg:		isGroupMsg,
		IsMediaMsg:		isMediaMsg,
		IsSendByMySelf:	isSendByMySelf,
		MsgType:		int64(msgType),
		AtMe:			isAtMe,
		MediaUrl:		mediaUrl,
		Content:		content,
		FromUserName:	fromUserName,
		FromUserInfo:	fromUserInfo,
		SenderUserName:	senderUserName,
		SenderUserInfo:	senderUserInfo,
		SenderUserId:   senderUserId,
		ToUserName:		toUserName,
		ToUserInfo:		toUserInfo,
		OriginalMsg:	msg,
		LocationX:		"",
		LocationY:		"",
		LocationLabel:  "",
	}

	event := Event {
		Time:	time.Now().Unix(),
		Data:	data,
	}
	handler, found := this.handlers[MsgEvent]
	if found {
		handler(event)
	}
}

func (this *Wechat) emitContactChangeEvent(contactChangeType ContactChangeType, userName string) {
	data := ContactEventData {
		ChangeType: contactChangeType,
		UserName:	userName,
	}

	handler, found := this.handlers[ContactChangeEvent]
	if found {
		handler(Event{
			Time: time.Now().Unix(),
			Data: data,
		})
	}
}

func (this *Wechat) LocationMsgEvent() {
	/*
	OriContent
	<?xml version="1.0"?>
	<msg>
		<location x="36.095364" y="120.373940" scale="16" label="市北区鞍山路(青岛鞍山路小学南)" maptype="0" poiname="[位置]" />
	</msg>
	*/
}