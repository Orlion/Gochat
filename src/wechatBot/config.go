package wechatBot

var Config = map[string]string {
	"getUuidApi": "https://login.wx2.qq.com/jslogin?appid=wx782c26e4c19acffb&redirect_uri=https%3A%2F%2Fwx2.qq.com%2Fcgi-bin%2Fmmwebwx-bin%2Fwebwxnewloginpage&fun=new&lang=zh_CN&_=",
	"login_qrcode_api": "https://login.weixin.qq.com/qrcode/",
	"login_poll_api": "https://login.wx2.qq.com/cgi-bin/mmwebwx-bin/login?uuid={uuid}&tip={tip}&_={time}",
	"get_contact_api": "https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxgetcontact?lang=zh_CN&pass_ticket={pass_ticket}&seq=0&skey={skey}",
	"wx_init_api": "https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxinit?lang=zh_CN&pass_ticket={pass_ticket}",
}
