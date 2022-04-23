package db

import (
	"encoding/json"
	"filestore/config"
	"filestore/models"
	"filestore/payload"
	"filestore/service/cache_service"
	"filestore/service/file_service"
	"filestore/service/oss_service"
	"filestore/service/rabbitmq_service"
	"filestore/service/token_service"
	"filestore/service/user_service"
	"filestore/util"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"strconv"
)

func Upload(c *gin.Context) {
	// 检查token
	myClaims, err := token_service.CheckToken(c)
	if err != nil {
		c.JSON(200, payload.FailPayload("token无效"))
	}
	// 根据Md5判断是否已有文件
	chunkNumber, chunkSize, currentChunkSize, filePath := c.Query("chunkNumber"), c.Query("chunkSize"), c.Query("currentChunkSize"), c.Query("filePath")
	fileName, Md5, relativePath, totalChunks, totalSize := c.Query("filename"), c.Query("identifier"), c.Query("relativePath"), c.Query("totalChunks"), c.Query("totalSize")
	fmt.Println(relativePath, totalChunks, chunkNumber, chunkSize, currentChunkSize)
	var uploads []int
	file, err := file_service.GetFileMeta(Md5)
	if err != nil {
		log.Infof("数据库中没有该文件,需要上传,err为：%v", err)
		details, err := file_service.QueryUploadTaskDetails(Md5) // 查询chunk信息
		if err != nil {                                          // 没有chunk信息
			log.Errorf("查找task错误是：%v", err)
			if err.Error() == "record not found" {
				details = []*models.UploadTaskDetail{}
			} else {
				log.Errorf("查询任务切片失败：%s", err)
				c.JSON(200, payload.FailPayload("查询任务失败"))
				return
			}
		} else { //有chunk信息
			for _, v := range details {
				uploads = append(uploads, int(v.ChunkNumber))
			}
		}

		c.JSON(200, payload.NormalUpload(uploads, true))
		fmt.Println("uploads是：", uploads)
		return
		// 不分片
	}
	log.Infof("数据库中已有该文件，开始极速上传")

	// 获取用户信息
	userInfo, err := cache_service.GetUserCache(myClaims.Phone)
	if err != nil {
		c.JSON(200, payload.FailPayload("获取缓存错误"))
		return
	}

	// 判断已有文件是否在当前path
	userFile, err := file_service.GetUserFile(userInfo.ID, Md5)
	if err == nil && userFile != nil && userFile.FilePath == filePath {
		c.JSON(200, payload.ExistsUpload())
		return
	}
	filesize, _ := strconv.Atoi(totalSize)

	userFile = &models.UserFile{
		UserId:   userInfo.ID,
		FileMd5:  Md5,
		FileSize: uint64(filesize),
		FileName: fileName,
		FilePath: filePath,
		FileType: file.FileType,
	}

	// 插入用户文件表
	if totalChunks == "1" { // 说明文件太小
		err = file_service.CreateUserFile(userFile)
		if err != nil {
			log.Errorf("插入用户表出错：%v", err)
			c.JSON(200, payload.FailPayload("插入用户表出错"))
			return
		}
	}

	c.JSON(200, payload.FastUpload(nil, true))
}

