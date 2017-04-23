package gochat

import (
	"time"
	"fmt"
	"strings"
	"encoding/json"
	"io"
	"log"
	"os"
	"gochat/utils"
	"errors"
)

type  WeChat struct {
	Uuid 				string
	baseRequest 		baseRequest
	passTicket 			string
	syncKey   			syncKey
	syncHost			string
	host     			string
	me 					Contact
	contacts			map[string]*Contact
	httpClient			*httpClient
	storage				*storage
	logger 				*log.Logger
	listeners			map[EventType]func(Event)
}

type initRequest struct {
	BaseRequest baseRequest
}

type initResp struct {
	Response
	User    	Contact
	Skey    	string
	SyncKey 	syncKey
}

// New A WeChat
func NewWeChat(storageFilePath string, logFile io.Writer) *WeChat {
	storage := storage {
		filePath:	storageFilePath,
	}

	if logFile == nil {
		logFile = os.Stdout
	}
	logger := log.New(logFile, "", log.Ldate | log.Ltime)

	return & WeChat {
		httpClient: &httpClient{},
		storage: 	&storage,
		listeners: 	map[EventType]func(Event){},
		logger:		logger,
	}
}

// Run
func (weChat *WeChat) Run() error {

	err := weChat.login()
	if err != nil {
		return err
	}

	err = weChat.init()
	if err != nil {
		weChat.storage.delData()
		err = weChat.login()
		if err != nil {
			return err
		}
		err = weChat.init()
		if err != nil {
			return err
		}
	}

	go weChat.triggerInitEvent(weChat.me)
	weChat.logger.Println("[Info] WeChat Init.")

	err = weChat.initContact()
	if err != nil {
		return err
	}

	go weChat.triggerContactsInitEvent(len(weChat.contacts))
	weChat.logger.Println("[Info] Contacts Init.")

	weChat.beginListen()

	return nil
}

// init
func (weChat *WeChat) init() error {
	wxInitApi := strings.Replace(weChatApi["initApi"], "{r}", utils.GetUnixTime(), 1)
	wxInitApi = strings.Replace(wxInitApi, "{host}", weChat.host, 1)
	wxInitApi = strings.Replace(wxInitApi, "{pass_ticket}", weChat.passTicket, 1)

	postData, err := json.Marshal(initRequest{
		BaseRequest: weChat.baseRequest,
	})
	if err != nil {
		return err
	}

	content, err := weChat.httpClient.post(wxInitApi, postData, time.Second * 5, &httpHeader{
		Accept:				"application/json, text/plain, */*",
		ContentType:		"application/json;charset=UTF-8",
		Origin:				"https://" + weChat.host,
		Host: 				weChat.host,
		Referer: 			"https://"+ weChat.host +"/?&lang=zh_CN",
	})
	if err != nil {
		return err
	}

	var initRes initResp
	err = json.Unmarshal([]byte(content), &initRes)
	if err != nil {
		return err
	}

	if initRes.Response.BaseResponse.Ret != 0 {
		weChat.logger.Println("[Error] Init Failed. Res.Ret=" + string(initRes.Response.BaseResponse.Ret))
		return errors.New("Init Failed")
	}

	weChat.me = initRes.User
	weChat.baseRequest.Skey = initRes.Skey
	weChat.syncKey = initRes.SyncKey

	return nil
}

func (weChat *WeChat) skeyKV() string {
	return fmt.Sprintf(`skey=%s`, weChat.baseRequest.Skey)
}