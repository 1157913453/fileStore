package api

import (
	dbApi "filestore/src/api/v1/db"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var Router *gin.Engine

func InitRouter() {
	Router = gin.Default()
	Api := Router.Group("/api")
	{
		// 文件接口
		Api.POST("/file/upload", dbApi.PostUpload) // 真实上传接口
		Api.GET("/file/upload", dbApi.Upload)
		Api.GET("/file/meta", dbApi.GetFileMeta)         // 获取文件
		Api.GET("/file/download", dbApi.DownLoadFile)    // 下载文件
		Api.PUT("/file/update", dbApi.UpdateFileMeta)    // 修改文件
		Api.DELETE("/file/delete", dbApi.DeleteFileMeta) // 软删除文件
		Api.GET("/file/list", dbApi.GetFileList)         // 获取文件列表
		//Api.GET("/file/selectFileByType", dbApi.GetFileList) // 获取对应类型文件列表
		Api.GET("/file/path/tree", dbApi.GetFilePathTree) //
		Api.POST("/file/mkdir", dbApi.Mkdir)
		Api.POST("/file/batchDeleteFile", dbApi.BatchDeleteFile)
		Api.POST("/recoveryFile/list", dbApi.GetRecoveryFileList)

		// 用户接口
		Api.POST("/user/register", dbApi.Register)
		Api.GET("/user/login", dbApi.Login)
		Api.GET("/user/checkuserlogininfo", dbApi.CheckUserLoginInfo)
		Api.POST("/user/login", dbApi.Login)
		Api.GET("/user/storage/info", dbApi.GetUserStorageInfo)
		Api.GET("/user/info", dbApi.GetUserInfo)

		// 其他
		Api.GET("/filetransfer/getstorage", dbApi.GetStorage)

		Api.GET("/user/home", dbApi.Login)
	}

	err := Router.Run(":8888")
	if err != nil {
		log.Errorf("运行错误：%v", err)
	}
}
