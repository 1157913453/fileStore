package db

import (
	"filestore/src/models"
	"filestore/src/payload"
	"filestore/src/service/cache_service"
	"filestore/src/service/file_service"
	"filestore/src/service/token_service"
	"filestore/src/service/user_service"
	"filestore/src/util"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

const (
	pwdSalt = "*923!kj"
)

type UserInfo struct {
	*models.User
	Token string `json:"token"`
	//Id            int    `json:"Id"`
}

type ResUserInfo struct {
	Code    int8      `json:"code"`
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Data    *UserInfo `json:"data"`
}

type RegisterModel struct {
	UserName string `json:"userName"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {
	var register RegisterModel
	err := c.ShouldBindJSON(&register)
	if err != nil {
		log.Errorf("绑定错误：%v", err)
		c.JSON(200, payload.FailPayload("参数绑定错误"))
		return
	}
	phone := register.Phone
	userName := register.UserName
	password := register.Password
	if len(phone) < 11 || len(password) < 5 {
		log.Errorf("手机号或密码太短")
		c.JSON(200, payload.FailPayload("手机号或密码太短"))
		return
	}

	// 对密码加密
	encPassword := util.Sha1([]byte(password + pwdSalt))
	err = user_service.CreateUser(phone, encPassword, userName)
	if err != nil {
		log.Errorf("注册失败：%v", err)
		c.JSON(200, payload.FailPayload("注册失败,已有该用户名"))
		return
	}

	path := "/tmp/fileStore/" + phone
	exists, err := util.PathExists(path)
	if err != nil {
		log.Errorf("判断用户目录是否存在失败:%v", err)
		c.JSON(200, payload.FailPayload("判断用户是否存在失败"))
		return
	}
	if exists {
		log.Infof("%s目录已存在", phone)
	} else {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Errorf("创建%s用户文件夹失败:%v", phone, err)
			c.JSON(200, payload.FailPayload("创建用户文件夹失败"))
			return
		}
		log.Infof("创建用户%s文件夹成功", phone)
	}
	c.JSON(200, payload.SucPayload("注册成功"))
}

func CheckUserLoginInfo(c *gin.Context) {
	token := c.GetHeader("token")

	// 判断token是否有效
	myClaims, err := token_service.ParseToken(token)
	//_, err := token_service.ParseToken(token)
	if err != nil {
		log.Errorf("token 无效：%v", err)
		c.JSON(200, payload.FailPayload("token无效"))
		return
	}

	// 判断缓存是否存在
	_, err = cache_service.GetUserCache(myClaims.Phone)
	if err != nil {
		// 设置缓存
		userInfo, err := user_service.GetUserByPhone(myClaims.Phone)
		11
		if err != nil {
			return
		}
		err = cache_service.AddUserCache(userInfo)
		if err != nil {
			log.Errorf("设置用户%s缓存错误:%v", userInfo.UserName, err)
			return
		}
	}

	//// 获取用户信息
	//	models.LoginUser, err = user_service.GetUserByPhone(myClaims.Phone)
	////user, err := user_service.GetUserInfoByToken(token)
	//if err != nil {
	//	c.JSON(200, payload.FailPayload("获取用户信息失败"))
	//	return
	//}

	data := ResUserInfo{
		Code:    0,
		Success: true,
		Message: "检查成功",
		Data: &UserInfo{
			//User:  models.LoginUser,
			Token: token,
		},
	}

	c.JSON(200, data)

	//	data := []byte(`{
	//	"code": 0,
	//	"data": {
	//		"available": 0,
	//		"birthday": "",
	//		"email": "",
	//		"industry": "",
	//		"intro": "",
	//		"lastLoginTime": "",
	//		"modifyTime": "",
	//		"modifyUserId": 0,
	//		"openId": "",
	//		"password": "",
	//		"position": "",
	//		"registerTime": "",
	//		"roles": [
	//			{
	//				"available": 0,
	//				"createTime": "",
	//				"createUserId": 0,
	//				"description": "",
	//				"modifyTime": "",
	//				"modifyUserId": 0,
	//				"permissions": [
	//					{
	//						"createTime": "",
	//						"createUserId": 0,
	//						"modifyTime": "",
	//						"modifyUserId": 0,
	//						"orderNum": 0,
	//						"parentId": 0,
	//						"permissionCode": "",
	//						"permissionId": 0,
	//						"permissionName": "",
	//						"resourceType": 0
	//					}
	//				],
	//				"roleId": 0,
	//				"roleName": ""
	//			}
	//		],
	//		"salt": "",
	//		"phone": "",
	//		"token": "",
	//		"userId": 0,
	//		"username": "",
	//		"verificationCode": ""
	//	},
	//	"message": "成功"
	//}`)
	//
	//	js, _ := simplejson.NewJson(data)
	//	c.JSON(200, js)
}

func Login(c *gin.Context) {
	phone := c.Query("phone")
	password := c.Query("password")

	// 查询是否有该用户
	userInfo, err := user_service.GetUserByPhone(phone)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(200, payload.FailPayload("该用户不存在"))
			return
		}
		c.JSON(200, payload.FailPayload(fmt.Sprintf("出错了：%v", err)))
		return
	}

	// 查询密码是否正确
	encPassword := util.Sha1([]byte(password + pwdSalt))
	err = user_service.CheckPassword(phone, encPassword)
	if err != nil {
		log.Errorf("登陆失败:%v", err)
		c.JSON(200, payload.FailPayload(fmt.Sprintf("登陆失败：%v", err)))
		return
	}

	token, err := token_service.MakeToken(phone)
	if err != nil {
		log.Errorf("生成token失败%v", err)
		return
	}
	err = token_service.UpdateToken(phone, token)
	if err != nil {
		log.Errorf("更新token失败:%v", err)
		c.JSON(200, payload.FailPayload(fmt.Sprintf("更新token失败:%v", err)))
		return
	}

	data := ResUserInfo{
		Code:    0,
		Success: true,
		Message: "登陆成功",
		Data: &UserInfo{
			Token: token,
		},
	}

	c.JSON(200, data)
	// 缓存用户信息
	defer cache_service.AddUserCache(userInfo)

	//data := []byte(`{
	//	"code": 0,
	//	"data": {
	//		"email": "116****483@qq.com",
	//		"lastLoginTime": "2019-12-23 14:21:52",
	//		"registerTime": "2019-12-23 14:21:52",
	//		"phone": "187****1817",
	//		"token": "",
	//		"userId": 1,
	//		"username": "奇文网盘"
	//	},
	//	"message": "成功",
	//	"success": true
	//}`)
	//js, _ := simplejson.NewJson(data)
	//c.JSON(200, js)

}

