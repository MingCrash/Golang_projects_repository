package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const (
	success_status = 200
	notresourcefound_status = 403
)

type loginresp struct {
	Code  		int
	Msg   		string
	Data  		interface{}    	 //interface{} 在golang语言中表示任何类型
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

func userLogin(writer http.ResponseWriter, request *http.Request)  {
		writer.Header().Set("`Content`-Type","json/xml")
		//_ = request.ParseForm()
		//mobile := request.PostForm.Get("mobile")		//读取参数前需要解析
		//passwd := request.PostForm.Get("passwd")

		mobile := request.PostFormValue("mobile")	//调用时，已自动解析参数
		passwd := request.PostFormValue("passwd")

		if (mobile == "13697413574" && passwd == "123456"){
			writer.WriteHeader(success_status)
			Resp(success_status,"Login success",`{"id":100001,"token":"test"}`,writer)
		}else {
			writer.WriteHeader(notresourcefound_status)
 			Resp(notresourcefound_status,"Not resource found",`{"id":none,"token":"test"}`,writer)
		}
}

func UserTypeView()  {
	tpl, err := template.ParseGlob("./**/**/*")   //ParseGlob
	if err != nil {log.Fatal("template ParseFiles Failture" + err.Error())}
	for _,v := range tpl.Templates() {
		tplname := v.Name()
		//注册handler
		http.HandleFunc(tplname, func(writer http.ResponseWriter, request *http.Request) {
				err = v.ExecuteTemplate(writer, tplname, nil)
				if err != nil {log.Fatal("Execute Template Failture"+err.Error())}
		})
	}
}

func main() {
	//处理 通过API访问 的函数
	http.HandleFunc("/user/login",userLogin)

	//提供静态资源目录支持
	http.Handle("/asset/",http.FileServer(http.Dir(("."))))  //提供静态资源的目录地址 = http.Dir + pattern

	//使用template模板渲染
	//http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
	//	//template解析
	//	tpl, err := template.ParseFiles("./view/user/login.html")
	//	if err != nil {
	//		//打印并直接退出
	//		log.Fatal("template ParseFiles Failture" + err.Error())
	//	}
	//	err = tpl.ExecuteTemplate(w, "/user/login.shtml", nil)
	//	if err != nil {
	//		//打印并直接退出
	//		log.Fatal("Execute Template Failture"+err.Error())
	//	}
	//})
	////注册
	//http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
	//	//template解析
	//	tpl, err := template.ParseFiles("./view/user/register.html")
	//	if err != nil {
	//		//打印并直接退出
	//		log.Fatal("template ParseFiles Failture" + err.Error())
	//	}
	//	err = tpl.ExecuteTemplate(w, "/user/register.shtml", nil)
	//	if err != nil {
	//		//打印并直接退出
	//		log.Fatal("Execute Template Failture"+err.Error())
	//	}
	//})

	UserTypeView()

	_ = http.ListenAndServe(":9090", nil)
}


