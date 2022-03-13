package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"hash"
	"io"
	"os"
	"strconv"
)

type Sha1Stream struct {
	_sha1 hash.Hash
}

func FileSha1(file *os.File) string {
	file.Seek(0, 0)
	_sha1 := sha1.New()
	io.Copy(_sha1, file)
	return hex.EncodeToString(_sha1.Sum(nil))
}

func Sha1(data []byte) string {
	sha1 := sha1.New()
	sha1.Write(data)
	return hex.EncodeToString(sha1.Sum([]byte("")))
}

func MD5(data []byte) string {
	md5 := md5.New()
	md5.Write(data)
	return hex.EncodeToString(md5.Sum([]byte("")))
}

func FileMD5(file *os.File) string {
	md5 := md5.New()
	io.Copy(md5, file)
	return hex.EncodeToString(md5.Sum(nil))
}

func PathMd5(path string) string {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		fmt.Printf("打开文件%s错误:%s\n", path, err)
	}
	md5 := FileMD5(file)
	return md5
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func mergeFile(rfileName string, wfile *os.File) (err error) {
	rfile, err := os.OpenFile(rfileName, os.O_RDWR, 0666)
	defer rfile.Close()
	if err != nil {
		fmt.Println("合并时打开临时文件错误", err)
		return err
	}

	stat, err := rfile.Stat()
	if err != nil {
		fmt.Println("获取stat错误:", err)
		return err
	}
	num := stat.Size()
	buf := make([]byte, 1024*1024)
	for i := 0; int64(i) < num; {
		length, err := rfile.Read(buf)
		if err != nil {
			fmt.Println("读取文件错误：", err)
			return err
		}
		i += length
		wfile.Write(buf[:length])
	}
	return
}

// MainMergeFile 将xxx_1,xxx_2 等合并为xxx，后缀从_1开始
func MainMergeFile(connumber int, filename, targetPath string) error {
	targetFile, err := os.Create(targetPath)
	if err != nil {
		fmt.Println("创建有效文件错误：", err)
		return err
	}
	defer targetFile.Close()

	//依次对文件进行合并
	for i := 1; i <= connumber; i++ {
		mergeFile(filename+"_"+strconv.Itoa(i), targetFile)
	}

	// 删除文件
	for i := 1; i <= connumber; i++ {
		err = os.Remove(filename + "_" + strconv.Itoa(i))
		if err != nil {
			fmt.Printf("删除文件%s失败:%v", filename+"_"+strconv.Itoa(i), err)
		}
	}
	return err
}

func HashAndSalt(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePassword(hashPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Errorf("密码比较错误：", err)
		return false
	}
	return true
}
