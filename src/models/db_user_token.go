package models

import log "github.com/sirupsen/logrus"

type UserToken struct {
	Phone     string `gorm:"type:varchar(64);not null;uniqueIndex:idx_phone" json:"phone"`
	UserToken string `gorm:"type:varchar(256);not null;" json:"userToken"`
}

func UpdateToken(phone, token string) (err error) {
	tokenModel := &UserToken{
		Phone:     phone,
		UserToken: token,
	}
	_, err = GetToken(phone)
	if err != nil { // 如果没token,就创建
		err = OrmDb.Create(tokenModel).Error
		if err != nil {
			log.Errorf("创建token失败：%v", err)
			return
		}
	} else { // 如果有token，就更新
		err = OrmDb.Model(&UserToken{}).Where("phone = ?", phone).Update("user_token", token).Error
		if err != nil {
			return
		}
	}
	return
}

func GetToken(phone string) (token *UserToken, err error) {
	token = &UserToken{}
	err = OrmDb.First(token, "phone = ?", phone).Error
	return
}
