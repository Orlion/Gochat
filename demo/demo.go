package demo

import (
	"gochat"
	"os"
	"fmt"
)

func main() {
	// new 一个微信实例
	weChat := gochat.NewWeChat("storage.json", os.Stdout)
	// 注册事件监听
	RegListener(weChat)
	// 登录
	err := weChat.Login()
	if err != nil {
		fmt.Println(err.Error())
	}
	// Run 开始监听微信服务器
	err = weChat.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
}

// 注册事件监听器
func RegListener(weChat *gochat.WeChat) {

	// 注册生成Uuid事件监听器
	weChat.SetListener(gochat.GEN_UUID_EVENT, func(event gochat.Event){
		eventData, ok := event.Data.(gochat.GenUuidEventData)
		if ok {
			sendEmail("Uuid=" + eventData.Uuid)
		}
	})

	// 注册已扫码事件监听器
	weChat.SetListener(gochat.SCAN_CODE_EVENT, func(event gochat.Event){
		eventData, ok := event.Data.(gochat.ScanCodeEventData)
		if ok {
			sendEmail("UserAvatar=" + eventData.UserAvatar)
		}
	})

	// 注册授权登录的事件监听器
	weChat.SetListener(gochat.CONFIRM_AUTH_EVENT, func(event gochat.Event){
		eventData, ok := event.Data.(gochat.ConfirmAuthEventData)
		if ok {
			sendEmail("RedirectUrl=" + eventData.RedirectUrl)
		}
	})

	// 注册已登录事件监听器
	weChat.SetListener(gochat.LOGIN_EVENT, func(event gochat.Event){
		eventData, ok := event.Data.(gochat.LoginEventData)
		if ok {
			sendEmail("DeviceID=" + eventData.DeviceID)
		}
	})

	// 注册初始化完成事件监听器
	weChat.SetListener(gochat.INIT_EVENT, func(event gochat.Event){
		eventData, ok := event.Data.(gochat.InitEventData)
		if ok {
			sendEmail("MemberCount=" + eventData.Me.MemberCount)
		}
	})

	// 注册联系人初始化完成事件监听器
	weChat.SetListener(gochat.CONTACTS_INIT_EVENT, func(event gochat.Event){
		eventData, ok := event.Data.(gochat.ContactsInitEventData)
		if ok {
			sendEmail("ContactsCount=" + eventData.ContactsCount)
		}
	})

	// 注册同步微信失败事件监听器
	weChat.SetListener(gochat.LISTEN_FAILED_EVENT, func(event gochat.Event){
		eventData, ok := event.Data.(gochat.ListenFailedEventData)
		if ok {
			sendEmail("ListenFailedCount=" + eventData.ListenFailedCount)
			if (eventData.ListenFailedCount > 10) {
				// 连续同步失败10次后重新登录
				weChat.Login()
			}
		}
	})

	// 注册联系人修改事件监听器
	weChat.SetListener(gochat.CONTACT_MODIFY_EVENT, func(event gochat.Event){
		eventData, ok := event.Data.(gochat.ContactModifyEventData)
		if ok {
			sendEmail("UserNames=" + eventData.UserNames)
		}
	})

	// 注册联系人删除事件监听器
	weChat.SetListener(gochat.CONTACT_DELETE_EVENT, func(event gochat.Event){
		eventData, ok := event.Data.(gochat.ContactDeleteEventData)
		if ok {
			sendEmail("UserNames=" + eventData.UserNames)
		}
	})

	// 注册消息事件监听器
	weChat.SetListener(gochat.MESSAGE_EVENT, func(event gochat.Event){
		eventData, ok := event.Data.(gochat.MessageEventData)
		if ok {
			if eventData.IsGroupMessage {
				if eventData.IsAtMe {
					weChat.SendTextMsg(tuling(eventData.Content, eventData.SenderUserId), eventData.SenderUserInfo.UserName)
				}
			} else {
				weChat.SendTextMsg(tuling(eventData.Content, eventData.SenderUserId), eventData.SenderUserInfo.UserName)
			}
		}
	})
}