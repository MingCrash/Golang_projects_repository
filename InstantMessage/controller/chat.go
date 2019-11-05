package controller

import (
	"../unit"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

type ConnNode struct {
	Conn 		*websocket.Conn 	`json:"conn"`
	//将结点接收到的可能性并行信息转换成串行信息,
	DataQueue	chan []byte 		`json:"data_queue"`
	setGroup
}

//本核心在于形成userid和Node的映射关系

//	127.0.0.1/chat?id=xxx&token=xxx
func Chat(writer http.ResponseWriter, request *http.Request)  {
	//检查接入是否合法
	query := request.URL.Query()
	id := query.Get("id")
	toverified_token := query.Get("token")
	toverified_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		unit.RespFail(writer,err)
	}
	isvalid, err := CheckToken(toverified_id,toverified_token)
	if isvalid == false || err != nil {log.Println(err.Error())}

	//把 http 请求升级为长连接的 WebSocket
	conn, err := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isvalid
		},
	}.Upgrade(writer, request, nil)

	//生成conn连接结点对象
	node := ConnNode{

	}

}

func CheckToken(userId int64,token string) (bool,error) {
	user ,err := userService.FindUserBy(userId)
	if user != nil && user.Token != token {
		return false,err
	}
	return true,nil
}

