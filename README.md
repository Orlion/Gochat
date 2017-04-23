# Gochat
WeChat Bot

# Demo
```
package main

import (
	"gochat"
	"os"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"errors"
	"bytes"
	"net"
	"time"
)
func main() {
	weChat := gochat.NewWeChat("storage.json", os.Stdout)
	MessageListener(weChat)
	err := weChat.Login()
	if err != nil {
		fmt.Println(err.Error())
	}

	err = weChat.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func MessageListener(weChat *gochat.WeChat) {
	weChat.SetListener(gochat.MESSAGE_EVENT, func(event gochat.Event){
		eventData, ok := event.Data.(gochat.MessageEventData)
		if ok {
			if eventData.IsGroupMessage {
				if eventData.IsAtMe {
					res, err := tuling(eventData.Content, "青岛",eventData.SenderUserId)
					if err != nil {
						weChat.SendTextMsg("@"+ eventData.SenderUserInfo.NickName +" "+"短路了...快通知我主人修修我...", eventData.FromUserName)
					} else {
						weChat.SendTextMsg("@"+ eventData.SenderUserInfo.NickName +" "+res, eventData.FromUserName)
					}
				}
			} else {

				if gochat.FriendReqMessage == eventData.MessageType {
					reqUserName, okU := eventData.RecommendInfo["UserName"].(string)
					reqTicket, okT := eventData.RecommendInfo["Ticket"].(string)
					if okU && okT {
						weChat.VerifyUser(reqUserName, reqTicket, "你好, I am Oosten.")
					}
				} else {
					res, err := tuling(eventData.Content, eventData.FromUserInfo.City, eventData.SenderUserId)
					if err != nil || res == ""{
						weChat.SendTextMsg("短路了...快通知我主人修修我...", eventData.SenderUserInfo.UserName)
					} else {
						weChat.SendTextMsg(res, eventData.SenderUserInfo.UserName)
					}
				}
			}
		}
	})
}
```