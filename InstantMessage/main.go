package main

import (
	"net/http"
	"strconv"
)

const (
	success_status = 200
	notresourcefound_status = 403
)

type loginresp struct {
	Code  int
	Msg   string
	data  interface{}     //interface{} 在golang语言中表示任何类型
}



func userLogin(writer http.ResponseWriter, request *http.Request)  {
		writer.Header().Set("`Content`-Type","json/xml")
		//_ = request.ParseForm()						//读取参数前需要解析
		//mobile := request.PostForm.Get("mobile")
		//passwd := request.PostForm.Get("passwd")

		mobile := request.PostFormValue("mobile")	//调用时，已自动解析参数
		passwd := request.PostFormValue("passwd")

		var body string
		if (mobile == "13697413574" && passwd == "123456"){
			writer.WriteHeader(success_status)
			body = `{"code":`+strconv.Itoa(success_status)+`,"msg":"Login success","data":{"id":100001,"token":"test"}}`
		}else {
			writer.WriteHeader(notresourcefound_status)
			body = `{"code":`+strconv.Itoa(notresourcefound_status)+`,"msg":"Not resource found","data":{"id":none,"token":"test"}}`
		}
		_, _ = writer.Write([]byte(body))
}

func main() {
	http.HandleFunc("/user/login",userLogin)

	_ = http.ListenAndServe(":9090", nil)
}


