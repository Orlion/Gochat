package main

import (
	"wechatBot"
)

func main() {
	 // 获取uuid
	var wechatBot = wechatBot.WechatBot{}
	wechatBot.Init()
	wechatBot.Login()
}