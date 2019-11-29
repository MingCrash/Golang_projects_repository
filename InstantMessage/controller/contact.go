package controller

import (
	"../unit"
	"net/http"
	"strconv"
	"../model"
	"time"
)


func ContactAddfriend(writer http.ResponseWriter, request *http.Request)  {
	writer.Header().Set("Content-Type","json/xml")
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

func ContactJoincommunity(writer http.ResponseWriter, request *http.Request)  {

}

func ContactCreatecommunity(writer http.ResponseWriter, request *http.Request)  {
	ownerid, _ := strconv.ParseInt(request.PostFormValue("ownerid"),10,64)
	commname := request.PostFormValue("name")
	icon := request.PostFormValue("icon")
	memo := request.PostFormValue("memo")
	tmpcommunity := model.Community{
		Id:       ownerid,
		Name:     commname,
		Ownerid:  ownerid,
		Icon:     icon,
		Cate:     model.COMMUNITY_CATE_COM,
		Memo:     memo,
		Createat: time.Now(),
	}

	err := contactService.Createcommunity(&tmpcommunity)
	if err!=nil{
		unit.RespFail(writer,err)
	}else {
		unit.RespSuccess(writer,"")
	}
}


