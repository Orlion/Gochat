package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net"
	"time"
	"errors"
	"bytes"
	"encoding/json"
	"gochat"
)

func main() {
	wechat := gochat.NewWechat(gochat.Option{
		StorageDirPath: "",
	})
	wechat.Handle(gochat.MsgEvent, func(event gochat.Event){
		eventData, ok := event.Data.(gochat.MsgEventData)
		if ok {
			if eventData.IsGroupMsg {
				if eventData.AtMe {
					res, err := tuling(eventData.Content, "青岛",eventData.SenderUserId)
					if err != nil {
						wechat.SendTextMsg("@"+ eventData.SenderUserInfo.NickName +" "+"短路了...快通知我主人修修我...", eventData.FromUserName)
					} else {
						wechat.SendTextMsg("@"+ eventData.SenderUserInfo.NickName +" "+res, eventData.FromUserName)
					}
				}
			} else {
				res, err := tuling(eventData.Content, eventData.FromUserInfo.City, eventData.SenderUserId)
				fmt.Println(res)
				if err != nil || res == ""{
					wechat.SendTextMsg("短路了...快通知我主人修修我...", eventData.SenderUserName)
				} else {
					wechat.SendTextMsg(res, eventData.SenderUserName)
				}
			}
		}
	})

	wechat.Handle(gochat.FriendReqEvent, func(event gochat.Event){
		fmt.Println("frient event")
		eventData, ok := event.Data.(gochat.FriendReqEventData)
		if ok {
			wechat.VerifyUser(eventData.UserName, eventData.Ticket)
		}
	})
	err := wechat.Run()
	if err != nil {
		fmt.Println(err)
	}
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

	data := `{"perception": {"inputText": {"text": "`+ input +`"},"selfInfo": {"location": {"city": "`+ city +`"},}},"userInfo": {"apiKey": "x","userId": `+ userId +`}}`
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
	fmt.Println(string(body))
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
