package gochat

import (
	"strings"
	"time"
	"encoding/json"
	"errors"
	"strconv"
)

type ContactType int	// 联系人类型

const (
	_ ContactType = iota
	Offical				// 公众号
	Friend				// 好友
	Group				// 群组
	GroupMember			// 群成员
)

type Contact struct {
	Uin					float64
	UserName        	string
	NickName        	string
	HeadImgUrl      	string
	ContactFlag     	float64
	MemberCount			float64
	MemberList      	[]*Memeber
	MemberMap			map[string]*Memeber
	RemarkName			string
	HideInputBarFlag	float64
	Sex					float64
	Signature       	string
	VerifyFlag      	float64
	PYInitial			string
	PYQuanPin			string
	RemarkPYInitial 	string
	RemarkPYQuanPin 	string
	StarFriend      	float64
	AppAccountFlag		float64
	Statues				float64
	AttrStatus			float64
	Province        	string
	City            	string
	Alias           	string
	SnsFlag				float64
	UniFriend			float64
	DisplayName     	string
	ChatRoomId			float64
	KeyWord				string
	EncryChatRoomId 	string
	IsOwner				float64
	Type            	ContactType
}

type Memeber struct {
	Uin 			float64
	UserName		string
	NickName		string
	AttrStatus		float64
	PYInitial		string
	PYQuanPin		string
	RemarkPYInitial	string
	RemarkPYQuanPin	string
	MemberStatus	float64
	DisplayName		string
	KeyWord			string
}

type getContactResponse struct {
	Response
	MemberCount int
	MemberList  []Contact
	Seq         float64
}

type batchGetContactResponse struct {
	Response
	Count       int
	ContactList []Contact
}

func (this *Wechat) initContact() error {
	seq := float64(-1)

	// 初始化为空
	var cts = []Contact{}
	this.contacts = map[string]*Contact{}

	for seq != 0 {
		if -1 == seq {
			seq = 0
		}
		contactList, s, err := this.getContacts(seq)
		if err != nil {
			return err
		}
		seq = s
		cts = append(cts, contactList...)
	}

	// 初始化群的成员列表
	var groupUserNames []string

	for _, v := range cts {
		verifyFlag := v.VerifyFlag
		userName := v.UserName

		if verifyFlag / 8 != 0 {
			v.Type = Offical
		} else if strings.HasPrefix(userName, "@@") {
			v.Type = Group
			groupUserNames = append(groupUserNames, userName)
		} else {
			v.Type = Friend
		}
		this.contacts[userName] = &v
	}

	groups, _ := this.fetchContacts(groupUserNames)
	for _, group := range groups {
		group.MemberMap = map[string]*Memeber{}
		for _, contact := range group.MemberList {
			group.MemberMap[contact.UserName] = contact
		}
		this.contacts[group.UserName] = &group
	}

	return nil
}

func (this *Wechat) getContacts(seq float64) ([]Contact, float64, error) {

	getContactsApiUrl := strings.Replace(Config["getcontact_api"], "{pass_ticket}", this.passTicket, 1)
	getContactsApiUrl = strings.Replace(getContactsApiUrl, "{seq}", strconv.FormatInt(int64(seq), 10), 1)
	getContactsApiUrl = strings.Replace(getContactsApiUrl, "{skey}", this.baseRequest.Skey, 1)
	getContactsApiUrl = strings.Replace(getContactsApiUrl, "{r}", this.utils.getUnixTime(), 1)
	getContactsApiUrl = strings.Replace(getContactsApiUrl, "{host}", this.host, 1)

	content, err := this.httpClient.get(getContactsApiUrl, time.Second * 5, &HttpHeader{
		Accept: 			"application/json, text/plain, */*",
		AcceptEncoding: 	"gzip, deflate, br",
		AcceptLanguage: 	"zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3",
		Connection: 		"keep-alive",
		Host: 				"login.wx2.qq.com",
		Referer: 			"https://wx2.qq.com/?&lang=zh_CN",
	})

	var resp getContactResponse
	err = json.Unmarshal([]byte(content), &resp)
	if err != nil {
		return nil, float64(0), err
	}

	return resp.MemberList, resp.Seq, nil
}

/**
 * 根据qunUserName获取MemberList
 */
func (this *Wechat) fetchContacts(userNames []string) ([]Contact, error) {

	var list []map[string]string

	for _, u := range userNames {
		list = append(list, map[string]string{
			"UserName": u,
			"ChatRoomId": "",
		})
	}

	data, err := json.Marshal(map[string]interface{}{
		"BaseRequest": 	this.baseRequest,
		"Count":		len(list),
		"List":			list,
	})

	if err != nil {
		return nil, err
	}

	batchGetcontactApi := strings.Replace(Config["batchgetcontact_api"], "{r}", this.utils.getUnixMsTime(), 1)
	batchGetcontactApi = strings.Replace(batchGetcontactApi, "{host}", this.host, 1)
	content, err := this.httpClient.post(batchGetcontactApi, data, time.Second * 5, &HttpHeader{
		Host: 				this.host,
		Referer: 			"https://wx2.qq.com/?&lang=zh_CN",
	})

	var resp batchGetContactResponse
	err = json.Unmarshal([]byte(content), &resp)
	if err != nil {
		return nil, err
	}

	return resp.ContactList, nil
}

/**
 * 根据UserName添加新成员
 */
func (this *Wechat) updateOrAddContact(userNames []string) (int, error) {

	addNum := 0

	contacts, err := this.fetchContacts(userNames)

	if err != nil || len(contacts) != 1 {
		return 0, errors.New("Fetch contacts failed.")
	}

	for _, contact := range contacts {
		contact.MemberMap = map[string]*Memeber{}
		for _, member := range contact.MemberList {
			contact.MemberMap[member.UserName] = member
		}

		if contact.VerifyFlag / 8 != 0 {
			contact.Type = Offical
		} else if strings.HasPrefix(contact.UserName, "@@") {
			contact.Type = Group
		} else {
			contact.Type = Friend
		}
		// 添加到通讯录
		this.contacts[contact.UserName] = &contact
		addNum++
	}

	return addNum, nil
}

/**
 *
 */
func (this *Wechat) contactsModify(cts []map[string]interface{}) (int, error) {
	userNames := []string{}
	for _, newContact := range cts {
		userNames = append(userNames, newContact["UserName"].(string))
	}

	return this.updateOrAddContact(userNames)
}

func (this *Wechat) contactsDelete(cts []map[string]interface{}) {
	for _, contact := range cts {
		delete(this.contacts, contact["UserName"].(string))
	}
}