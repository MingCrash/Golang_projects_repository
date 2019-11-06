package model

import "time"

const (
	COMMUNITY_CATE_COM = 1
)
type Community struct {
	Id  		int64             `json:"id"`		
	Name    	string          `json:"name"`//名称
	Ownerid 	int64            `json:"ownerid"`//群主ID
	Icon  		string 			`json:"icon"`//群logo
	Cate		int             `json:"cate"`
	Memo  		string       	`json:"memo"`//群描述
	Createat	time.Time 		`json:"createat"`//群创建时间
}