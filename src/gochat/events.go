package gochat

import (
	"strings"
	"fmt"
	"time"
	"gochat/utils"
)

// 事件类型
type EventType int

const (
	_ EventType = iota
	GenUuidEvent		// 生成Uuid
	ScanCodeEvent		// 已扫码，未确认
	ConfirmAuthEvent	// 已确认授权登录
	LoginEvent 			// 已登录
	InitEvent			// 初始化完成
	ContactsInitEvent	// 联系人初始化完
	ListenFailedEvent 	// 同步微信失败,可能为客户端已退出 | 被微信反爬虫
	ContactChangeEvent	// 联系人改变了
	ContactDeleteEvent  // 联系人删除事件
	MessageEvent		// 消息
)

// 事件体
type Event struct {
	Time 		int64
	EventType 	EventType
	Data 		interface{}
}

// 生成Uuid事件的数据
type GenUuidEventData struct {
	Uuid 	string
}

// 扫码事件数据
type ScanCodeEventData struct {
	UserAvatar 	string
}

// 授权事件数据
type ConfirmAuthEventData struct {
	RedirectUrl 	string
}

// 登录事件数据
type LoginEventData struct {
	DeviceID 	string
}

// 初始化事件数据
type InitEventData struct {
	Me 	Contact
}

// 通讯录初始化事件数据
type ContactsInitEventData struct {
	ContactsCount 	int
}

// 同步微信失败事件数据
type ListenFailedEventData struct {
	ListenFailedCount 	int
	Host 				string
}

// 联系人修改事件数据
type ContactChangeEventData struct {
	UserNames 	[]string
}

// 联系人删除事件数据
type ContactDeleteEventData struct {
	UserNames 	[]string
}

// 消息事件数据
type MessageEventData struct{
	MessageType       	MessageType
	IsGroupMessage      bool
	IsSendByMySelf   	bool
	IsAtMe             	bool
	MediaUrl         	string
	Content          	string
	FromUserName     	string
	FromUserInfo     	Contact
	SenderUserInfo   	SenderUserInfo
	SenderUserId	 	string					// 根据SendUserName生成ID
	ToUserName       	string
	ToUserInfo       	Contact
	RecommendInfo		map[string]interface{}
	LocationInfo		LocationInfo
	OriginalMsg      	map[string]interface{}
}

// 消息类型
type MessageType int
const (
	_ MessageType = iota
	TextMessage
	ImgMessage
	VoiceMessage
	VideoMessage
	CardMessage
	LocationMessage
	FriendReqMessage
)

// 发送人信息
type SenderUserInfo struct {
	UserName	string
	NickName	string
}

// 位置信息数据
type LocationInfo struct {
	X 		string
	Y 		string
	Label 	string
	Img 	string
}

// 设置事件监听器
func (weChat *WeChat) SetListener(eventType EventType, listener func(Event)) {
	weChat.listeners[eventType] = listener
}

// 处理从微信服务器拉过来的响应数据
func (weChat *WeChat) handleSyncResponse(resp *syncMessageResponse) {

	if resp.ModContactCount > 0 {
		userNames := []string{}
		for _, v := range resp.ModContactList {
			userNames = append(userNames, v["UserName"].(string))
		}
		go weChat.triggerContactChangeEvent(userNames)
	}

	if resp.DelContactCount > 0 {
		userNames := []string{}
		for _, v := range resp.ModContactList {
			userNames = append(userNames, v["UserName"].(string))
		}
		go weChat.triggerContactDeleteEvent(userNames)
	}

	if resp.AddMsgCount > 0 {
		for _, v := range resp.AddMsgList {
			go weChat.triggerMessageEvent(v)
		}
	}
}

// 触发生成uuid的事件
func (weChat *WeChat) triggerGenUuidEvent(uuid string) {
	listener, isReg := weChat.listeners[GenUuidEvent]
	if isReg {
		listener(Event{
			Time: 		time.Now().Unix(),
			EventType:	GenUuidEvent,
			Data: 		GenUuidEventData {
				Uuid:	uuid,
			},
		})
	}
}

// 触发扫码事件(未确认)
func (weChat *WeChat) triggerScanCodeEvent(userAvatar string) {
	listener, isReg := weChat.listeners[ScanCodeEvent]
	if isReg {
		listener(Event{
			Time: 		time.Now().Unix(),
			EventType: 	ScanCodeEvent,
			Data: 		ScanCodeEventData {
				UserAvatar:		userAvatar,
			},
		})
	}
}

