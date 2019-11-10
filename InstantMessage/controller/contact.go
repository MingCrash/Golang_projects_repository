package controller

import (
	"../unit"
	"net/http"
	"strconv"
	_ "xorm.io/core"
)


func ContactAddfriend(writer http.ResponseWriter, request *http.Request)  {
	writer.Header().Set("Content-Type","json/xml")

	//var args model.Args
	//err := unit.Bind(request, &args)		//绑定对象（将接收到不同类型的数据归一化格式）
	//if err != nil{
	//	unit.RespFail(writer,err)
	//	return
	//}
	//err = contactService.Addfriend(args.UserId, args.DistId)

	userid := request.PostFormValue("userid")
	distid := request.PostFormValue("distid")
	uid, _ := strconv.ParseInt(userid,10,64)
	did, _ := strconv.ParseInt(distid,10,64)
	err := contactService.Addfriend(uid, did)

	if err!=nil{
		unit.RespFail(writer,err)
	}else {
		unit.RespSuccess(writer,nil)
	}
}

func ContactLoadFriend(writer http.ResponseWriter, request *http.Request)  {
	writer.Header().Set("Content-Type","json/xml")

	userid := request.PostFormValue("userid")
	intid, _ := strconv.ParseInt(userid,10,64)
	//根据输入的id查询contact
	frids, err := contactService.Loadfriend(intid)

	if err != nil {
		unit.RespFail(writer,err)
	}else{
		unit.RespSuccess(writer,*frids)
	}
}

