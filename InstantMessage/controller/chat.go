package controller

import (
	"../unit"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

//	127.0.0.1/chat?id=xxx&token=xxx
func Chat(writer http.ResponseWriter, request *http.Request)  {
	//检查接入是否合法
	query := request.URL.Query()							//读取url后带参数
	id := query.Get("id")
	toverified_token := query.Get("token")
	toverified_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		unit.RespFail(writer,err)
		return
	}
	//fmt.Println(fmt.Sprintf("toverified_id:%d,toverified_token:%s", toverified_id, toverified_token))
	//unit.RespSuccess(writer,isvalid)

	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			isvalid,err := CheckToken(toverified_id,toverified_token)
			if err != nil {	fmt.Println(err.Error())}
			return isvalid
		},
	}).Upgrade(writer, request, nil)

}

func CheckToken(userId int64,token string) (bool,error) {
	user ,err := userService.FindUserBy(userId)
	if user != nil && user.Token != token {
		fmt.Println(user.Token)
		return false,err
	}
	return true,err

}


