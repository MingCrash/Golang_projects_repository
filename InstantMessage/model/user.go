package model

import "time"

const (
	SEX_WOMEN = "W"
	SEX_MAN = "M"
	SEX_UNKNOWN = "U"
)

type User struct {
	Id    		int64  		`json:"id"`		//用户ID
	Mobile  	string						//手机号码
	Passwd		string						//用户密码=f(plainwd+salt)
	Avatar 		string
	Sex 		string
	Nickname	string
	Salt		string						//随机数
	Online		bool
	Token   	string						//chat?id=1&token=x
	Memo        string
	Createat	time.Time					//统计每天用户增量
}
