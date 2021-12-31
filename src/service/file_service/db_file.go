package file_service

import "filestore/src/models"

func CreateFileMeta(file *models.File) error {
	return models.CreateFileMeta(file)
}

func GetFileMeta(fileMd5 string) (*models.File, error) {
	return models.GetFileMeta(fileMd5)
}

func DeleteFileChunk(md5 string) error {
	return models.DeleteFileChunk(md5)
}

func GetRecoveryFileList() ([]*models.ListData, error) {
	return models.GetRecoveryFileList()
}

func GetFileById(id int) (*models.File, error) {
	return models.GetFileById(id)
}
