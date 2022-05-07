package file_service

import "filestore/models"

func QueryUploadTaskDetails(Md5 string) ([]*models.UploadTaskDetail, error) {
	return models.QueryUploadTaskDetails(Md5)
}

func CreateUploadTaskDetail(detail *models.UploadTaskDetail) error {
	return models.CreateUploadTaskDetail(detail)
}