// 触发授权登录事件
func (weChat *WeChat) triggerConfirmAuthEvent(redirectUrl string) {
	listener, isReg := weChat.listeners[ConfirmAuthEvent]
	if isReg {
		listener(Event{
			Time: 		time.Now().Unix(),
			EventType: 	ConfirmAuthEvent,
			Data: 		ConfirmAuthEventData {
				RedirectUrl:		redirectUrl,
			},
		})
	}
}

// 触发登录事件
func (weChat *WeChat) triggerLoginEvent(deviceID string) {
	listener, isReg := weChat.listeners[LoginEvent]
	if isReg {
		listener(Event{
			Time: 		time.Now().Unix(),
			EventType: 	LoginEvent,
			Data: 		LoginEventData {
				DeviceID:		deviceID,
			},
		})
	}
}

// 触发初始化事件
func (weChat *WeChat) triggerInitEvent(me Contact) {
	listener, isReg := weChat.listeners[InitEvent]
	if isReg {
		listener(Event{
			Time: 		time.Now().Unix(),
			EventType: 	InitEvent,
			Data: 		InitEventData {
				Me:		me,
			},
		})
	}
}

// 触发通讯录初始化事件
func (weChat *WeChat) triggerContactsInitEvent(contactsCount int) {
	listener, isReg := weChat.listeners[ContactsInitEvent]
	if isReg {
		listener(Event{
			Time: 		time.Now().Unix(),
			EventType: 	ContactsInitEvent,
			Data: 		ContactsInitEventData {
				ContactsCount:		contactsCount,
			},
		})
	}
}

func (weChat *WeChat) triggerListenFailedEvent(listenFailedCount int, host string) {
	listener, isReg := weChat.listeners[ListenFailedEvent]
	if isReg {
		listener(Event{
			Time: 		time.Now().Unix(),
			EventType: 	ListenFailedEvent,
			Data: 		ListenFailedEventData {
				ListenFailedCount:		listenFailedCount,
				Host:					host,
			},
		})
	}
}

// 触发通讯录修改事件
func (weChat *WeChat) triggerContactChangeEvent(userNames []string) {
	listener, isReg := weChat.listeners[ContactChangeEvent]
	if isReg {
		listener(Event{
			Time: 		time.Now().Unix(),
			EventType: 	ContactChangeEvent,
			Data: 		ContactChangeEventData {
				UserNames:	userNames,
			},
		})
	}
}

// 触发通讯录删除事件
func (weChat *WeChat) triggerContactDeleteEvent(userNames []string) {
	listener, isReg := weChat.listeners[ContactDeleteEvent]
	if isReg {
		listener(Event{
			Time: 		time.Now().Unix(),
			EventType: 	ContactDeleteEvent,
			Data: 		ContactDeleteEventData {
				UserNames:	userNames,
			},
		})
	}
}

