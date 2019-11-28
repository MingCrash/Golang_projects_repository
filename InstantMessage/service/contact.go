package service

import (
	"../model"
	"github.com/pkg/errors"
	"time"
)

type ContactService struct {
}

func (us *ContactService) Addfriend(userId,distId int64) (err error) {
	//判断是否自己加自己
	if userId == distId{
		return errors.New("不能够加自己为好友")
	}

	//判断distId是否已注册用户
	var tmpuser model.User
	_, err = DBEngine.ID(distId).Get(&tmpuser)
	if !(tmpuser.Id>0)	{
		return errors.New("目标用户非注册用户")
	}

	//查询是否已经是好友
	//链式操作
	var tmpcontact model.Contact
	//_, err = DBEngine.Where("ownerid=?", userId).And("dstodj=?", distId).And("cate",model.CONCAT_CATE_USER).Get(&tmpcontact)
	_, err = DBEngine.Where("ownerid=? and dstodj= ? and cate= ?", userId, distId,model.CONCAT_CATE_USER).Get(&tmpcontact)
	//如果存在记录说明已经是好友了不加
	if tmpcontact.Id>0	{
		return errors.New("目标用户已是你的好友")
	}

	//添加好友关系
	//多数据操作都需要成功时候，需要使用事务来确保所有操作均成功后再做提交操作
	//在有需要的批量操作数据时，事务的使用往往是必要的
	session := DBEngine.NewSession()
	_ = session.Begin()
	defer session.Close()
	_, err1 := session.InsertOne(model.Contact{
		Ownerid:  userId,
		Dstodj:   distId,
		Cate:     model.CONCAT_CATE_USER,
		Memo:     "",
		Createat: time.Now(),
	})

	_, err2 := session.InsertOne(model.Contact{
		Ownerid:  distId,
		Dstodj:   userId,
		Cate:     model.CONCAT_CATE_USER,
		Memo:     "",
		Createat: time.Now(),
	})
	if err1 == nil && err2 == nil {
		_ = session.Commit()
	}else {
		//当在执行事务过程中遇到任何错误时，应该及时停止事务，将已经执行的进行回滚
		_ = session.Rollback()
		if err1 != nil{
			return err1
		}else {
			return err2
		}
	}

	//返回数据
	return nil
}

func (us *ContactService) Joincommunity()  {

}

func (us *ContactService) Loadfriend(userId int64) (*[]model.Contact,error) {
	tmpcontact := make([]model.Contact,0)
	//链式操作
	err := DBEngine.Where("ownerid=?", userId).Find(&tmpcontact)
	if err != nil {
		return nil,err
	}else if len(tmpcontact) == 0{
		return nil, errors.New("用户没有任何好友")
	}
	//返回数据
	return &tmpcontact,nil
}

