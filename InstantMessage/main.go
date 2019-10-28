package main

import (
	"./controller"
	"html/template"
	"log"
	"net/http"
)

//注册视图
func UserTypeView()  {
	tpl, err := template.ParseGlob("./view/user/*.html")   //ParseGlob找到符合
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
	http.HandleFunc("/user/login",controller.UserLogin)
	http.HandleFunc("/user/register",controller.UserRegister)

	//提供静态资源目录支持,js,css等文件引用就靠这个了
	http.Handle("/asset/",http.FileServer(http.Dir(("."))))  //提供静态资源的目录地址 = http.Dir + pattern

	UserTypeView()

	_ = http.ListenAndServe(":9090", nil)
}


