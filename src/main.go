package main

import (
	"filestore/src/api"
	"filestore/src/models"
)

func main() {
	models.InitDB()
	api.InitRouter()
}
