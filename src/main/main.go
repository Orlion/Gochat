package main

import (
	"gochat"
	"fmt"
)

func main() {

	wechat := gochat.New()
	uuid, _ := wechat.Begin()
	fmt.Println("https://login.weixin.qq.com/qrcode/" + uuid)
	wechat.Login()
}
