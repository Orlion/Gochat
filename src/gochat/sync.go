package gochat

import (
	"strings"
	"fmt"
	"time"
	"regexp"
)

func (this *Wechat) sync() {

}

func (this *Wechat) syncCheck() error {
	syncCheckApi := strings.Replace(Config["synccheck_api"], "{r}", this.Utils.GetUnixMsTime(), 1)
	syncCheckApi = strings.Replace(syncCheckApi, "{skey}", this.BaseRequest.Skey, 1)
	syncCheckApi = strings.Replace(syncCheckApi, "{sid}", this.BaseRequest.Sid, 1)
	syncCheckApi = strings.Replace(syncCheckApi, "{uid}", this.BaseRequest.Uin, 1)
	syncCheckApi = strings.Replace(syncCheckApi, "{deviceid}", this.BaseRequest.DeviceID, 1)
	syncCheckApi = strings.Replace(syncCheckApi, "{synckey}", this.formattedSyncCheckKey(), 1)

	syncCheckRes, _, err := this.HttpClient.Get(syncCheckApi, time.Second * 26)
	if err != err {
		return err
	}

	selector, err := this.analysisSelector(syncCheckRes)
	if err != nil {
		return err
	}

	if selector > 0 { // 有消息
		
	}

	return nil
}

func (this *Wechat) analysisSelector(syncCheckRes string) (int, error) {

	reg, err := regexp.Compile(`window.synccheck=\{retcode:"0",selector:"(\d+)"}`)
	if err != nil {
		return ``, err
	}
	selectorArr := reg.FindSubmatch([]byte(syncCheckRes))
	if len(selectorArr) != 2 {
		return ``, nil
	}
	selector := int(selectorArr[1])
	return selector, nil
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