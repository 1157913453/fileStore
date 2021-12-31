package user_service

import (
	"filestore/src/models"
)

func GetUserByPhone(phone string) (*models.User, error) {
	return models.GetUserByPhone(phone)
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

func GetUserInfoByToken(token string) (*models.User, error) {
	return models.GetUserByToken(token)
}