func PostUpload(c *gin.Context) {
	// 检查token
	myClaims, err := token_service.CheckToken(c)
	if err != nil {
		c.JSON(200, payload.FailPayload("token无效"))
	}

	//接受file数据流
	file, err := c.FormFile("file")
	if err != nil {
		log.Errorf("接受文件流错误：%v", err)
		c.JSON(200, payload.FailPayload("接受文件流错误"))
		return
	}

	// 存储chunk数据
	chunkNumber, chunkSize, _, filePath := c.PostForm("chunkNumber"), c.PostForm("chunkSize"), c.PostForm("currentChunkSize"), c.PostForm("filePath")
	fileName, Md5, relativePath, totalChunks, totalSize := c.PostForm("filename"), c.PostForm("identifier"), c.PostForm("relativePath"), c.PostForm("totalChunks"), c.PostForm("totalSize")
	// 给切片文件重命名
	file.Filename = fileName + "_" + chunkNumber
	chunkNum, _ := strconv.Atoi(chunkNumber)
	Size, _ := strconv.Atoi(chunkSize)
	totalsize, _ := strconv.Atoi(totalSize)
	chunks, _ := strconv.Atoi(totalChunks)
	totalchunks, _ := strconv.Atoi(totalChunks)

	err = c.SaveUploadedFile(file, config.ChunkPath+"/"+file.Filename) // 保存文件
	if err != nil {
		log.Errorf("保存文件%s失败：%v", file.Filename, err)
		c.JSON(200, payload.FailPayload("保存文件失败："+err.Error()))
		return
	}
	// 将chunk数据存储到数据库
	uploadTaskDetail := &models.UploadTaskDetail{
		FileMd5:      Md5,
		ChunkNumber:  int64(chunkNum),
		ChunkSize:    int64(Size),
		TotalChunks:  int32(chunks),
		TotalSize:    int64(totalsize),
		FilePath:     filePath,
		RelativePath: relativePath,
	}
	err = file_service.CreateUploadTaskDetail(uploadTaskDetail)
	if err != nil {
		log.Errorf("数据库增加%s失败：%v", filePath, err)
		return
	}

	// 如果是最后一个分片文件
	if chunkNum == totalchunks {
		// 合并所有文件
		err = file_service.MergeFile(myClaims, fileName, config.ChunkPath, Md5, totalchunks)
		if err != nil {
			log.Errorf("合并文件出错：%v", err)
			c.JSON(200, payload.FailPayload("合并文件出错"))
			return
		}

		fileAddr := config.BasePath + myClaims.Phone + "/" + fileName
		// 发送到RabbitMQ队列中
		err = rabbitmq_service.SendMQ(fileAddr, chunkNum)
		if err != nil {
			c.JSON(200, "RabbitMQ消息队列异常")
			return
		}

		//go func() {
		//	err := oss_service.OssUploadPart(fileAddr, chunkNum)
		//	if err != nil {
		//		log.Errorf("上传OSS错误%v", err)
		//		return
		//	}
		//
		//}()

		// 更新数据库
		err = file_service.UpdateDbFile(myClaims, fileName, filePath, Md5, totalsize)
		if err != nil {
			log.Errorf("更新数据库文件失败:%v", err)
			c.JSON(200, payload.FailPayload("更新数据库文件失败"))
			return
		}
	}
	c.JSON(200, payload.UploadRes(true))
}

func GetFileMeta(c *gin.Context) {
	fileMd5 := c.Query("fileMd5")
	fMeta, err := file_service.GetFileMeta(fileMd5)
	if err != nil {
		c.JSON(200, payload.FailPayload("获取FileMeta失败："+err.Error()))
		return
	}
	data, err := json.Marshal(fMeta)
	if err != nil {
		c.JSON(200, payload.FailPayload("序列化失败："+err.Error()))
		return
	}

	c.JSON(200, payload.SucDataPayload("获取成功", string(data)))
}

