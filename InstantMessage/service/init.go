package service

import (
	"../model"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
)
var  DBEngine  *xorm.Engine
//init函数不能被其他函数调用，在main之前执行
func init()  {
	driverName := "mysql"
	dataSourceName := "root:zzm441538235@(127.0.0.1:3306)/user_account?charset=utf8"
	var err = errors.New("")
	DBEngine,err = xorm.NewEngine(driverName,dataSourceName)
	if err != nil {
		log.Fatal(err.Error())
	}
	//设置DBEngine
	DBEngine.ShowSQL(true)
	DBEngine.ShowExecTime(true)
	DBEngine.SetMaxOpenConns(2)

	err = DBEngine.Sync2(new(model.User))
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("db init success")
}