package models

import "time"

type UploadTask struct {
	Id           uint      `gorm:"primarykey"`
	ExtendName   string    `gorm:"type:varchar(32)" json:"extendName'"`
	FileName     string    `gorm:"type:varchar(64);not null;default:''" json:"fileName"`
	FilePath     string    `gorm:"type:varchar(1024);not null;default:'/'" json:"filePath"`
	Md5          string    `gorm:"type:varchar(32);uniqueIndex:idx_Md5;not null;'" json:"md5"`
	UploadStatus uint      `json:"uploadStatus"`
	UploadTime   time.Time `json:"uploadTime"`
	UserId       uint      `json:"userId"`
}
