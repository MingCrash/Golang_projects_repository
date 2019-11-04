package controller

import (
	"../model"
	"../unit"
	"net/http"
)


func ContactAddfriend(writer http.ResponseWriter, request *http.Request)  {
	writer.Header().Set("Content-Type","json/xml")
	var args model.Args
	err := unit.Bind(request, &args)		//绑定对象（将接收到不同类型的数据归一化格式）
	if err != nil{
		unit.RespFail(writer,err)
		return
	}
	err = contactService.Addfriend(args.UserId, args.DistId)
	if err!=nil{
		unit.RespFail(writer,err)
	}else {
		unit.RespSuccess(writer,nil)
	}
}
