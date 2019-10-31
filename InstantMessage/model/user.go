package model

import "time"

const (
	SEX_WOMEN = "W"
	SEX_MAN = "M"
	SEX_UNKNOWN = "U"
)

type User struct {
	Id    		int64  		`json:"id"`					//用户ID
	Mobile  	string      `json:"mobile"`				//手机号码
	Passwd		string 		`json:"passwd"`				//用户密码=f(plainwd+salt)
	Avatar 		string 		`json:"avatar"`
	Sex 		string 		`json:"sex"`
	Nickname	string      `json:"nickname"`
	Salt		string 		`json:"salt"`				//随机数
	Online		bool        `json:"online"`
	Token   	string 		`json:"token"`				//chat?id=1&token=x
	Memo        string      `json:"memo"`
	Createat	time.Time 	`json:"createat"`			//统计每天用户增量
}