func GetUserInfo(c *gin.Context) {
	c.JSON(200, gin.H{
		"id":           1,
		"username":     "richard",
		"email":        "974102233@qq.com",
		"imgUrl":       "http://xxx.apple.png",
		"registerTime": "2021-12-19T11:29:42.241Z",
		"role":         "admin",
	})
}

//type Dd1 struct {
//	Id         string    `json:"id"`
//	Filename   string    `json:"filename"`
//	Extension  string    `json:"extension"`
//	Pid        string    `json:"pid"`
//	FileSize   int       `json:"fileSize"`
//	UpdateTime time.Time `json:"updateTime"`
//	CreateTime time.Time `json:"createTime"`
//	DeleteTime time.Time `json:"deleteTime"`
//	Dir        bool      `json:"dir"`
//}
//
//type Dd struct {
//	Status int    `json:"status"`
//	Msg    string `json:"msg"`
//	UserInfo   []Dd1  `json:"data"`
//}

func GetFileList(c *gin.Context) {
	myClaims, err := token_service.CheckToken(c)
	if err != nil {
		c.JSON(200, payload.FailPayload("token无效"))
	}
	filePath, page, pageCount, fileType := c.DefaultQuery("filePath", "/"), c.DefaultQuery("currentPage", "1"), //fileType: 0为全部文件，1、2、3、4、5分别对应图片，视频，文档，音乐，其他
		c.DefaultQuery("pageCount", "50"), c.DefaultQuery("fileType", "0")
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

func GetFilePathTree(c *gin.Context) {
	c.JSON(200, gin.H{})
}

func GetUserStorageInfo(c *gin.Context) {
	c.JSON(200, gin.H{
		"maxStorage":  0,
		"usedStorage": 0,
		"overflow":    true,
	})
}
