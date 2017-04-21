package gochat

var Config = map[string]string {
	"get_uuid_api": 		"https://login.wx2.qq.com/jslogin?appid=wx782c26e4c19acffb&redirect_uri=https%3A%2F%2Fwx2.qq.com%2Fcgi-bin%2Fmmwebwx-bin%2Fwebwxnewloginpage&fun=new&lang=zh_CN&_=",
	"login_qrcode_api": 	"https://login.weixin.qq.com/qrcode/",
	"login_poll_api": 		"https://login.wx2.qq.com/cgi-bin/mmwebwx-bin/login?uuid={uuid}&tip={tip}&_={time}",
	"get_contact_api": 		"https://{host}/cgi-bin/mmwebwx-bin/webwxgetcontact?lang=zh_CN&pass_ticket={pass_ticket}&seq=0&skey={skey}",
	"wx_init_api":			"https://{host}/cgi-bin/mmwebwx-bin/webwxinit?r={r}&lang=zh_CN&pass_ticket={pass_ticket}",
	"wx_statusnotify_api": 	"https://{host}/cgi-bin/mmwebwx-bin/webwxstatusnotify?lang=zh_CN&pass_ticket={pass_ticker}",
	"getcontact_api": 		"https://{host}/cgi-bin/mmwebwx-bin/webwxgetcontact?lang=zh_CN&pass_ticket={pass_ticket}&r={r}&seq={seq}&skey={skey}",
	"batchgetcontact_api": 	"https://{host}/cgi-bin/mmwebwx-bin/webwxbatchgetcontact?type=ex&r={r}",
	"sync_api": 			"https://{host}/cgi-bin/mmwebwx-bin/webwxsync?sid={sid}&skey={skey}",
	"synccheck_api": 		"https://{host}/cgi-bin/mmwebwx-bin/synccheck?r={r}&skey={skey}&sid={sid}&uin={uin}&deviceid={deviceid}&synckey={synckey}",
	"sendmsg_api": 			"https://{host}/cgi-bin/mmwebwx-bin/webwxsendmsg?lang=zh_CN&pass_ticket={pass_ticket}",
	// 同意好友申请
	"verify_user_api": 		"https://{host}/cgi-bin/mmwebwx-bin/webwxverifyuser?r={r}&pass_ticket={pass_ticket}",
	// 上传文件
	"upload_media_api": 	"https://{prefix}.{host}/cgi-bin/mmwebwx-bin/webwxuploadmedia?f=json",
	// 发送文件
	"sendappmsg_api":		"https://{host}/cgi-bin/mmwebwx-bin/webwxsendappmsg?fun=async&f=json",
	// 发送图片
	"sendimgmsg_api":		"https://{host}/cgi-bin/mmwebwx-bin/webwxsendmsgimg?fun=async&f=json",
	// 暂时没用到
	"statreport_api": 		"https://{host}/cgi-bin/mmwebwx-bin/webwxstatreport?fun=new&lang=zh_CN&pass_ticket={pass_ticket}",
}