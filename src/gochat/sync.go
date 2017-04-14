package gochat

import (
	"strings"
	"fmt"
	"time"
	"regexp"
	"encoding/json"
	"strconv"
)

type syncMessageRequest struct {
	SyncKey     map[string]interface{}
	RR          int64 `json:"rr"`
	BaseRequest BaseRequest
}

type syncMessageResponse struct {
	Response
	SyncKey      map[string]interface{}
	SyncCheckKey map[string]interface{}
	SKey         string
	ContinueFlag int

	// Content
	AddMsgCount            int
	AddMsgList             []map[string]interface{}
	ModContactCount        int
	ModContactList         []map[string]interface{}
	DelContactCount        int
	DelContactList         []map[string]interface{}
	ModChatRoomMemberCount int
	ModChatRoomMemberList  []map[string]interface{}
}

func (this *Wechat) beginSync() error {
	for {
		code, selector, err := this.syncCheck()
		if err != nil {
			return err
		}

		if code != "0" {
			return error("Syncing failed, please relogin. [code]:" + code)
		}

		// 接收到了消息
		if selector != "0" {
			continueFlag := -1
			// 持续接收消息直到continueFlag为0
			for continueFlag != "0" {
				resp, err := this.sync()
				if err != nil {
					return err
				}

				continueFlag = resp.ContinueFlag
				if resp.ModContactCount > 0 {

				}
				if resp.DelContactCount >0 {

				}
				if resp.ModChatRoomMemberCount >0 {

				}

				go this.handleServerEvent(resp)
			}
		}
	}
}

func (this *Wechat) sync() (*syncMessageResponse, error) {
	syncApi := strings.Replace(Config["sync_api"], "{sid}", this.BaseRequest.Sid, 1)
	syncApi = strings.Replace(syncApi, "{skey}", this.BaseRequest.Skey, 1)

	syncKeyf := make(map[string]interface{}, 0)
	keys := strings.Split(this.formattedSyncCheckKey(), "|")
	syncKeyf["Count"] = len(keys)
	list := make([]map[string]int64, 0)

	for _, key := range keys {
		kv := strings.Split(key, "_")
		k, _ := strconv.ParseInt(kv[0], 10, 64)
		v, _ := strconv.ParseInt(kv[1], 10, 64)
		kvmap := map[string]int64{"Key": k, "Val": v}
		list = append(list, kvmap)
	}
	syncKeyf["List"] = list
	data, err := json.Marshal(syncMessageRequest{
		BaseRequest: this.BaseRequest,
		SyncKey:     syncKeyf,
		RR:          ^time.Now().Unix(),
	})

	if err != nil {
		return nil, err
	}

	content, _, err := this.HttpClient.Post(syncApi, data, time.Second * 5)
	if err != nil {
		return "", err
	}

	var smr syncMessageResponse
	err = json.Unmarshal([]byte(content), &smr)
	if err != nil {
		return "", err
	}

	if smr.SyncCheckKey != nil {
		this.syncKey = smr.SyncCheckKey
	} else {
		this.syncKey = smr.SyncKey
	}

	return smr, err
}

func (this *Wechat) syncCheck() (string, string, error) {

	hosts := [...]string{
		`webpush.wx2.qq.com`,
		`webpush.wx.qq.com`,
		`wx2.qq.com`,
		`wx8.qq.com`,
		`webpush.wx8.qq.com`,
		`qq.com`,
		`webpush.wx.qq.com`,
		`web2.wechat.com`,
		`webpush.web2.wechat.com`,
		`wechat.com`,
		`webpush.web.wechat.com`,
		`webpush.weixin.qq.com`,
		`webpush.wechat.com`,
		`webpush1.wechat.com`,
		`webpush2.wechat.com`,
		`webpush2.wx.qq.com`,
	}

	for _, host := range hosts {
		syncCheckApi := strings.Replace(Config["synccheck_api"], "{host}", host, 1)
		syncCheckApi = strings.Replace(syncCheckApi, "{r}", this.Utils.GetUnixMsTime(), 1)
		syncCheckApi = strings.Replace(syncCheckApi, "{skey}", this.BaseRequest.Skey, 1)
		syncCheckApi = strings.Replace(syncCheckApi, "{sid}", this.BaseRequest.Sid, 1)
		syncCheckApi = strings.Replace(syncCheckApi, "{uid}", this.BaseRequest.Uin, 1)
		syncCheckApi = strings.Replace(syncCheckApi, "{deviceid}", this.BaseRequest.DeviceID, 1)
		syncCheckApi = strings.Replace(syncCheckApi, "{synckey}", this.formattedSyncCheckKey(), 1)

		syncCheckRes, _, err := this.HttpClient.Get(syncCheckApi, time.Second * 26)
		if err != err {
			return err
		}

		code, selector, err := this.analysisSelector(syncCheckRes)
		if err != nil {
			return err
		}

		if code == `0` {
			return code, selector, nil
		}
	}

	return ``, ``, nil
}

func (this *Wechat) analysisSelector(syncCheckRes string) (string, string, error) {

	reg, err := regexp.Compile(`window.synccheck=\{retcode:"(\.+)",selector:"(\.+)"}`)
	if err != nil {
		return ``, ``, err
	}
	selectorArr := reg.FindSubmatch([]byte(syncCheckRes))
	if len(selectorArr) != 3 {
		return ``, ``, nil
	}

	return selectorArr[1], selectorArr[2], nil
}

func (this *Wechat) formattedSyncCheckKey() string {

	keys, _ := this.syncKey["List"].([]interface{})

	synckeys := make([]string, 0)

	for _, key := range keys {
		kv, ok := key.(map[string]interface{})
		if ok {
			f, _ := kv["Val"].(float64)
			synckeys = append(synckeys, fmt.Sprintf("%v_%v", kv["Key"], int64(f)))
		}
	}

	return strings.Join(synckeys, "|")
}

func (this *Wechat) choseAvalibleSyncHost() bool {
	hosts := [...]string{
		`webpush.wx.qq.com`,
		`wx2.qq.com`,
		`webpush.wx2.qq.com`,
		`wx8.qq.com`,
		`webpush.wx8.qq.com`,
		`qq.com`,
		`webpush.wx.qq.com`,
		`web2.wechat.com`,
		`webpush.web2.wechat.com`,
		`wechat.com`,
		`webpush.web.wechat.com`,
		`webpush.weixin.qq.com`,
		`webpush.wechat.com`,
		`webpush1.wechat.com`,
		`webpush2.wechat.com`,
		`webpush2.wx.qq.com`}

	for _, host := range hosts {
		this.syncHost = host
		code, _, _ := this.syncCheck()
		if code == `0` {
			return true
		}
	}

	return false
}