package controller

import (
	"../model"
	"../unit"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
)
const (
	successStatus           = 200
	notresourcefound_status = 403
)


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
		unit.RespSuccess(writer,&newuser)
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
		unit.RespSuccess(writer,*user)
	}
}

func UserLogout(writer http.ResponseWriter, request *http.Request)  {
	writer.Header().Set("Content-Type","json/xml")
	//_ = request.ParseForm()
	//mobile := request.PostForm.Get("mobile")		//读取参数前需要解析
	//passwd := request.PostForm.Get("passwd")
	id := request.PostFormValue("id")	//调用时，已自动解析参数
	token := request.PostFormValue("token")
	toverified_id, _ := strconv.ParseInt(id, 10, 64)
	err := userService.Logout(toverified_id,token)
	if err!=nil{
		unit.RespFail(writer,err)
	}else {
		unit.RespSuccess(writer,nil)
	}
}


