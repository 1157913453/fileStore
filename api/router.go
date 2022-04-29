package api

import (
	dbApi "filestore/api/v1/db"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var Router *gin.Engine

func InitRouter() {
	Router = gin.Default()
	Api := Router.Group("/api")
	{
		// 文件接口
		Api.POST("/file/upload", dbApi.PostUpload)    // 上传文件数据
		Api.GET("/file/upload", dbApi.Upload)         // 上传文件
		Api.GET("/file/meta", dbApi.GetFileMeta)      // 获取文件
		Api.GET("/file/download", dbApi.DownLoadFile) // 下载文件
		Api.GET("/file/list", dbApi.GetFileList)      // 获取文件列表
		//Api.GET("/file/selectFileByType", dbApi.GetFileList)              // 获取对应类型文件列表
		Api.GET("/file/path/tree", dbApi.GetFilePathTree)         // 获取文件树
		Api.POST("/file/mkdir", dbApi.Mkdir)                      // 创建文件夹
		Api.POST("/file/batchDeleteFile", dbApi.BatchDeleteFile)  // 批量删除
		Api.POST("/recoveryFile/list", dbApi.GetRecoveryFileList) // 获取回收站文件列表

		// 用户接口
		Api.POST("/user/register", dbApi.Register)                    // 注册
		Api.GET("/user/login", dbApi.Login)                           // 登录
		Api.GET("/user/checkUserLoginInfo", dbApi.CheckUserLoginInfo) // 检查用户信息
		//Api.GET("/user/storage/info", dbApi.GetUserStorageInfo)               // 获取用户存储空间信息
		Api.GET("/user/info", dbApi.GetUserInfo) // 获取用户信息

		// 其他
		Api.GET("/filetransfer/getstorage", dbApi.GetStorage)       // 获取用户存储空间信息
		Api.GET("/filetransfer/preview", dbApi.GetImagePreview)     // 预览图片
		Api.POST("/office/previewofficefile", dbApi.GetFilePreview) // 预览

		Api.GET("/user/home", dbApi.Login)
	}
	ginpprof.Wrap(Router)
	err := Router.Run(":8888")
	if err != nil {
		log.Errorf("运行错误：%v", err)
	}

}
