package file_service

import (
	"bytes"
	"encoding/hex"
	"errors"
	"filestore/models"
	"filestore/service/token_service"
	"filestore/service/user_service"
	"filestore/util"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
)

var (
	fileTypeMap sync.Map
	Image       = []string{"png", "jpg", "jpeg", "gif", "webp", "bmp", "pcx", "tif", "tga", "exif", "fpx", "svg", "ico", "psd", "cdr", "pcd", "dxf", "ufo", "eps", "ai", "hdri", "raw", "wfm", "flic", "emf", "avif", "apng"}
	Video       = []string{"wmv", "asf", "asxc", "rmvb", "rm", "mpg", "mpeg", "mpe", "mp4", "3gp", "mov", "avi", "m4v", "dat", "mkv", "flv", "vob", "qt", "ogg", "mod", "ram", "viv"}
	Document    = []string{"doc", "docx", "docv", "xml", "xls", "xlsx", "pdf", "txt", "ppt"}
	Music       = []string{"mp3", "ape", "wav", "flac", "ape"}
)

func init() {
	//longFileTypeMap.Store(strings.ToLower("D0CF11E0A1B11AE1000000000000000000000000000000003E000300FEFF090006000000000000000000000001000000"),"xls")
	//longFileTypeMap.Store(strings.ToLower("D0CF11E0A1B11AE1000000000000000000000000000000003E000300FEFF090006000000000000000000000003000000"),"doc")
	fileTypeMap.Store("ffd8ff", "jpg")         //JPEG (jpg)
	fileTypeMap.Store("89504e47", "png")       //PNG (png)
	fileTypeMap.Store("47494638", "gif")       //GIF (gif)
	fileTypeMap.Store("49492a00", "tif")       //TIFF (tif)
	fileTypeMap.Store("424d", "bmp")           //16色位图(bmp)
	fileTypeMap.Store("41433130", "dwg")       //CAD (dwg)
	fileTypeMap.Store("3c21444f43545", "html") //HTML (html)   3c68746d6c3e0  3c68746d6c3e0
	fileTypeMap.Store("3c68746d6c3e0", "html") //HTML (html)   3c68746d6c3e0  3c68746d6c3e0
	fileTypeMap.Store("3c21646f63747", "htm")  //HTM (htm)
	fileTypeMap.Store("48544d4c207b0", "css")  //css
	fileTypeMap.Store("696b2e71623d6", "js")   //js
	fileTypeMap.Store("7B5C727466315", "rtf")  // 我（IBAS）猜想的rtf
	fileTypeMap.Store("7b5c727466315", "rtf")  //Rich Text Format (rtf)
	fileTypeMap.Store("38425053", "psd")       //Photoshop (psd)
	fileTypeMap.Store("46726f6d3a203d", "eml") //Email [Outlook Express 6] (eml)
	fileTypeMap.Store("d0cf11e0a1b11a", "doc") //MS Excel 注意：word、msi 和 excel的文件头一样
	fileTypeMap.Store("d0cf11e0a1b11a", "vsd") //Visio 绘图
	fileTypeMap.Store("5374616E646172", "mdb") //MS Access (mdb)
	fileTypeMap.Store("252150532D4164", "ps")
	fileTypeMap.Store("255044462d312e", "pdf")  //Adobe Acrobat (pdf)
	fileTypeMap.Store("2e524d46000000", "rmvb") //rmvb/rm相同
	fileTypeMap.Store("464c5601050000", "flv")  //flv与f4v相同
	fileTypeMap.Store("00000020667479", "mp4")
	fileTypeMap.Store("49443303000000", "mp3")
	fileTypeMap.Store("000001ba210001", "mpg") //
	fileTypeMap.Store("3026b2758e66cf", "wmv") //wmv与asf相同
	fileTypeMap.Store("52494646e27807", "wav") //Wave (wav)
	fileTypeMap.Store("52494646d07d60", "avi")
	fileTypeMap.Store("4d546864000000", "mid") //MIDI (mid)
	fileTypeMap.Store("504b03040a", "zip")     // 我（IBAS）看到的zip
	fileTypeMap.Store("504b030414", "zip")
	fileTypeMap.Store("526172211a0700", "rar") // 我（IBAS）看到的rar
	fileTypeMap.Store("526172211a0700", "rar")
	fileTypeMap.Store("23546869732063", "ini")
	fileTypeMap.Store("504b0304140008", "jar") // 我（IBAS）看到的jar
	fileTypeMap.Store("504b03040a0000", "jar")
	fileTypeMap.Store("4d5a9000030000", "exe")        //可执行文件
	fileTypeMap.Store("3c254020706167", "jsp")        //jsp文件
	fileTypeMap.Store("4d616e69666573", "mf")         //MF文件
	fileTypeMap.Store("3c3f786d6c2076", "xml")        //xml文件
	fileTypeMap.Store("494e5345525420", "sql")        //xml文件
	fileTypeMap.Store("7061636b616765", "java")       //java文件
	fileTypeMap.Store("406563686f206f", "bat")        //bat文件
	fileTypeMap.Store("1f8b0800000000", "gz")         //gz文件
	fileTypeMap.Store("6c6f67346a2e72", "properties") //bat文件
	fileTypeMap.Store("cafebabe000000", "class")      //bat文件
	fileTypeMap.Store("49545346030000", "chm")        //bat文件
	fileTypeMap.Store("04000000010000", "mxp")        //bat文件
	fileTypeMap.Store("504b0304140006", "docx")       //docx文件
	fileTypeMap.Store("d0cf11e0a1b11a", "wps")        //WPS文字wps、表格et、演示dps都是一样的
	fileTypeMap.Store("6431303a637265", "torrent")
	fileTypeMap.Store("6D6F6F76", "mov")         //Quicktime (mov)
	fileTypeMap.Store("FF575043", "wpd")         //WordPerfect (wpd)
	fileTypeMap.Store("CFAD12FEC5FD746F", "dbx") //Outlook Express (dbx)
	fileTypeMap.Store("2142444E", "pst")         //Outlook (pst)
	fileTypeMap.Store("AC9EBD8F", "qdf")         //Quicken (qdf)
	fileTypeMap.Store("E3828596", "pwl")         //Windows Password (pwl)
	fileTypeMap.Store("2E7261FD", "ram")         //Real Audio (ram)
}

