package model

import "time"

const (
	CONCAT_CATE_USER = 1
	CONCAT_CATE_COMUNITY = 2
)
//用户关系表
//好友和群都存在这个表里面
type Contact struct {
	Id 			int64
	Ownerid		int64
	Dstodj		int64
	Cate		int
	Memo        string
	Createat 	time.Time
}

type Args struct {
	UserId  	int64
	DistId		int64
}