package cache_service

import (
	"filestore/src/models"
	json "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"time"
)

func AddUserCache(userInfo *models.User) (err error) {
	userB, err := json.Marshal(userInfo)
	if err != nil {
		log.Errorf("userInfo序列化失败：%v", err)
		return
	}

	err = Rdb.SetNX(Ctx, "user_phone:"+userInfo.Phone, string(userB), 4*time.Hour).Err()
	if err != nil {
		log.Errorf("添加用户缓存失败：%v", err)
		return
	}

	return
}

func GetUserCache(phone string) (userInfo *models.User, err error) {
	res, err := Rdb.Get(Ctx, "user_phone:"+phone).Result()
	if err != nil {
		log.Errorf("获取%s缓存错误：%v", phone, err)
		return nil, err
	}
	userB := []byte(res)
	err = json.Unmarshal(userB, &userInfo)
	if err != nil {
		log.Errorf("反序列化用户%s失败%v", phone, err)
		return
	}
	return userInfo, nil
}
