package controller

import (
	"../unit"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)


//	127.0.0.1/chat?id=xxx&token=xxx
func Chat(writer http.ResponseWriter, request *http.Request)  {
	//检查接入是否合法
	//toverified_id := request.PostFormValue("id")
	//toverified_token := request.PostFormValue("token")
	query := request.URL.Query()
	id := query.Get("id")
	toverified_token := query.Get("token")
	toverified_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		unit.RespFail(writer,err)
	}

	_, _ = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			isvalid,_ := CheckToken(toverified_id,toverified_token)
			if isvalid == false {log.Println(err.Error())}
			return isvalid
		},
	}.Upgrade(writer, request, nil)
}

func CheckToken(userId int64,token string) (bool,error) {
	user ,err := userService.FindUserBy(userId)
	if user != nil && user.Token != token {
		return false,err
	}
	return true,nil

}

