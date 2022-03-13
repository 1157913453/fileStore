package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	UserName      string     `gorm:"type:varchar(64);not null;default:'';" json:"userName"`
	Password      string     `gorm:"type:varchar(256);not null;" json:"password"`
	Email         string     `gorm:"type:varchar(64);default:'';" json:"email"`
	Phone         string     `gorm:"type:varchar(64);uniqueIndex:idx_phone;default:''" json:"phone"`
	RoleId        int8       `json:"roleId"`
	LastLoginTime *time.Time `json:"lastLoginTime"`

	gorm.Model
}

func GetUserByEmail(email string) {

}

func CheckUser(phone, encPassword string) (err error) {
	user := User{}
	err = OrmDb.First(user, "phone = ? AND password = ?", phone, encPassword).Error
	return
}

func GetPwdByPhone(phone string) (Pwd string, err error) {
	user := new(User)
	err = OrmDb.First(user, "phone = ?", phone).Error
	return user.Password, err
}

func GetUserByPhone(phone string) (user *User, err error) {
	err = OrmDb.First(&user, "phone = ?", phone).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByToken(token string) (user *User, err error) {
	var userToken *UserToken
	// 获取token对应的Phone
	err = OrmDb.First(&userToken, "user_token = ?", token).Error
	if err != nil {
		return nil, err
	}

	// 获取用户信息
	user, err = GetUserByPhone(userToken.Phone)
	if err != nil {
		return nil, err
	}
	return
}

func CreateUser(user *User) (err error) {
	err = OrmDb.Create(&user).Error
	return
}
