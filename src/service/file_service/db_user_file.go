package file_service

import (
	"filestore/src/models"
	"filestore/src/service/user_service"
)

func GetFileList(phone, fileType, path string, page, pageCount int) (models.FileList, error) {
	userInfo, err := user_service.GetUser(phone)
	if err != nil {
		return models.FileList{}, err
	}
	return models.GetFileList(fileType, path, (page-1)*pageCount, pageCount, userInfo)
}

func CreateUserFile(file *models.UserFile) error {
	return models.CreateUserFile(file)
}

func GetFolder(userId uint, path, name string) error {
	return models.GetFolder(userId, path, name)
}

func GetUserFile(userId uint, Md5 string) (*models.UserFile, error) {
	return models.GetUserFile(userId, Md5)
}

func BatchDeleteFile(files models.DeleteFiles) error {
	return models.BatchDeleteFile(files.Files)
}

func GetUserFileById(id int) (*models.UserFile, error) {
	return models.GetUserFileById(id)
}
