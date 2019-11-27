package controller

import (
	"../model"
	"../unit"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"log"
	"net/http"
	"strconv"
	"sync"
)

//chat服务的核心就是将userid与ConnNode 形成映射关系
//定义client = map[userId][ConnNode]
var clientMap map[int64]*ConnNode = make(map[int64]*ConnNode)
var rwlocker sync.RWMutex

type ConnNode struct {
	Conn *websocket.Conn `json:"conn"`
	DataQueue chan []byte   `json:"data_queue"`			//将结点接收到的可能性并行信息转换成串行信息,
	KillQueue chan bool 	`json:"kill_queue"`
	GroupSet  set.Interface `json:"group_set"`
}

//	127.0.0.1/chat?id=xxx&token=xxx
func Chat(writer http.ResponseWriter, request *http.Request)  {
	//检查接入是否合法
	query := request.URL.Query()			//读取url后带参数
	id := query.Get("id")
	toverified_token := query.Get("token")
	toverified_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		unit.RespFail(writer,err)
		return
	}

	isvalid:= CheckToken(toverified_id,toverified_token)
	if isvalid {
		//golang判断key是否在map中
		if _, ok := clientMap[toverified_id]; ok{
			unit.RespSuccess(writer,"该UserId已加入即时对话名单中")
			return
		}
	}

	//把 http 请求升级为长连接的 WebSocket
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isvalid
		},
	}).Upgrade(writer, request, nil)
	if err != nil{
		unit.RespFail(writer,err)
		return
	}

	//生成conn连接结点对象
	var newnode = &ConnNode{
		Conn:      		conn,
		DataQueue: 		make(chan []byte,30),
		KillQueue:      make(chan bool),
		GroupSet: 		set.New(set.ThreadSafe),
	}
	//建立映射关系
	rwlocker.Lock()
	clientMap[toverified_id] = newnode   //map同时read不会引发异常，同时read和write会异常，同时write会异常,Chat服务可能同时被调用，所以需要加读写锁
	rwlocker.Unlock()

	//启动目前用户的接收器协程与发送器协程
	go SendCoro(newnode)
	go RecvCoro(newnode)

	conn.SetCloseHandler(func(code int, text string) error {

		newnode.KillQueue <- true


		rwlocker.Lock()
		delete(clientMap, toverified_id)
		rwlocker.Unlock()
		return error()
	})
}

func SendCoro(node *ConnNode)  {
	for {
		select {
			case data := <- node.DataQueue:
				err := node.Conn.WriteMessage(websocket.TextMessage, data)
				if err != nil{
					log.Println(err.Error())
					return
				}
			case <- node.KillQueue:
				return
		}
	}
}

func RecvCoro(node *ConnNode)  {
	for {
		select {
			case <- node.KillQueue:
				return
			default:
				_,data,err := node.Conn.ReadMessage()
				Dispach(&data,node)
				if err != nil{
					log.Println(err.Error())
					return
				}
			}
		}
}

//解析接收到的json信息
func Dispach(data *[]byte,selfnode *ConnNode)  {
	msg := model.Message{}
	err := json.Unmarshal(*data,&msg)
	fmt.Println(msg)
	if err != nil{
		log.Println(err.Error())
		return
	}
	switch msg.Cmd{
		//判断为私聊类型信息，根据目标id,转发消息到目标connNode的dataqueue信息通道上
		case model.CMD_SINGLE_MSG:
			TranferMsgto(msg.Dstid,data)
		//判断为群聊类型信息，根据目标群id,转发消息到目标connNode的dataqueue信息通道上
		case model.CMD_ROOM_MSG:
			fmt.Println("接收到群聊类型消息")
		//判断为心跳类型信息，
		case model.CMD_HEART:
			Pong(selfnode)
	}
}

func Pong(selfnode *ConnNode)  {
	msg := model.Message{}
	msg.Content = "pong"
	msg.Cmd = model.CMD_HEART
	msg.Media = model.MEDIA_TYPE_TEXT
	bytelist, err := json.Marshal(msg)
	if err!=nil {
		fmt.Println(err.Error())
	}else{
		selfnode.DataQueue <- bytelist
	}
}

func TranferMsgto(distid int64, data *[]byte) {
	rwlocker.RLock()
	distUserNode,ok := clientMap[distid]
	rwlocker.RUnlock()
	if ok {
		distUserNode.DataQueue <- *data
	}else{
		fmt.Println("Warning：消息转发失败，在用户列表中找不到对应的目标Id")
	}
}

func CheckToken(userId int64,token string) (bool) {
	user ,_ := userService.FindUserBy(userId)
	if user != nil && user.Token != token {
		return false
	}
	return true
}


