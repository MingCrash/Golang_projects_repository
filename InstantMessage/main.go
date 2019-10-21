package main

import (
	"encoding/json"
	_ "golang.org/x/net/websocket"
	_ "io"
	"net/http"
)

func main() {
	//绑定请求与处理函数
	http.HandleFunc("/user/hell", func(writer http.ResponseWriter, request *http.Request) {

		str :=
		//获取request参数
		mobile := request.PostForm.Get("mobile")
		password := request.PostForm.Get("password")
		_, _ = writer.Write([]byte(mobile+password))
		writer.WriteHeader(200)
		writer.Header().Set("Content_type","json/xml")	})



	_ = http.ListenAndServe(":4545", nil) //阻塞
}
