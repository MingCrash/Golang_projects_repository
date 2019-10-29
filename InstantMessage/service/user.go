package service

import (
	"../model"
	"../unit"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type UserService struct {
}

func (us *UserService) Register(mobile,plainpwd,nickname,avatar,sex string) (user model.User,err error) {
	//检测手机号码是否存在,
	tmp := model.User{}
	_,err = DBEngine.Where("mobile=?",mobile).Get(&tmp)
	if err!=nil {
		return tmp,err
	}

	//如果存在则返回提示已经注册
	if tmp.Id > 0{
		return tmp,errors.New("该手机号已经注册")
	}

	//否则拼接并插入数据 (Id由数据库自动添加)
	tmp.Mobile = mobile
	tmp.Avatar = avatar
	tmp.Nickname = nickname
	tmp.Sex = sex
	tmp.Createat = time.Now()
	tmp.Salt = fmt.Sprintf("%06d",rand.Int31n(10000))
	tmp.Passwd = unit.MakePasswd(plainpwd,tmp.Salt)
	_, err = DBEngine.InsertOne(&tmp)

	//最后返回用户信息
	return tmp, err
}

func (us *UserService) Login(mobile,plainpwd string) (user model.User,err error) {
	//首先通过手机号查询用户
	tmp := model.User{}
	 _, err = DBEngine.Where("mobile = ?", mobile).Get(&tmp)
	if err!=nil {
		return tmp,err
	}
	//查询到对比密码
	if !unit.ValidatePasswd(plainpwd,tmp.Salt,tmp.Passwd){
		return tmp,errors.New("账号密码不正确")
	}
	//刷新token
	tmp.Token = unit.MD5Encode(strconv.Itoa(int(time.Now().Unix())))
	tmp.Online = true
	_, _ = DBEngine.ID(tmp.Id).Cols("token,online").Update(&tmp)

	//返回数据
	return tmp, nil
}