// 获取前面结果字节的二进制
func bytesToHexString(src []byte) string {
	res := bytes.Buffer{}
	if src == nil || len(src) <= 0 {
		return ""
	}
	temp := make([]byte, 0)
	for _, v := range src {
		sub := v & 0xFF
		hv := hex.EncodeToString(append(temp, sub))
		if len(hv) < 2 {
			res.WriteString(strconv.FormatInt(int64(0), 10))
		}
		res.WriteString(hv)
	}
	return res.String()
}

// GetFileType 用文件前面几个字节来判断
// fSrc: 文件字节流（就用前面几个字节）
func GetFileType(fSrc []byte) (fileType string) {
	fileCode := bytesToHexString(fSrc)

	fileTypeMap.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		if strings.HasPrefix(fileCode, strings.ToLower(k)) ||
			strings.HasPrefix(k, strings.ToLower(fileCode)) {
			fileType = v
			return false
		}
		return true
	})
	if fileType == "" {
		fileType = "其他"
		return
	}
	for _, v := range Image {
		if fileType == v {
			fileType = "图片"
			return
		}
	}
	for _, v := range Video {
		if fileType == v {
			fileType = "视频"
			return
		}
	}
	for _, v := range Document {
		if fileType == v {
			fileType = "文档"
			return
		}
	}
	for _, v := range Music {
		if fileType == v {
			fileType = "音乐"
			return
		}
	}
	return
}

func MergeFile(myClaims *token_service.MyClaims, fileName, ChunkPath, Md5 string, totalchunks int) (err error) {
	targetPath := "/tmp/fileStore/" + myClaims.Phone + "/" + fileName
	err = util.MainMergeFile(totalchunks, ChunkPath+"/"+fileName, targetPath)
	if err != nil {
		log.Errorf("合并文件出错：%v", err)
		return
	}

	currentMd5 := util.PathMd5(targetPath)
	if currentMd5 != Md5 {
		log.Errorf("上传前后Md5不一致")
		return errors.New("上传前后Md5不一致")
	}

	return nil
}

func UpdateDbFile(myClaims *token_service.MyClaims, fileName, filePath, Md5 string, totalsize int) error {
	targetPath := "/tmp/fileStore/" + myClaims.Phone + "/" + fileName
	f, err := ioutil.ReadFile(targetPath)
	if err != nil {
		log.Errorf("读取target文件错误：%v", err)
		return err
	}
	fileType := GetFileType(f[:10])
	log.Infof("文件类型为：%s", fileType)

	// 存储文件信息到数据库
	fileMeta := &models.File{
		FileMd5:    Md5,
		FileName:   fileName,
		FileSize:   uint64(totalsize),
		FileAddr:   targetPath,
		FileStatus: 0,
		FileType:   fileType,
		PointCount: 1,
	}

	var userInfo *models.User

	userInfo, err = user_service.GetUser(myClaims.Phone)
	if err != nil {
		log.Errorf("获取%s用户信息失败:%v", myClaims.Phone, err)
		return err
	}
	userFile := &models.UserFile{
		UserId:     userInfo.ID,
		FileMd5:    Md5,
		FileSize:   uint64(totalsize),
		FileName:   fileName,
		FilePath:   filePath,
		FileType:   fileType,
		IsDir:      0,
		FileStatus: 0,
	}

	// 更新file表
	err = CreateFileMeta(fileMeta)
	if err != nil {
		log.Errorf("保存文件出错：%v", err)
		return err
	}

	// 删除chunk中的缓存
	err = DeleteFileChunk(Md5)
	if err != nil {
		log.Errorf("清除chunk缓存失败：%v", err)
		return err
	}

	// 更新userFile表
	err = CreateUserFile(userFile)
	if err != nil {
		log.Errorf("插入用户表出错：%v", err)
		return err
	}
	return nil
}
