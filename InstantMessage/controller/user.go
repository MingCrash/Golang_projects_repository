package controller

import (
	"../model"
	"../service"
	"../unit"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
)
const (
	successStatus           = 200
	notresourcefound_status = 403
)

var userService  *service.UserService
var contactService  *service.ContactService

func UserRegister(writer http.ResponseWriter, request *http.Request)  {
	writer.Header().Set("Content-Type","json/xml")
	//_ = request.ParseForm()
	//mobile := request.PostForm.Get("mobile")		//读取参数前需要解析
	//passwd := request.PostForm.Get("passwd")

	mobile := request.PostFormValue("mobile")	//调用时，已自动解析参数
	plainpwd := request.PostFormValue("passwd")
	nickname := fmt.Sprintf("user%06d",rand.Int31())
	avatar := ""
	sex := model.SEX_UNKNOWN

	newuser, err := userService.Register(mobile, plainpwd, nickname, avatar, sex)
	if err != nil{
		writer.WriteHeader(successStatus)
		unit.RespFail(writer,err)
	}else{
		writer.WriteHeader(successStatus)
		unit.RespSuccess(writer,newuser)
	}
}

func UserLogin(writer http.ResponseWriter, request *http.Request)  {
	writer.Header().Set("Content-Type","json/xml")
	//_ = request.ParseForm()
	//mobile := request.PostForm.Get("mobile")		//读取参数前需要解析
	//passwd := request.PostForm.Get("passwd")

	mobile := request.PostFormValue("mobile")	//调用时，已自动解析参数
	plainpwd := request.PostFormValue("passwd")
	var user, err = userService.Login(mobile,plainpwd)
	if err!=nil{
		unit.RespFail(writer,err)
	}else {
		unit.RespSuccess(writer,user)
	}
}

func ContactAddfriend(writer http.ResponseWriter, request *http.Request)  {
	writer.Header().Set("Content-Type","json/xml")
	//绑定对象（将接收到不同类型的数据归一化格式）
	var args model.Args
	//err := unit.Bind(request, &args)
	v, _ := ioutil.ReadAll(request.Body)
	fmt.Println(string(v))
	err := json.Unmarshal(v, &args)
	if err != nil{fmt.Println(err.Error())}
	fmt.Println(args)
	//fmt.Println(fmt.Sprintf("reqBody:%s reqForm:%s args:%s", request.Body, request.PostForm, args))

	err = contactService.Addfriend(args.UserId, args.DistId)
	if err!=nil{
		unit.RespFail(writer,err)
	}else {
		unit.RespSuccess(writer,nil)
	}
}