func GetFileList(c *gin.Context) {
	myClaims, err := token_service.CheckToken(c)
	if err != nil {
		c.JSON(200, payload.FailPayload("token无效"))
		return
	}
	filePath, page, pageCount, fileType := c.DefaultQuery("filePath", "/"),
		c.DefaultQuery("currentPage", "1"), //fileType: 0为全部文件，1、2、3、4、5分别对应图片，视频，文档，音乐，其他
		c.DefaultQuery("pageCount", "50"),
		c.DefaultQuery("fileType", "0")
	Page, _ := strconv.Atoi(page)
	PageCount, _ := strconv.Atoi(pageCount)
	listData, err := file_service.GetFileList(myClaims.Phone, fileType, filePath, Page, PageCount)
	if err != nil {
		log.Errorf("获取文件列表出错:%v", err)
		c.JSON(200, payload.FailPayload("获取文件列表出错"))
		return
	}

	c.JSON(200, payload.SucFileListPayload("获取文件列表成功", true, listData))

	//data := []byte(`{
	//"code":0,
	//"data":{
	//	"total": 1,
	//"list":[{
	//	"fileId":1,
	//	"deleteFlag":0,
	//	"extendName":"gg",
	//	"fileName":"444",
	//	"filePath":"/",
	//	"fileSize":4554,
	//	"fileUrl":"upload/20211223/d77ba387-fdfa-48bc-885b-0a4599e4ef37.gg",
	//	"identifier":"d77ba387-fdfa-48bc-885b",
	//	"isDir" :0,
	//	"storageType": 1,
	//	"uploadTime":"2021-12-23 01:26:23",
	//	"userId":789,
	//	"userFileId":1234
	//}]
	//},
	//"message": "成功",
	//"success": true
	//}`)
	//js, err := simplejson.NewJson(data)
	//if err != nil {
	//	log.Errorf("e是：%v", err)
	//}
	//
	////d1 := &Dd{}
	////err := json.Unmarshal(data, d1)
	////log.Errorf("err是：%v", err)
	////log.Infof("错误是：%v", d1)
	//c.JSON(200, js)
}

func DownLoadFile(c *gin.Context) {
	token := c.Query("token")
	myClaims, err := token_service.ParseToken(token)
	if err != nil {
		log.Errorf("token无效:%v", err)
		c.JSON(200, payload.FailPayload("token无效"))
		return
	}

	fileId := c.Query("fileId")
	log.Infof("下载的文件id是:%s", fileId)
	downloadFileId, _ := strconv.Atoi(fileId)

	file, err := file_service.GetFileById(downloadFileId)
	if err != nil {
		log.Errorf("查找文件id错误：%v", err)
		c.JSON(200, payload.FailPayload("查找文件id错误"))
		return
	}

	// 如果OSS未上传完就马上下载，则从本地返回文件
	exit, err := util.PathExists(file.FileAddr)
	if err != nil {
		log.Errorf("查询文件是否存在失败：%v", err)
		return
	}
	if exit {
		// 从本地返回文件
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.FileName))
		c.File(file.FileAddr)
		return
	}

	// 从OSS返回文件流
	data, err := oss_service.OssDownLoadFile(myClaims, file.FileName)
	if err != nil {
		c.JSON(200, payload.FailPayload("从OSS下载文件失败"))
		return
	}
	c.Data(200, "application/octet-stream", data)
	log.Infof("文件从OSS下载成功")
}

func Mkdir(c *gin.Context) {
	myClaims, err := token_service.CheckToken(c)
	if err != nil {
		c.JSON(200, payload.FailPayload("token无效"))
		return
	}
	folder := new(models.Folder)
	err = c.ShouldBindJSON(folder)
	if err != nil {
		log.Errorf("绑定folder参数错误：%v", err)
		c.JSON(200, payload.FailPayload("绑定folder参数错误"))
		return
	}

	// 获取用户信息
	var userInfo *models.User

	userInfo, err = user_service.GetUser(myClaims.Phone)
	if err != nil {
		c.JSON(200, payload.FailPayload("查询用户失败"))
		return
	}

	// 判断当前路径是否有同名文件夹
	err = file_service.GetFolder(userInfo.ID, folder.FilePath, folder.FoldName)
	if err == nil {
		log.Errorf("有同名文件夹：%v", err)
		c.JSON(200, payload.FailPayload("有同名文件夹"))
		return
	} else if err.Error() != "record not found" {
		log.Errorf("查询文件夹出错：%v", err)
		c.JSON(200, payload.FailPayload("查询文件夹出错"))
		return
	}

	//foldName, filePath := c.PostForm("foldName"), c.PostForm("filePath")
	//log.Infof("foldName is :%v and path is :%v", foldName, filePath)

	file := &models.UserFile{
		UserId:   userInfo.ID,
		FileName: folder.FoldName,
		FilePath: folder.FilePath,
		IsDir:    1,
	}
	err = file_service.CreateUserFile(file)
	if err != nil {
		log.Errorf("创建文件夹失败：%v", err)
		c.JSON(200, payload.FailPayload("创建文化夹失败"))
		return
	}
	c.JSON(200, payload.SucDataPayload("创建文件夹成功", nil))
}

