package gochat

import (
	"time"
	"strings"
	"fmt"
)

type Contacts struct {

}

func (this *Wechat) Getcontacts() error {
	cookies := this.HttpClient.HttpHeader.Cookies
	this.HttpClient.HttpHeader = HttpHeader{
		"",
		"",
		"",
		"",
		"",
		"",
		nil,
		"",
		"",
		"",
	}
	this.HttpClient.HttpHeader.Cookies = cookies
	getContactApi := strings.Replace(Config["getcontact_api"], "{pass_ticket}", this.PassTicket, 1)
	getContactApi = strings.Replace(getContactApi, "{skey}", this.Skey, 1)
	content, _, err := this.HttpClient.Get(getContactApi, time.Second * 5)
	if err != nil {
		return err
	}
	fmt.Println(content)

	return err
}