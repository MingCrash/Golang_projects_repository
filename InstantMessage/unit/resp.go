package unit

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type loginresp struct {
	Code  		int
	Msg   		string
	Data  		interface{}    	 //interface{} 在golang语言中表示任何类型
}


func RespSuccess(writer http.ResponseWriter,data interface{})  {
	Resp(0,"",data,writer)
}

func RespFail(writer http.ResponseWriter,err error)  {
	Resp(-1,err.Error(),nil,writer)
}

func Resp(code int,msg string,data interface{},writer http.ResponseWriter){
	h := &loginresp{
		Code:code,
		Msg:msg,
		Data:data,
	}
	b,err := json.Marshal(h)
	if err != nil {
		fmt.Println("json编码出错")
	}
	_, err = writer.Write(b)
	if err != nil {
		fmt.Println("写入器出错")
	}
}
