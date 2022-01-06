package main

import (
	"filestore/api"
	"filestore/models"
	"filestore/service/cache_service"
)

func main() {
	models.InitDB()
	cache_service.InitCache()
	api.InitRouter()
}
