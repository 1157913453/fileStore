package models

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Model struct {
	Id        uint      `gorm:"primarykey"`
	CreateAt  time.Time `gorm:"autoCreateTime" json:"createAt"`
	UpdateAt  time.Time `gorm:"autoUpdateTime" json:"updateAt"`
	IsDeleted uint      `gorm:"type:tinyint;not null;default:0" json:"isDeleted"`
}

var OrmDb = &gorm.DB{}

func InitDB() {
	var err error
	dsn := "root:aaliyy46821593@tcp(118.178.237.76:3306)/fileStore?charset=utf8mb4&parseTime=True&loc=Local"
	OrmDb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Errorf("连接数据库失败：%v", err)
		return
	}

	err = OrmDb.AutoMigrate(&File{}, &User{}, &UserToken{}, &UserFile{}, &UploadTask{}, UploadTaskDetail{}, &Chunk{})
	if err != nil {
		log.Panicf("迁移表格失败：%v", err)
	}

}
