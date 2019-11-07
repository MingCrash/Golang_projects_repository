package model

const (
	CMD_SINGLE_MSG = 10	//私聊信息
	CMD_ROOM_MSG = 11	//群聊信息
	CMD_HEART = 0		//心跳信息，不处理
)

const (
	MEDIA_TYPE_TEXT = 1
	MEDIA_TYPE_NEWS = 2
	MEDIA_TYPE_VOICE = 3
	MEDIA_TYPE_IMG = 4
	MEDIA_TYPE_REDPACKGR = 5
	MEDIA_TYPE_EMOJ = 6
	MEDIA_TYPE_LINK = 7
	MEDIA_TYPE_VIDEO = 8
	MEDIA_TYPE_CONCAT = 9
	MEDIA_TYPE_UDEF = 100
)

type Message struct {
	Id      int64  `json:"id,omitempty" form:"id"` //消息ID
	Userid  int64  `json:"userid,omitempty" form:"userid"` //谁发的
	Cmd     int    `json:"cmd,omitempty" form:"cmd"` //群聊还是私聊
	Dstid   int64  `json:"dstid,omitempty" form:"dstid"`//对端用户ID/群ID
	Media   int    `json:"media,omitempty" form:"media"` //消息按照什么样式展示
	Content string `json:"content,omitempty" form:"content"` //消息的内容
	Pic     string `json:"pic,omitempty" form:"pic"` //预览图片
	Url     string `json:"url,omitempty" form:"url"` //服务的URL
	Memo    string `json:"memo,omitempty" form:"memo"` //简单描述
	Amount  int    `json:"amount,omitempty" form:"amount"` //其他和数字相关的
}

//消息发送结构体,点对点单聊为例
//1、MEDIA_TYPE_TEXT
//{id:1,userid:2,dstid:3,cmd:10,media:1,
//content:"hello"}
//
//3、MEDIA_TYPE_VOICE,amount单位秒
//{id:1,userid:2,dstid:3,cmd:10,media:3,
//url:"http://www.a,com/dsturl.mp3",
//amount:40}
//
//4、MEDIA_TYPE_IMG
//{id:1,userid:2,dstid:3,cmd:10,media:4,
//url:"http://www.baidu.com/a/log.jpg"}
//
//
//2、MEDIA_TYPE_News
//{id:1,userid:2,dstid:3,cmd:10,media:2,
//content:"标题",
//pic:"http://www.baidu.com/a/log,jpg",
//url:"http://www.a,com/dsturl",
//"memo":"这是描述"}
//
//
//5、MEDIA_TYPE_REDPACKAGR //红包amount 单位分
//{id:1,userid:2,dstid:3,cmd:10,media:5,url:"http://www.baidu.com/a/b/c/redpackageaddress?id=100000","amount":300,"memo":"恭喜发财"}

//6、MEDIA_TYPE_EMOJ 6
//{id:1,userid:2,dstid:3,cmd:10,media:6,"content":"cry"}
//
//7、MEDIA_TYPE_Link 7
//{id:1,userid:2,dstid:3,cmd:10,media:7,
//"url":"http://www.a.com/dsturl.html"
//}
//
//8、MEDIA_TYPE_VIDEO 8
//{id:1,userid:2,dstid:3,cmd:10,media:8,
//pic:"http://www.baidu.com/a/log,jpg",
//url:"http://www.a,com/a.mp4"
//}
//
//9、MEDIA_TYPE_CONTACT 9
//{id:1,userid:2,dstid:3,cmd:10,media:9,
//"content":"10086",
//"pic":"http://www.baidu.com/a/avatar,jpg",
//"memo":"胡大力"}