package models

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	p "path"
	"strings"
)

var FileIdMap = map[string]uint{} // key=Md5,val=id

type UserFile struct {
	UserId     uint   `gorm:"not null;index:idx_userId;" json:"userId"`
	FileMd5    string `gorm:"type:varchar(32);" json:"fileMd5"`
	DeleteFlag uint8  `gorm:"not null;default:0;index:idx_deleteFlag" json:"delete"`
	FileSize   uint64 `gorm:"not null;default:0" json:"fileSize"`
	FileName   string `gorm:"type:varchar(256);not null;" json:"fileName"`
	FilePath   string `gorm:"type:varchar(1024);not null;default:'/'" json:"filePath"`
	IsDir      uint8  `gorm:"not null;default:0" json:"isDir"`
	FileType   string `gorm:"type:varchar(8)" json:"fileType"`
	FileStatus uint8  `json:"fileStatus"`

	gorm.Model
}

type FileList struct {
	Total int64       `json:"total"`
	List  []*ListData `json:"list"`
}

type ListData struct {
	UserFileId uint   `json:"userFileId"`
	FileId     uint   `json:"fileId"`
	IsDeleted  uint8  `json:"deleteFlag"`
	ExtendName string `json:"extendName"`
	Md5        string `json:"identifier"`
	UserId     uint   `json:"userId"`
	FileSize   uint64 `json:"fileSize"`
	FileType   string `json:"storageType"`
	CreateTime string `json:"uploadTime"`
	FileName   string `json:"fileName"`
	FilePath   string `json:"filePath"`
	IsDir      uint8  `json:"isDir"`
	FileStatus uint8  `json:"fileStatus"`
}

var FileTypeMap = map[string]string{"1": "图片", "2": "文档", "3": "视频", "4": "音乐", "5": "其他"}

func GetFileList(fileType, path string, offset, limit int) (FileList, error) {
	userFileList := []*UserFile{}
	var count int64
	var err error
	if fileType == "0" { // 全部文件
		_ = OrmDb.Model(&UserFile{}).Where("user_id = ? AND file_path = ? AND delete_flag = 0", LoginUser.ID, path).Count(&count)
		err = OrmDb.Limit(limit).Offset(offset).Find(&userFileList, "file_path = ? AND user_id = ? AND delete_flag = 0", path, LoginUser.ID).Error
		if err != nil {
			return FileList{}, err
		}
	} else { // 根据类型查询文件
		_ = OrmDb.Model(&UserFile{}).Where("user_id = ? AND file_type = ? AND delete_flag = 0", LoginUser.ID, FileTypeMap[fileType]).Count(&count)
		err = OrmDb.Limit(limit).Offset(offset).Find(&userFileList, "file_type = ? AND user_id = ? AND delete_flag = 0", FileTypeMap[fileType], LoginUser.ID).Error
		if err != nil {
			return FileList{}, err
		}
	}

	if count != 0 {
		log.Infof("当前的count为%d", count)
		listData := []*ListData{}
		l := int(count)

		for _, v := range userFileList {
			if v.FileMd5 != "" { // 不是文件夹
				meta, err := GetFileMeta(v.FileMd5)
				if err != nil {
					log.Errorf("查找出错：%v", err)
					return FileList{}, err
				}
				FileIdMap[meta.FileMd5] = meta.ID
			}
		}

		for i := 0; i < l; i++ {
			data := &ListData{}
			data.CreateTime = userFileList[i].CreatedAt.String()[:19] // 只用显示到秒
			data.FilePath = userFileList[i].FilePath
			data.UserFileId = userFileList[i].ID

			data.Md5 = userFileList[i].FileMd5
			data.FileId = FileIdMap[data.Md5] //后期用redis
			data.UserId = userFileList[i].UserId
			data.FileSize = userFileList[i].FileSize
			data.IsDir = userFileList[i].IsDir
			data.FileType = userFileList[i].FileType
			extendName := p.Ext(userFileList[i].FileName) // .exe 带点
			if len(extendName) > 1 {
				data.ExtendName = extendName[1:]
			} else {
				data.ExtendName = ""
			}
			data.FileName = strings.TrimSuffix(userFileList[i].FileName, extendName)
			if userFileList[i].DeletedAt.Valid == false {
				data.IsDeleted = 0
			}

			listData = append(listData, data)
		}

		fileList := FileList{
			Total: count,
			List:  listData,
		}

		return fileList, nil
	}
	return FileList{List: []*ListData{}}, err
}

func CreateUserFile(file *UserFile) (err error) {
	err = OrmDb.Create(file).Error
	return
}

func GetUserFile(Md5 string) (userFile *UserFile, err error) {
	err = OrmDb.First(&userFile, "file_md5 = ? AND user_id = ? AND delete_flag = 0", Md5, LoginUser.ID).Error
	return userFile, err
}

func GetFolder(path, name string) (err error) {
	folder := &UserFile{}
	err = OrmDb.First(folder, "file_path = ? AND file_name = ? AND user_id = ? AND delete_flag = 0", path, name, LoginUser.ID).Error
	return
}

func BatchDeleteFile(files []*ListData) (err error) {
	for _, v := range files {
		if v.ExtendName != "" { // 有后缀名
			v.FileName = v.FileName + "." + v.ExtendName
		}
		err = OrmDb.Model(&UserFile{}).Where("file_name = ? AND user_id = ? AND file_path = ?", v.FileName, v.UserId, v.FilePath).Update("delete_flag", 1).Error
		if err != nil {
			return
		}
	}
	return
}

type DeleteFiles struct {
	Files []*ListData `json:"files"`
}

func GetRecoveryFileList() (list []*ListData, err error) {
	userFileList := []*UserFile{}
	err = OrmDb.Find(&userFileList, "user_id = ? AND delete_flag = 1", LoginUser.ID).Error
	if err != nil {
		return []*ListData{}, err
	}

	count := len(userFileList)
	if count != 0 {
		log.Infof("当前的count为%d", count)
		listData := []*ListData{}
		l := int(count)
		for i := 0; i < l; i++ {
			data := &ListData{}
			data.CreateTime = userFileList[i].CreatedAt.String()[:19] // 只用显示到秒
			data.FilePath = userFileList[i].FilePath
			data.UserFileId = userFileList[i].ID
			data.Md5 = userFileList[i].FileMd5
			data.UserId = userFileList[i].UserId
			data.FileSize = userFileList[i].FileSize
			data.IsDir = userFileList[i].IsDir
			data.FileType = userFileList[i].FileType
			extendName := p.Ext(userFileList[i].FileName) // .exe 带点
			if len(extendName) > 1 {
				data.ExtendName = extendName[1:]
			} else {
				data.ExtendName = ""
			}
			data.FileName = strings.TrimSuffix(userFileList[i].FileName, extendName)
			if userFileList[i].DeletedAt.Valid == false {
				data.IsDeleted = 0
			}

			listData = append(listData, data)
		}

		return listData, nil
	}
	return []*ListData{}, err
}

func GetUserFileById(id int) (file *UserFile, err error) {
	err = OrmDb.First(file, "id = ? AND delete_flag = 0", id).Error
	if err != nil {
		return nil, err
	}
	return file, nil
}
