package file_service

import "filestore/src/models"

func GetFileList(phone, fileType, path string, page, pageCount int) (models.FileList, error) {
	return models.GetFileList(phone, fileType, path, (page-1)*pageCount, pageCount)
}

func CreateUserFile(file *models.UserFile) error {
	return models.CreateUserFile(file)
}

func GetFolder(path, name string) error {
	return models.GetFolder(path, name)
}

func GetUserFile(Md5 string) (*models.UserFile, error) {
	return models.GetUserFile(Md5)
}

func BatchDeleteFile(files models.DeleteFiles) error {
	return models.BatchDeleteFile(files.Files)
}

func GetUserFileById(id int) (*models.UserFile, error) {
	return models.GetUserFileById(id)
}
