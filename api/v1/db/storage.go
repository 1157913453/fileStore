package db

import (
	"github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func GetStorage(c *gin.Context) {
	data := []byte(`{
	"code": 0,
	"data": {
		"modifyTime": "",
		"modifyUserId": 1,
		"storageId": 1,
		"storageSize": 10033,
		"totalStorageSize": 1233440000,
		"userId": 1
	},
	"message": "成功",
	"success": true
}`)
	js, err := simplejson.NewJson(data)
	if err != nil {
		log.Errorf("errrrrr是：%v", err)
	}
	c.JSON(200, js)
}
