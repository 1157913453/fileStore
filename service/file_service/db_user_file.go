package file_service

import (
	"filestore/models"
	"filestore/service/oss_service"
	"filestore/service/user_service"
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

func BatchDelete(files models.DeleteRecoveryFiles) error {
	return models.BatchDelete(files.Files)
}

func PermanentlyDelete(files models.DeleteRecoveryFiles) error {
	DeleteOssFiles, err := models.PermanentlyDelete(files.Files)
	go oss_service.OssDeleteFiles(DeleteOssFiles)
	return err
}

func GetUserFileById(id int) (*models.UserFile, error) {
	return models.GetUserFileById(id)
}
