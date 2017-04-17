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
)

func main() {
	tuling("你好", "1")
	//wechat := gochat.NewWechat(gochat.Option{
	//	StorageDirPath: "D:/develop/",
	//})
	//wechat.Handle(gochat.MsgEvent, func(event gochat.Event){
	//	eventData, ok := event.Data.(gochat.MsgEventData)
	//	if ok {
	//		if eventData.IsGroupMsg {
	//			if eventData.AtMe {
	//				wechat.SendTextMsg(eventData.FromUserName, "@" + eventData.SenderUserInfo.NickName + "你好")
	//			}
	//		} else {
	//			wechat.SendTextMsg(eventData.SenderUserName, "你好！")
	//		}
	//	}
	//})
	//err := wechat.Run()
	//if err != nil {
	//	fmt.Println(err)
	//}
}

func tuling(input string, userId string) (string, error) {
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

	data := `{"perception": {"inputText": {"text": "`+ input +`"},"selfInfo": {"location": {"city": "青岛"},}},"userInfo": {"apiKey": "aaaaaa","userId": `+ userId +`}}`
	req, err := http.NewRequest("POST", api, bytes.NewReader([]byte(data)))
	if err != nil {
		return ``, err
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	type Resp struct {
		Intent interface{}
		Results []interface{}
	}
	var resp Resp
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil
	}
	return resp, nil
}