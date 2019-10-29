package model

import "time"

const (
	COMMUNITY_CATE_COM = 1
)
type Community struct {
	Id  	int64
	Name    string			//名称
	Ownerid 	int64		//群主ID
	Icon  		string		//群logo
	Cate		int
	Memo  		string		//群描述
	Createat	time.Time	//群创建时间
}