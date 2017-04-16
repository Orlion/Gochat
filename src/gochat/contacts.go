package gochat

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