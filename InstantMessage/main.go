package main

import (
	"encoding/json"
	"fmt"
	"./controller"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type test struct {
	UserId  	int64 		`json:"user_id" form:"userid"`
	DistId		int64 		`json:"dist_id"	form:"distid"`
}

//注册视图
func RegistetViews()  {
	tpl, err := template.ParseGlob("./view/**/*.html")   //ParseGlob找到符合
	if err != nil {log.Fatal("template ParseFiles Failture" + err.Error())}
	for _,v := range tpl.Templates() {
		tplname := v.Name()
		//注册handler
		http.HandleFunc(tplname, func(writer http.ResponseWriter, request *http.Request) {
				if strings.Contains(tplname,"/chat/index") {

				}
				err = v.ExecuteTemplate(writer, tplname, nil)
				if err != nil {log.Fatal("Execute Template Failture"+err.Error())}
		})
	}
}

func main(){
	//处理 通过API访问 的函数
	http.HandleFunc("/user/login",controller.UserLogin)
	http.HandleFunc("/user/logout",controller.UserLogout)
	http.HandleFunc("/user/register",controller.UserRegister)
	http.HandleFunc("/user/find",controller.UserFind)
	http.HandleFunc("/contact/addfriend",controller.ContactAddfriend)
	http.HandleFunc("/contact/deletefriend",controller.ContactDeletefriend)

	//http.HandleFunc("/contact/createcommunity",controller.ContactCreatecommunity)
	//http.HandleFunc("/contact/joincommunity",controller.ContactJoincommunity)
	http.HandleFunc("/contact/friend",controller.ContactLoadFriend)
	http.HandleFunc("/attach/upload",controller.AttachUpload)
	http.HandleFunc("/chat",controller.Chat)

	http.HandleFunc("/testing/url", func(writer http.ResponseWriter, request *http.Request) {
		var tmp test
		v, _ := ioutil.ReadAll(request.Body)
		fmt.Println(string(v))
		err := json.Unmarshal(v, &tmp)
		if err != nil{fmt.Println(err.Error())}
		fmt.Println(tmp)
	})

	//提供静态资源目录支持,js,css等文件引用就靠这个了
	http.Handle("/asset/",http.FileServer(http.Dir(("."))))  //提供静态资源的目录地址 = http.Dir + pattern
	http.Handle("/mnt/",http.FileServer(http.Dir(("."))))  //提供静态资源的目录地址 = http.Dir + pattern

	//模板引擎将html与数据结合，并注册（需要静态资源目录 支持）
	RegistetViews()

	_ = http.ListenAndServe(":9090", nil)
}
