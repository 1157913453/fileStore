package db

import (
	"errors"
	"filestore/config"
	"filestore/middleware/token"
	"filestore/models"
	"filestore/payload"
	"filestore/service/cache_service"
	"filestore/service/user_service"
	"filestore/util"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io/fs"
	"os"
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
	phone, userName, password := register.Phone, register.UserName, register.Password
	if len(phone) < 11 || len(password) < 5 {
		log.Errorf("手机号或密码太短")
		c.JSON(200, payload.FailPayload("手机号或密码太短"))
		return
	}

	// 对密码加密
	encPassword, err := util.HashAndSalt([]byte(password))
	if err != nil {
		log.Errorf("密码加密失败：%v", err)
		c.JSON(200, payload.FailPayload("密码加密失败"))
		return
	}
	//encPassword := util.Sha1([]byte(password + pwdSalt))
	err = user_service.CreateUser(phone, encPassword, userName)
	if err != nil {
		log.Errorf("注册失败：%v", err)
		c.JSON(200, payload.FailPayload("注册失败,已有该用户名"))
		return
	}

	path := config.BasePath + phone
	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		if !errors.Is(err, fs.ErrExist) {
			log.Errorf("创建%s用户文件夹失败:%v", phone, err)
			c.JSON(200, payload.FailPayload("创建用户文件夹失败"))
			return
		}
	}

	log.Infof("创建用户%s文件夹成功", phone)

	c.JSON(200, payload.SucPayload("注册成功"))
	defer func() {
		userInfo, err := user_service.GetUserByPhone(phone)
		if err != nil {
			log.Errorf("查询%s用户信息失败:%v", phone, err)
			return
		}
		err = cache_service.AddUserCache(userInfo)
		if err != nil {
			return
		}
	}()
}

func CheckUserLoginInfo(c *gin.Context) {
	token := c.GetHeader("token")

	Phone, exists := c.Get("Phone")
	if !exists {
		_, err := middleware.ParseToken(token)
		if err != nil {
			log.Errorf("token Err:%v", err)
			return
		}

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
		return
	}
	phone := Phone.(string)
	// 判断缓存是否存在
	_, err := cache_service.GetUserCache(phone)
	if err != nil {
		// 设置缓存
		userInfo, err := user_service.GetUserByPhone(phone)
		if err != nil {
			log.Errorf("查询用户信息失败:%v", err)
			c.JSON(200, payload.FailPayload("查询用户信息失败"))
			return
		}
		err = cache_service.AddUserCache(userInfo)
		if err != nil {
			log.Errorf("设置用户缓存失败:%v", err)
			c.JSON(200, payload.FailPayload("设置用户缓存失败"))
			return
		}
	}

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
	//encPassword := util.Sha1([]byte(password + pwdSalt))
	err = user_service.CheckPassword(phone, password)
	if err != nil {
		log.Errorf("登陆失败:%v", err)
		c.JSON(200, payload.FailPayload(fmt.Sprintf("登陆失败：%v", err)))
		return
	}
	c.Set("Phone", phone)

	token, err := middleware.MakeToken(phone)
	if err != nil {
		log.Errorf("生成token失败%v", err)
		return
	}
	//err = token.UpdateToken(phone, token)
	//if err != nil {
	//	log.Errorf("更新token失败:%v", err)
	//	c.JSON(200, payload.FailPayload(fmt.Sprintf("更新token失败:%v", err)))
	//	return
	//}

	data := ResUserInfo{
		Code:    0,
		Success: true,
		Message: "登陆成功",
		Data: &UserInfo{
			Token: token,
		},
	}
	log.Infof("%s登陆成功", phone)
	c.JSON(200, data)
	// 缓存用户信息
	defer cache_service.AddUserCache(userInfo)
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
