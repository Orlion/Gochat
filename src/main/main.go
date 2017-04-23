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

func tuling(input string, city string, userId string) (string, error) {
	api := "http://openapi.tuling123.com/openapi/api/v2"

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(time.Second * 5)
				c, err := net.DialTimeout(netw, addr, time.Second * 5)
				if err != nil {
					return nil, err
				}

				c.SetDeadline(deadline)
				return c, nil
			},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("Cannot Redirect")
		},
	}

	data := `{"perception": {"inputText": {"text": "`+ input +`"},"selfInfo": {"location": {"city": "`+ city +`"},}},"userInfo": {"apiKey": "","userId": `+ userId +`}}`
	req, err := http.NewRequest("POST", api, bytes.NewReader([]byte(data)))
	if err != nil {
		return ``, err
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	type Resp struct {
		Intent interface{}
		Results []struct{
			ResultType string
			Values	map[string]string
		}
	}
	var tulingResp = new(Resp)
	err = json.Unmarshal(body, &tulingResp)
	if err != nil {
		return "", nil
	}

	resultText, resultUrl := "", ""
	for _,v := range tulingResp.Results {
		if (v.ResultType == "text") {
			resultText, _ = v.Values["text"]
		}

		if (v.ResultType == "url") {
			resultUrl, _ = v.Values["url"]
		}
	}

	return resultText + resultUrl, nil
}