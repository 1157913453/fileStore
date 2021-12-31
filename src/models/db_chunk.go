package models

type Chunk struct {
	FileMd5  string `gorm:"type:varchar(32);not null;" json:"fileMd5"`
	ChunkSum int32  `json:"chunkSum"`
}

func GetChunkSum(Md5 string) (chunk *Chunk, err error) {
	err = OrmDb.First(chunk, "file_md5 = ?", Md5).Error
	return
}

func DeleteFileChunk(Md5 string) (err error) {
	taskDetail := new(UploadTaskDetail)
	//var chunk *Chunk
	err = OrmDb.Where("file_md5 = ?", Md5).Delete(taskDetail).Error
	//if err != nil {
	//	return
	//}
	//err = OrmDb.Where("file_md5 = ?").Delete(chunk).Error
	return
}
