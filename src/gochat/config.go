package gochat

var Config = map[string]string {
	"getUuidApi": 			"https://login.wx2.qq.com/jslogin?appid=wx782c26e4c19acffb&redirect_uri=https%3A%2F%2Fwx2.qq.com%2Fcgi-bin%2Fmmwebwx-bin%2Fwebwxnewloginpage&fun=new&lang=zh_CN&_=",
	"login_qrcode_api": 	"https://login.weixin.qq.com/qrcode/",
	"login_poll_api": 		"https://login.wx2.qq.com/cgi-bin/mmwebwx-bin/login?uuid={uuid}&tip={tip}&_={time}",
	"get_contact_api": 		"https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxgetcontact?lang=zh_CN&pass_ticket={pass_ticket}&seq=0&skey={skey}",
	"wx_init_api": 			"https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxinit?r=-{r}",
	"wx_statusnotify_api": 	"https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxstatusnotify?lang=zh_CN&pass_ticket={pass_ticker}",
	"getcontact_api": 		"https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxgetcontact?lang=zh_CN&pass_ticket={pass_ticket}&r={r}&seq=0&skey={skey}",
	"batchgetcontact_api": 	"https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxbatchgetcontact?type=ex&r={r}&lang=zh_CN&pass_ticket={pass_ticket}",
	"sync_api": 			"https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxsync?sid={sid}&skey={skey}&lang=zh_CN&pass_ticket={pass_ticket}",
	// 暂时没用到
	"statreport_api": 		"https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxstatreport?fun=new&lang=zh_CN&pass_ticket={pass_ticket}",
	"synccheck_api": 		"https://webpush.wx2.qq.com/cgi-bin/mmwebwx-bin/synccheck?r={r}&skey={skey}&sid={sid}&uin={uid}&deviceid={deviceid}&synckey=1_660902725%7C2_660902839%7C3_660902841%7C1000_1491905341&_=1491916243535",
	"sendmsg_api": 			"https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxsendmsg?lang=zh_CN&pass_ticket={pass_ticket}",
}