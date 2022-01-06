package file_service

import (
	"bytes"
	"encoding/hex"
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
	fileTypeMap.Store("ffd8ff", "jpg")                //JPEG (jpg)
	fileTypeMap.Store("89504e47", "png")              //PNG (png)
	fileTypeMap.Store("47494638", "gif")              //GIF (gif)
	fileTypeMap.Store("49492a00", "tif")              //TIFF (tif)
	fileTypeMap.Store("424d", "bmp")                  //16色位图(bmp)
	fileTypeMap.Store("41433130", "dwg")              //CAD (dwg)
	fileTypeMap.Store("3c21444f435459504520", "html") //HTML (html)   3c68746d6c3e0  3c68746d6c3e0
	fileTypeMap.Store("3c68746d6c3e0", "html")        //HTML (html)   3c68746d6c3e0  3c68746d6c3e0
	fileTypeMap.Store("3c21646f637479706520", "htm")  //HTM (htm)
	fileTypeMap.Store("48544d4c207b0d0a0942", "css")  //css
	fileTypeMap.Store("696b2e71623d696b2e71", "js")   //js
	fileTypeMap.Store("7B5C727466315c", "rtf")        // 我（IBAS）猜想的rtf
	fileTypeMap.Store("7b5c727466315c616e73", "rtf")  //Rich Text Format (rtf)
	fileTypeMap.Store("38425053", "psd")              //Photoshop (psd)
	fileTypeMap.Store("46726f6d3a203d3f6762", "eml")  //Email [Outlook Express 6] (eml)
	fileTypeMap.Store("d0cf11e0a1b11ae10000", "doc")  //MS Excel 注意：word、msi 和 excel的文件头一样
	fileTypeMap.Store("d0cf11e0a1b11ae10000", "vsd")  //Visio 绘图
	fileTypeMap.Store("5374616E64617264204A", "mdb")  //MS Access (mdb)
	fileTypeMap.Store("252150532D41646F6265", "ps")
	fileTypeMap.Store("255044462d312e350d0a", "pdf")  //Adobe Acrobat (pdf)
	fileTypeMap.Store("2e524d46000000120001", "rmvb") //rmvb/rm相同
	fileTypeMap.Store("464c5601050000000900", "flv")  //flv与f4v相同
	fileTypeMap.Store("00000020667479706d70", "mp4")
	fileTypeMap.Store("49443303000000002176", "mp3")
	fileTypeMap.Store("000001ba210001000180", "mpg") //
	fileTypeMap.Store("3026b2758e66cf11a6d9", "wmv") //wmv与asf相同
	fileTypeMap.Store("52494646e27807005741", "wav") //Wave (wav)
	fileTypeMap.Store("52494646d07d60074156", "avi")
	fileTypeMap.Store("4d546864000000060001", "mid") //MIDI (mid)
	fileTypeMap.Store("504b03040a", "zip")           // 我（IBAS）看到的zip
	fileTypeMap.Store("504b030414", "zip")
	fileTypeMap.Store("526172211a07009bd773", "rar") // 我（IBAS）看到的rar
	fileTypeMap.Store("526172211a0700cf9073", "rar")
	fileTypeMap.Store("235468697320636f6e66", "ini")
	fileTypeMap.Store("504b0304140008080800", "jar") // 我（IBAS）看到的jar
	fileTypeMap.Store("504b03040a0000000000", "jar")
	fileTypeMap.Store("4d5a9000030000000400", "exe")        //可执行文件
	fileTypeMap.Store("3c25402070616765206c", "jsp")        //jsp文件
	fileTypeMap.Store("4d616e69666573742d56", "mf")         //MF文件
	fileTypeMap.Store("3c3f786d6c2076657273", "xml")        //xml文件
	fileTypeMap.Store("494e5345525420494e54", "sql")        //xml文件
	fileTypeMap.Store("7061636b616765207765", "java")       //java文件
	fileTypeMap.Store("406563686f206f66660d", "bat")        //bat文件
	fileTypeMap.Store("1f8b0800000000000000", "gz")         //gz文件
	fileTypeMap.Store("6c6f67346a2e726f6f74", "properties") //bat文件
	fileTypeMap.Store("cafebabe0000002e0041", "class")      //bat文件
	fileTypeMap.Store("49545346030000006000", "chm")        //bat文件
	fileTypeMap.Store("04000000010000001300", "mxp")        //bat文件
	fileTypeMap.Store("504b0304140006000800", "docx")       //docx文件
	fileTypeMap.Store("d0cf11e0a1b11ae10000", "wps")        //WPS文字wps、表格et、演示dps都是一样的
	fileTypeMap.Store("6431303a637265617465", "torrent")
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

// 用文件前面几个字节来判断
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
