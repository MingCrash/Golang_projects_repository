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
	Conn *websocket.Conn 	`json:"conn"`
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
		KillQueue:      make(chan bool,1),
		GroupSet: 		set.New(set.ThreadSafe),
	}
	//建立映射关系
	//如果设置了一个写锁，那么其它读的线程以及写的线程都拿不到锁，这个时候，与互斥锁的功能相同
	//如果设置了一个读锁，那么其它写的线程是拿不到锁的，但是其它读的线程是可以拿到锁
	rwlocker.Lock()
	clientMap[toverified_id] = newnode   //map同时read不会引发异常，同时read和write会异常，同时write会异常,Chat服务可能同时被调用，所以需要加读写锁
	rwlocker.Unlock()

	//启动目前用户的接收器协程与发送器协程
	go SendCoro(newnode)
	go RecvCoro(newnode)

	conn.SetCloseHandler(func(code int, text string) error {
		//给 管道KillQueue 发送信号关闭收发器
		newnode.KillQueue <- true
		//清理map无用node
		rwlocker.Lock()
		clientMap[toverified_id] = nil
		rwlocker.Unlock()
		return nil
	})
}

func SendCoro(node *ConnNode)  {
	for {
		select {
			case <- node.KillQueue:
				return
			case data := <- node.DataQueue:
				err := node.Conn.WriteMessage(websocket.TextMessage, data)
				if err != nil{
					log.Println(err.Error())
					return
				}
		}
	}
}

func RecvCoro(node *ConnNode)  {
	for {
			_,data,err := node.Conn.ReadMessage()
			Dispach(&data,node)
			if err != nil{
				log.Println(err.Error())
				return
			}
		}
}

//解析接收到的json信息
func Dispach(data *[]byte,selfnode *ConnNode)  {
	msg := model.Message{}
	err := json.Unmarshal(*data,&msg)
	if err != nil{
		log.Println(err.Error())
		return
	}
	switch msg.Cmd{
		//判断为私聊类型信息，根据目标id,转发消息到目标connNode的dataqueue信息通道上
		case model.CMD_SINGLE_MSG:
			rwlocker.RLock()						//设置写锁
			distUserNode,exists := clientMap[msg.Dstid]
			rwlocker.RUnlock()
			if exists {
				distUserNode.DataQueue <- *data
			}else{
				log.Println(fmt.Sprintf("[UserId:%v DistId:%v] msg:消息转发失败，在用户列表中找不到对应的目标Id",msg.Userid,msg.Dstid))
			}
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


func CheckToken(userId int64,token string) (bool) {
	user ,_ := userService.FindUserBy(userId)
	if user != nil && user.Token != token {
		return false
	}
	return true
}


