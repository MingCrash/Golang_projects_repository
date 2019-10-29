package model

import "time"

const (
	CONCAT_CATE_USER = 1
	CONCAT_CATE_COMUNITY = 2
)
//通信表
type Contact struct {
	Id 			int64
	Ownerid		int64
	Dstodj		int64
	Cate		int
	Memo        string
	Createat 	time.Time
}