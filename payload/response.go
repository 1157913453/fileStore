package payload

import (
	"filestore/models"
	"time"
)

type PayloadUpload struct {
	Code    int32        `json:"code"`
	Data    *PayloadData `json:"data"`
	Message string       `json:"message"`
	Success bool         `json:"success"`
}

type PayloadData struct {
	NeedMerge     bool      `json:"needMerge"`
	SkipUpload    bool      `json:"skipUpload"`
	TimeStampName time.Time `json:"timeStampName"`
	Uploaded      []int     `json:"uploaded"`
}

type PayloadFileList struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Success bool            `json:"success"`
	Data    models.FileList `json:"data"`
}

func NormalUpload(uploads []int, isSuccess bool) PayloadUpload {
	payloadData := &PayloadData{
		NeedMerge:     false,
		SkipUpload:    false,
		TimeStampName: time.Now(),
		Uploaded:      uploads,
	}
	data := PayloadUpload{
		Code:    0,
		Data:    payloadData,
		Message: "普通上传",
		Success: isSuccess,
	}
	return data
}

func FastUpload(uploads []int, isSuccess bool) PayloadUpload {
	payloadData := &PayloadData{
		NeedMerge:     false,
		SkipUpload:    true,
		TimeStampName: time.Now(),
		Uploaded:      uploads,
	}
	data := PayloadUpload{
		Code:    0,
		Data:    payloadData,
		Message: "极速上传",
		Success: isSuccess,
	}
	return data
}

func ExistsUpload() PayloadUpload {
	payloadData := &PayloadData{
		NeedMerge:     false,
		SkipUpload:    true,
		TimeStampName: time.Now(),
	}
	data := PayloadUpload{
		Code:    0,
		Data:    payloadData,
		Message: "当前目录已有文件",
		Success: true,
	}
	return data
}

func UploadRes(isSuccess bool) PayloadUpload {
	payloadData := &PayloadData{
		NeedMerge:     false,
		SkipUpload:    false,
		TimeStampName: time.Now(),
	}
	data := PayloadUpload{
		Code:    0,
		Data:    payloadData,
		Message: "普通上传成功",
		Success: isSuccess,
	}
	return data
}

type Payload struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SucPayload(sucMsg string) Payload {
	return Payload{
		Success: true,
		Code:    0,
		Message: sucMsg,
	}
}

func SucDataPayload(sucMsg string, data interface{}) Payload {
	return Payload{
		Success: true,
		Code:    0,
		Message: sucMsg,
		Data:    data,
	}
}

func FailPayload(errMsg string) Payload {
	return Payload{
		Success: false,
		Message: errMsg,
	}
}

func SucFileListPayload(msg string, isSuccess bool, data models.FileList) (fileList PayloadFileList) {
	fileList = PayloadFileList{
		Code:    0,
		Message: msg,
		Success: isSuccess,
		Data:    data,
	}
	return
}

func FailFileListPayload(path, msg string, isSuccess bool, data models.FileList) (fileList PayloadFileList) {

	fileList = PayloadFileList{
		Message: msg,
		Success: isSuccess,
		Data:    data,
	}
	return
}
