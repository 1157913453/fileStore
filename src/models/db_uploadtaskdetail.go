package models

type UploadTaskDetail struct {
	Id           uint   `gorm:"primarykey" json:"id"`
	FileMd5      string `gorm:"type:char(32);not null;" json:"md5"`
	ChunkNumber  int64  `gorm:"not null;default:1" json:"chunk_number"`
	ChunkSize    int64  `gorm:"not null;default:1048576" json:"chunkSize"`
	TotalChunks  int32  `gorm:"not null;default:1" json:"totalChunks"`
	TotalSize    int64  `gorm:"not null;default:1" json:"totalSize"`
	FilePath     string `gorm:"type:varchar(1024);not null;default:'/" json:"filePath"`
	RelativePath string `gorm:"type:varchar(1024);not null;default:'/" json:"relativePath"`
}

func QueryUploadTaskDetails(Md5 string) (details []*UploadTaskDetail, err error) {
	err = OrmDb.Find(&details, "file_md5 = ?", Md5).Error
	return
}

func CreateUploadTaskDetail(detail *UploadTaskDetail) (err error) {
	err = OrmDb.Create(detail).Error
	return
}
