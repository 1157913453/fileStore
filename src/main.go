package main

import (
	"filestore/src/api"
	"filestore/src/models"
	"filestore/src/service/cache_service"
)

func main() {
	models.InitDB()
	cache_service.InitCache()
	api.InitRouter()
}
