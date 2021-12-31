package models

import (
	"errors"
	"gorm.io/gorm"
)

type File struct {
	FileMd5    string `gorm:"type:char(32);not null;uniqueIndex:idx_fileMd5;default:''" json:"fileMd5"`
	FileName   string `gorm:"type:varchar(256);not null;index" json:"fileName"`
	FileSize   uint64 `gorm:"not null;default:0" json:"fileSize"`
	FileAddr   string `gorm:"type:varchar(1024);not null;default:''" json:"fileAddr"`
	FileStatus int8   `gorm:"default:0" json:"fileStatus"`
	FileType   string `gorm:"type:varchar(8)" json:"fileType"`
	PointCount uint32 `gorm:"not null;default:0" json:"pointCount"`
	gorm.Model
}

type Folder struct {
	FoldName string `json:"foldName"`
	FilePath string `json:"filePath"`
}

func CreateFileMeta(file *File) (err error) {
	fileMeta, err := GetFileMeta(file.FileMd5)
	if err != nil && err.Error() == "record not found" {
		err = OrmDb.Create(file).Error
		return err
	}
	if fileMeta != nil {
		return errors.New("已有该文件")
	}
	return
}

func GetFileMeta(Md5 string) (file *File, err error) {
	err = OrmDb.First(&file, "file_md5 = ?", Md5).Error
	if err != nil {
		return nil, err
	}
	return file, nil
}

func GetFileById(id int) (file *File, err error) {
	file = new(File)
	err = OrmDb.First(file, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return file, nil

}
