package wechatBot

import "fmt"

type WechatBot struct {
	Wechat Wechat
}

func (wechatBot *WechatBot) Init() {
	var wechat = Wechat{}
	wechat.Init()
	wechatBot.Wechat = wechat
}

func (wechatBot *WechatBot) Login() {
	uuid, _ := wechatBot.Wechat.GetUuid()
	fmt.Println(uuid)
	wechatBot.Wechat.GetRedirectUrl()
	xmlStr, _ := wechatBot.Wechat.WaitForLogin()
	wechatBot.Wechat.AnalysisLoginXml(xmlStr)
	wechatBot.Wechat.WxInit()
}