// 触发消息事件
func (weChat *WeChat) triggerMessageEvent(msg map[string]interface{}) {

	messageType := TextMessage
	isGroupMessage := false
	isSendByMySelf := false
	isAtMe := false
	mediaUrl := ""
	content := msg["Content"].(string)
	fromUserName := msg["FromUserName"].(string)
	senderUserInfo := SenderUserInfo{}
	senderUserId := ""
	toUserName := msg["ToUserName"].(string)
	recommendInfo := map[string]interface{}{}
	locationInfo := LocationInfo{}
	senderUserName := fromUserName

	var groupUserName string
	if strings.HasPrefix(fromUserName, "@@") {
		groupUserName = fromUserName
	} else if strings.HasPrefix(toUserName, "@@") {
		groupUserName = toUserName
	}

	if len(groupUserName) > 0 {
		isGroupMessage = true
	}

	msgType := msg["MsgType"].(float64)
	mid := msg["MsgId"].(string)

	path := ""
	switch msgType {
	case 3: {
		messageType = ImgMessage
		path = "webwxgetmsgimg"
	}
	case 47: {
		pid, _ := msg["HasProductId"].(float64)
		if pid == 0 {
			messageType = ImgMessage
			path = "webwxgetmsgimg"
		}
	}
	case 34: {
		messageType = VoiceMessage
		path = "webwxgetvoice"
	}
	case 43: {
		messageType = VideoMessage
		path = "webwxgetvideo"
	}
	case 37: {
		messageType = FriendReqMessage
		recommendInfo, _ = msg["RecommendInfo"].(map[string]interface{})
	}
	case 42: {
		messageType = CardMessage
	}
	}
	if len(path) > 0 {
		mediaUrl = fmt.Sprintf(`https://wx2.qq.com/%s?msgid=%v&%v`, path, mid, weChat.skeyKV())
	}

	subMsgType, found := msg["SubMsgType"]
	if found && 48 == subMsgType.(float64){
		messageType = LocationMessage
		locationX, locationY, locationLabel, err := utils.GetLocationInfoFromOriContent(msg["OriContent"].(string))
		if err == nil {
			locationInfo.X = locationX
			locationInfo.Y = locationY
			locationInfo.Label = locationLabel
		}

		locationImg,err := utils.GetLocationImgFromContent(content)
		if err == nil {
			locationInfo.Img = "https://" + weChat.host + locationImg
		}
	}

	if isGroupMessage {
		atme := "@"
		if len(weChat.me.DisplayName) > 0 {
			atme += weChat.me.DisplayName
		} else {
			atme += weChat.me.NickName
		}
		isAtMe = strings.Contains(content, atme) // 标识是否是@我

		infos := strings.Split(content, ":<br/>")
		if len(infos) != 2{
			return
		}

		content = infos[1]
		if isAtMe && infos[0] == weChat.me.UserName {
			isSendByMySelf = true
		}

		contact := &Member{}
		for {
			fromGroup, found := weChat.contacts[fromUserName]
			if found {
				contact, found  = fromGroup.MemberMap[infos[0]] // 根据content中UserName(消息发布人)找到详细数据
				if found {
					break
				}
			}

			err := weChat.updateContact([]string{fromUserName})
			if err != nil {
				return
			}

			contact, found = weChat.contacts[fromUserName].MemberMap[infos[0]]
			if !found {
				return
			}
		}

		if nil == contact {
			return
		}

		senderUserName = infos[0]
		senderUserInfo = SenderUserInfo{
			UserName: infos[0],
			NickName: contact.NickName,
		}
	} else {

		isSendByMySelf = fromUserName == weChat.me.UserName
		if isSendByMySelf {
			senderUserInfo = SenderUserInfo{
				UserName: senderUserName,
				NickName: weChat.me.NickName,
			}
		} else {
			senderUserInfo = SenderUserInfo{
				UserName: senderUserName,
				NickName: "",
			}
			senderUser, found := weChat.contacts[senderUserName]
			if found {
				senderUserInfo.NickName = senderUser.NickName
			}
		}
	}


	fromUserInfo := weChat.me
	if !isSendByMySelf {
		fromUserInfoTemp, found := weChat.contacts[fromUserName]
		if found {
			fromUserInfo = *fromUserInfoTemp
		}
	}

	toUserInfo := weChat.me
	if toUserName != weChat.me.UserName {
		toUserInfoTemp, found := weChat.contacts[toUserName]
		if found {
			toUserInfo = *toUserInfoTemp
		}
	}

	senderUserId = utils.UserName2Id(senderUserName)

	event := Event {
		Time:		time.Now().Unix(),
		EventType: 	MessageEvent,
		Data:		MessageEventData {
			MessageType:		messageType,
			IsGroupMessage:		isGroupMessage,
			IsSendByMySelf:		isSendByMySelf,
			IsAtMe:				isAtMe,
			MediaUrl:			mediaUrl,
			Content:			content,
			FromUserName:		fromUserName,
			FromUserInfo:		fromUserInfo,
			SenderUserInfo:		senderUserInfo,
			SenderUserId:   	senderUserId,
			ToUserName:			toUserName,
			ToUserInfo:			toUserInfo,
			RecommendInfo:		recommendInfo,
			LocationInfo: 		locationInfo,
			OriginalMsg:		msg,
		},
	}

	weChat.logger.Println("[Info] Get Message. SenderNickName=[" + senderUserInfo.NickName + "], Content=[" + content + "]")
	listener, isReg := weChat.listeners[MessageEvent]
	if isReg {
		listener(event)
	}
}

