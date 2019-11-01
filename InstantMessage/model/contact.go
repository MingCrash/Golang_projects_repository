package model

import "time"

const (
	CONCAT_CATE_USER = 1
	CONCAT_CATE_COMUNITY = 2
)
//用户关系表
//好友和群都存在这个表里面
type Contact struct {
	Id 			int64 		`json:"id"`
	Ownerid		int64 		`json:"ownerid"`
	Dstodj		int64 		`json:"dstodj"`
	Cate		int 		`json:"cate"`
	Memo        string 		`json:"memo"`
	Createat 	time.Time	`json:"createat"`
}

type Args struct {
	UserId  	int64 		`json:"userid" form:"userid"`
	DistId		int64 		`json:"distid"	form:"distid"`
}