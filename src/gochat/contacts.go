package gochat

import (
	"time"
	"strings"
	"fmt"
	"encoding/json"
)

type Contacts struct {

}

type Contact struct {
	GGID            string
	UserName        string
	NickName        string
	HeadImgUrl      string
	HeadHash        string
	RemarkName      string
	DisplayName     string
	StarFriend      float64
	Sex             float64
	Signature       string
	VerifyFlag      float64
	ContactFlag     float64
	HeadImgFlag     float64
	Province        string
	City            string
	Alias           string
	EncryChatRoomId string
	Type            int
	MemberList      []*Contact
}

type getContactResponse struct {
	Response
	MemberCount int
	MemberList 	[]map[string]interface{}
	Seq  		float64
}

func (this *Wechat) getContacts(seq float64) ([]map[string]interface{}, float64, error) {
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
	getContactApi = strings.Replace(getContactApi, "{seq}", string(seq), 1)
	getContactApi = strings.Replace(getContactApi, "{skey}", this.BaseRequest.Skey, 1)
	content, _, err := this.HttpClient.get(getContactApi, time.Second * 5)
	if err != nil {
		return nil, nil, err
	}

	var resp getContactResponse
	err = json.Unmarshal([]byte(content), &resp)
	if err != nil {
		return nil, nil, err
	}

	return resp.MemberCount, resp.Seq, nil
}

func (this *Wechat) SyncContact() error {
	seq := float64(-1)

	var cts []map[string]interface{}

	for seq != 0 {
		if -1 == seq {
			seq = 0
		}
		memberList, s, err := this.getContacts(seq)
		if err != nil {
			return err
		}
		seq = s
		cts = append(cts, memberList...)
	}

	var groupUserNames []string

	var tempIdxMap = make(map[string]int)

	for idx, v := range cts {
		vf, _ := v["VerifyFlag"].(float64)
		un, _ := v["UserName"].(string)

		if vf/8 != 0 {
			v["Type"] = "Offical"
		} else if strings.HasPrefix(un, "@@") {
			v["Type"] = "Group"
			groupUserNames = append(groupUserNames, un)
		} else {
			v["Type"] = "Friend"
		}
		tempIdxMap[un] = idx
	}

	groups, _ := this.fatchGroups(groupUserNames)
	for _, group := range groups {
		groupUserName := group["UserName"].(string)
		contacts := group["MemberList"].([]interface{})

		for _, c := range contacts {
			ct := c.(map[string]interface{})
			un := ct["UserName"].(string)
			if idx, found := tempIdxMap[un]; found {
				ctx[idx]["Type"] = "FriendAndMember"
			} else {
				ct["HeadImgUrl"] = fmt.Sprintf(`/cgi-bin/mmwebwx-bin/webwxgeticon?seq=0&username=%s&chatroomid=%s&skey=`, un, groupUserName)
				ct["Type"] = "Member"
				cts = append(cts, ct)
			}
		}

		group["Type"] = "Group"
		idx := tempIdxMap[groupUserName]
		cts[idx] = group
	}

	this.SyncContact()

	return nil
}

func (this *Wechat) contactDidChange(cts []map[string]interface{}, changeType string) {
	if "modify" == changeType { // 修改
		var mcts []map[string]interface{}
		for _, v := range cts {
			vf, _ := v["VerifyFlay"].(float64)
			un, _ := v["UserName"].(string)

			if vf/8 != 0 {
				v["Type"] = Offical
				mcts = append(mcts, v)
			} else if strings.HasPrefix(un, "@@") {
				this.ForceUpdateGroup(un)
			} else {
				v["Type"] = Friend
				mcts = append(mcts, v)
			}
		}
		this.appendContacts(mcts)
	} else {
		for _, v := range cts {
			this.removeContact(v["UserName"].(string))
		}
	}
}