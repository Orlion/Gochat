package gochat

import (
	"fmt"
	"io"
	"log"
	"os"
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

// Login And Init
func (weChat *WeChat) Login() error {

	err := weChat.beginLogin()
	if err != nil {
		return err
	}

	err = weChat.init()
	if err != nil {
		weChat.storage.delData()
		err = weChat.beginLogin()
		if err != nil {
			return err
		}
		err = weChat.init()
		if err != nil {
			return err
		}
	}

	weChat.triggerInitEvent(weChat.me)
	weChat.logger.Println("[Info] WeChat Init.")

	err = weChat.initContact()
	if err != nil {
		return err
	}

	weChat.triggerContactsInitEvent(len(weChat.contacts))
	weChat.logger.Println("[Info] Contacts Init.")

	return nil
}

func (weChat *WeChat) Run() error {
	err := weChat.beginListen();
	return err
}

func (weChat *WeChat) skeyKV() string {
	return fmt.Sprintf(`skey=%s`, weChat.baseRequest.Skey)
}