func BatchDeleteFile(c *gin.Context) {
	lists := models.DeleteFiles{}
	err := c.ShouldBindJSON(&lists)

	if err != nil {
		log.Errorf("批量删除文件时绑定参数出错：%s\n", err.Error())
		c.JSON(200, payload.FailPayload("批量删除文件时绑定参数出错"))
		return
	}

	err = file_service.BatchDeleteFile(lists)
	if err != nil {
		log.Errorf("批量删除文件时出错：%v", err)
		c.JSON(200, payload.FailPayload("批量删除文件时出错"))
		return
	}
	c.JSON(200, payload.SucPayload("批量删除文件成功"))
}

func GetRecoveryFileList(c *gin.Context) {
	myClaims, err := token_service.CheckToken(c)
	if err != nil {
		c.JSON(200, payload.FailPayload("token无效"))
		return
	}

	userInfo, err := user_service.GetUser(myClaims.Phone)
	if err != nil {
		c.JSON(200, payload.FailPayload("获取用户信息失败"))
		return
	}

	list, err := file_service.GetRecoveryFileList(userInfo.ID)
	if err != nil {
		log.Errorf("查询回收站数据失败：%v", err)
		c.JSON(200, payload.FailPayload("查询回收站数据失败"))
		return
	}
	c.JSON(200, payload.SucDataPayload("查询回收站数据成功", list))
}

func GetFilePreview(c *gin.Context) {
	token := c.Query("token")
	myClaims, err := token_service.ParseToken(token)
	if err != nil {
		log.Errorf("token无效:%v", err)
		c.JSON(200, payload.FailPayload("token无效"))
		return
	}

	fileId := c.Query("fileId")
	previewFileId, _ := strconv.Atoi(fileId)

	file, err := file_service.GetFileById(previewFileId)
	if err != nil {
		log.Errorf("查找文件id错误：%v", err)
		c.JSON(200, payload.FailPayload("查找文件id错误"))
		return
	}

	// 如果OSS未上传完就马上下载，则从本地返回文件
	exit, err := util.PathExists(file.FileAddr)
	if err != nil {
		log.Errorf("查询文件是否存在失败：%v", err)
		return
	}
	if exit {
		// 从本地返回文件
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Content-Disposition", fmt.Sprintf("filename=%s", file.FileName))
		fileData, err := ioutil.ReadFile(file.FileAddr)
		if err != nil {
			log.Errorf("打开文件错误：%v", err)
			c.JSON(200, payload.FailPayload("打开文件失败"))
			return
		}
		c.Data(200, "", fileData)
		log.Infof("从本地返回文件数据成功")
		return
	}

	// 从OSS返回文件流
	data, err := oss_service.OssDownLoadFile(myClaims, file.FileName)
	if err != nil {
		c.JSON(200, payload.FailPayload("从OSS下载文件失败"))
		return
	}
	c.Data(200, "", data)
	log.Infof("文件从OSS预览成功")
}

//
//func GetFileListByType(c *gin.Context) {
//	filetype := c.DefaultQuery("fileType", "0")
//	page := c.DefaultQuery("currentPage", "1")
//	pagecount := c.DefaultQuery("pageCount", "50")
//	Page, _ := strconv.Atoi(page)
//	pageCount, _ := strconv.Atoi(pagecount)
//	fileList, err := file_service.GetFileListByType()
//}
