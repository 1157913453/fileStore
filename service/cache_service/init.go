package cache_service

import (
	"context"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

func InitClient() (err error) {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		log.Errorf("redis 连接失败：%v", err)
		return err
	}
	log.Infof("Redis连接成功:%s", pong)

	return nil
}

func InitCache() {
	err := InitClient()
	if err != nil {
		panic(err)
	}
}
