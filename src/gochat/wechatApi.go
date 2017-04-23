package gochat

var weChatApi = map[string]string {
	"getUuidApi": 			"https://login.wx2.qq.com/jslogin?appid=wx782c26e4c19acffb&redirect_uri=https%3A%2F%2Fwx2.qq.com%2Fcgi-bin%2Fmmwebwx-bin%2Fwebwxnewloginpage&fun=new&lang=zh_CN&_=",
	"qrcodeApi": 			"https://login.weixin.qq.com/qrcode/",
	"loginApi": 			"https://login.wx2.qq.com/cgi-bin/mmwebwx-bin/login?uuid={uuid}&tip={tip}&_={time}",
	"initApi":				"https://{host}/cgi-bin/mmwebwx-bin/webwxinit?r={r}&lang=zh_CN&pass_ticket={pass_ticket}",
	"getContactApi": 		"https://{host}/cgi-bin/mmwebwx-bin/webwxgetcontact?lang=zh_CN&pass_ticket={pass_ticket}&r={r}&seq={seq}&skey={skey}",
	"batchGetContactApi": 	"https://{host}/cgi-bin/mmwebwx-bin/webwxbatchgetcontact?type=ex&r={r}",
	"syncApi": 				"https://{host}/cgi-bin/mmwebwx-bin/webwxsync?sid={sid}&skey={skey}",
	"syncCheckApi": 		"https://{host}/cgi-bin/mmwebwx-bin/synccheck?r={r}&skey={skey}&sid={sid}&uin={uin}&deviceid={deviceid}&synckey={synckey}&_={_}",
	"sendMsgApi": 			"https://{host}/cgi-bin/mmwebwx-bin/webwxsendmsg?lang=zh_CN&pass_ticket={pass_ticket}",
	"verifyUserApi": 		"https://{host}/cgi-bin/mmwebwx-bin/webwxverifyuser?r={r}&pass_ticket={pass_ticket}",
	"uploadMediaApi": 		"https://{prefix}.{host}/cgi-bin/mmwebwx-bin/webwxuploadmedia?f=json",
	"sendAppMsgApi":		"https://{host}/cgi-bin/mmwebwx-bin/webwxsendappmsg?fun=async&f=json",
	"sendImgMsgApi":		"https://{host}/cgi-bin/mmwebwx-bin/webwxsendmsgimg?fun=async&f=json",
	"logoutApi":			"https://{host}/cgi-bin/mmwebwx-bin/webwxlogout?redirect=1&type=1&skey={skey}", /* sid:IF7Nv4Z1ci0uRAgz uin:2926158633 Upgrade-Insecure-Requests:1*/
	"pushLoginApi":			"https://{host}/cgi-bin/mmwebwx-bin/webwxpushloginurl?uin={uin}",
}