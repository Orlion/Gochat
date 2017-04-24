![](http://i4.buimg.com/567571/4452ae08fdc6880b.jpg)

# Gochat
一个Go微信机器人包

# 特点
1. **灵活**。在微信从登录到开始同步服务器消息的过程中的各个节点触发事件，从而通过注册时间监听器就可以灵活的实现很多功能。  
2. **失败重新登录**。可以通过注册同步失败的事件重新调用Login()方法来重新登录，从而达到失败自动重新登录, 也可以调用pushlogin的接口免扫码来登录。

# Demo
> 有部分伪代码,不能直接运行
```
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
```

# Something
&nbsp;&nbsp;写完之后就没什么兴趣用这个库做东西了，主要是用的openshift太慢,而且微信很不稳定经常变更规则。

&nbsp;&nbsp;不过用微信机器人还是可以做很多有趣的事情的，有不少同学在玩微信机器人。可以用java写个Android的App。

&nbsp;&nbsp;类似项目有:
* https://github.com/littlecodersh/ItChat
* https://github.com/youfou/wxpy
* https://github.com/liuwons/wxBot  
...  
..  
.  