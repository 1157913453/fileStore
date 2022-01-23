package oss_service

import (
	cfg "filestore/config"
	"filestore/service/token_service"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
)

var ossCli *oss.Client

func init() {
	var err error
	ossCli, err = oss.New(cfg.EndPoint, cfg.AccessKeyId, cfg.AccessKeySecret)
	if err != nil {
		log.Errorf("Oss Client创建失败：%v", err)
	}
}

func OssUploadPart(fileAddr string, chunkNum int) error {
	bucket, err := ossCli.Bucket(cfg.BucketName)
	if err != nil {
		log.Errorf("获取bucket:%s失败：%v", cfg.BucketName, err)
		return err
	}

	chunks, err := oss.SplitFileByPartNum(fileAddr, chunkNum)
	if err != nil {
		log.Errorf("文件：%s分片失败：%v", fileAddr, err)
		return err
	}
	fd, err := os.Open(fileAddr)
	if err != nil {
		log.Errorf("打开文件%s失败：%v", fileAddr, err)
		return err
	}
	defer fd.Close()

	options := []oss.Option{
		oss.MetadataDirective(oss.MetaReplace),
		// 指定该Object被下载时的网页缓存行为。
		// oss.CacheControl("no-cache"),
		// 指定该Object被下载时的名称。
		oss.ContentDisposition("attachment;filename=" + path.Base(fileAddr)),
		// 指定该Object的内容编码格式。
		// oss.ContentEncoding("gzip"),
		// 指定对返回的Key进行编码，目前支持URL编码。
		// oss.EncodingType("url"),
		// 指定Object的存储类型。
		oss.ObjectStorageClass(oss.StorageStandard),
	}
	log.Infof("fileAddr[15:]是%s", fileAddr[15:])
	imur, err := bucket.InitiateMultipartUpload(fileAddr[15:], options...)
	if err != nil {
		log.Errorf("初始化分片上传失败：%v", err)
		return err
	}
	var parts []oss.UploadPart
	for _, chunk := range chunks {
		fd.Seek(chunk.Offset, os.SEEK_SET)
		// 调用UploadPart方法上传每个分片。
		part, err := bucket.UploadPart(imur, fd, chunk.Size, chunk.Number)
		if err != nil {
			log.Errorf("上传oss分片失败：%v", err)
			os.Exit(-1)
			return err
		}
		parts = append(parts, part)
	}

	// 指定Object的读写权限为公共读，默认为继承Bucket的读写权限。
	//objectAcl := oss.ObjectACL(oss.ACLPrivate)

	// 步骤3：完成分片上传，指定文件读写权限为公共读。
	cmur, err := bucket.CompleteMultipartUpload(imur, parts)
	if err != nil {
		log.Errorf("合并oss分片失败：%v", err)
		os.Exit(-1)
		return err
	}
	fmt.Println("cmur:", cmur)
	return nil
}

func OssDownLoadFile(myClaims *token_service.MyClaims, fileName string) ([]byte, error) {
	bucket, err := ossCli.Bucket(cfg.BucketName)
	if err != nil {
		log.Errorf("获取bucket:%s失败：%v", cfg.BucketName, err)
		return nil, err
	}
	// 下载文件到流。
	body, err := bucket.GetObject(myClaims.Phone + "/" + fileName)
	if err != nil {
		log.Errorf("获取body流失败：%v", err)
		return nil, err
	}
	// 数据读取完成后，获取的流必须关闭，否则会造成连接泄漏，导致请求无连接可用，程序无法正常工作。
	defer body.Close()

	data, err := ioutil.ReadAll(body)
	if err != nil {
		log.Errorf("读取body流失败:%v", err)
		return nil, err
	}

	return data, nil
}
