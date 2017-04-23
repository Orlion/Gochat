package gochat

import (
	"strings"
	"fmt"
	"time"
	"regexp"
	"encoding/json"
	"gochat/utils"
	"errors"
)

type syncMessageRequest struct {
	SyncKey     syncKey
	RR          int64 `json:"rr"`
	BaseRequest baseRequest
}

type syncMessageResponse struct {
	Response
	SyncKey      			syncKey
	SyncCheckKey 			syncKey
	SKey         			string
	ContinueFlag 			int
	AddMsgCount            	int
	AddMsgList             	[]map[string]interface{}
	ModContactCount        	int
	ModContactList         	[]map[string]interface{}
	DelContactCount        	int
	DelContactList         	[]map[string]interface{}
	ModChatRoomMemberCount 	int
	ModChatRoomMemberList  	[]map[string]interface{}
}

type syncKey struct {
	Count 	int
	List 	[]map[string]int64
}

var hosts = []string{
	"webpush.wx2.qq.com",
	"wx2.qq.com",
	"wx.qq.com",
	"webpush.wx.qq.com",
	"wx8.qq.com",
	"webpush.wx8.qq.com",
	"qq.com",
	"web2.wechat.com",
	"webpush.web2.wechat.com",
	"wechat.com",
	"webpush.web.wechat.com",
	"webpush.weixin.qq.com",
	"webpush.wechat.com",
	"webpush1.wechat.com",
	"webpush2.wechat.com",
	"webpush2.wx.qq.com",
}

// 开始长轮询
func (weChat *WeChat) beginListen() error {

	// 待优化
	if "wx.qq.com" == weChat.host {
		hosts[0] = "wx.qq.com"
		hosts[1] = "webpush.wx.qq.com"
		hosts[2] = "webpush.wx2.qq.com"
		hosts[3] = "wx2.qq.com"
	}

	weChat.logger.Println("[Info] Being Listen ... ")

	listenFailedCount := 0
	for {
		_, selector, err := weChat.listen()
		if err != nil {
			listenFailedCount++
			weChat.logger.Println("[Error] Listen Failed. Msg:" + err.Error() + fmt.Sprintf(", ListenFailedCount=%d.", listenFailedCount))
			weChat.triggerListenFailedEvent(listenFailedCount, weChat.host)
		} else {
			listenFailedCount = 0
		}

		// 接收到了消息
		if selector != "0" {
			continueFlag := -1
			// 持续接收消息直到continueFlag为0
			for continueFlag != 0 {
				resp, err := weChat.sync()
				if err != nil {
					return err
				}
				continueFlag = resp.ContinueFlag

				// 联系人有修改
				if resp.ModContactCount > 0 {
					weChat.contactsModify(resp.ModContactList)
				}

				// 联系人删除
				if resp.DelContactCount > 0 {
					weChat.contactsDelete(resp.DelContactList)
				}

				go weChat.handleSyncResponse(resp)
			}
		}
	}
}

// 监听服务器
func (weChat *WeChat) listen() (string, string, error) {

	syncCheckApi := strings.Replace(weChatApi["syncCheckApi"], "{r}", utils.GetUnixMsTime(), 1)
	syncCheckApi = strings.Replace(syncCheckApi, "{skey}", weChat.baseRequest.Skey, 1)
	syncCheckApi = strings.Replace(syncCheckApi, "{sid}", weChat.baseRequest.Sid, 1)
	syncCheckApi = strings.Replace(syncCheckApi, "{uin}", weChat.baseRequest.Uin, 1)
	syncCheckApi = strings.Replace(syncCheckApi, "{deviceid}", weChat.baseRequest.DeviceID, 1)
	syncCheckApi = strings.Replace(syncCheckApi, "{synckey}", weChat.formattedSyncCheckKey(), 1)
	syncCheckApi = strings.Replace(syncCheckApi, "{_}", utils.GetUnixTime(), 1)

	for i, host := range hosts {
		syncCheckApi = strings.Replace(syncCheckApi, "{host}", host, 1)
		syncCheckResContent, err := weChat.httpClient.get(syncCheckApi, time.Second * 26, &httpHeader{
			Accept:				"*/*",
			AcceptEncoding:		"gzip, deflate, sdch, br",
			AcceptLanguage:		"zh-CN,zh;q=0.8",
			Connection: 		"keep-alive",
			Host: 				host,
			Referer: 			"https://"+ weChat.host +"/?&lang=zh_CN",
		})
		if err != err {
			return "", "0", err
		}

		fmt.Println(host + " | " + syncCheckResContent)
		code, selector, err := weChat.analysisSelector(syncCheckResContent)
		if err != nil {
			return "", "0", err
		}

		if code == "0" {
			hosts[i] = hosts[0]
			hosts[0] = host
			return code, selector, nil
		}
	}

	return "", "0", errors.New("Code != 0")
}

// 监听到服务器通知后拉取数据
func (weChat *WeChat) sync() (*syncMessageResponse, error) {
	syncApi := strings.Replace(weChatApi["syncApi"], "{sid}", weChat.baseRequest.Sid, 1)
	syncApi = strings.Replace(syncApi, "{skey}", weChat.baseRequest.Skey, 1)
	syncApi = strings.Replace(syncApi, "{host}", weChat.host, 1)

	data, err := json.Marshal(syncMessageRequest {
		SyncKey:     weChat.syncKey,
		RR:          ^time.Now().Unix(),
		BaseRequest: weChat.baseRequest,
	})

	if err != nil {
		return nil, err
	}

	content, err := weChat.httpClient.post(syncApi, data, time.Second * 5, &httpHeader{
		ContentType:		"application/json;charset=utf-8",
		Host: 				weChat.host,
		Referer: 			"https://"+ weChat.host +"/?&lang=zh_CN",
	})
	if err != nil {
		return nil, err
	}

	var smr syncMessageResponse
	err = json.Unmarshal([]byte(content), &smr)
	if err != nil {
		return nil, err
	}

	if smr.SyncCheckKey.Count > 0 {
		weChat.syncKey = smr.SyncCheckKey
	} else {
		weChat.syncKey = smr.SyncKey
	}

	return &smr, err
}

// 解析从微信服务器返回的信息
func (weChat *WeChat) analysisSelector(syncCheckRes string) (string, string, error) {

	reg, err := regexp.Compile(`window.synccheck=\{retcode:"(.+)",selector:"(.+)"\}`)
	if err != nil {
		return "", "", err
	}
	selectorArr := reg.FindSubmatch([]byte(syncCheckRes))
	if len(selectorArr) != 3 {
		return "", "", nil
	}

	return string(selectorArr[1]), string(selectorArr[2]), nil
}

// 格式化syncKey为请求参数
func (weChat *WeChat) formattedSyncCheckKey() string {

	syncKeys := []string{}

	for _, k2v := range weChat.syncKey.List {
		syncKeys = append(syncKeys, fmt.Sprintf("%v_%v", k2v["Key"], k2v["Val"]))
	}

	return strings.Join(syncKeys, "|")
}