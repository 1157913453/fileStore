package user_service

import (
	"filestore/models"
	"filestore/service/cache_service"
	log "github.com/sirupsen/logrus"
)

func GetUserByPhone(phone string) (userInfo *models.User, err error) {
	userInfo, err = models.GetUserByPhone(phone)
	if err != nil {
		log.Errorf("查询%s用户失败：%v", phone, err)
		return nil, err
	}
	return
}

func CreateUser(phone, password, userName string) error {
	user := &models.User{
		Phone:    phone,
		Password: password,
		UserName: userName,
	}

	return models.CreateUser(user)
}

func CheckPassword(phone, encPassword string) error {
	return models.CheckUser(phone, encPassword)
}

// 通过reds或数据库查询用户
func GetUser(phone string) (userInfo *models.User, err error) {
	userInfo, err = cache_service.GetUserCache(phone)
	if err != nil { // 缓存没找到就到数据库找并更新缓存
		userInfo, err = GetUserByPhone(phone)
		if err != nil {
			return nil, err
		}
		err = cache_service.AddUserCache(userInfo)
		if err != nil {
			return nil, err
		}

	}
	return
}