package main

import (
	"go-gin-gorm/config"
	"go-gin-gorm/models"
	"go-gin-gorm/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	config.Connect()
	config.DB.AutoMigrate(&models.Book{})

	r := gin.Default()
	routes.SetupRoutes(r)
	r.Run()
}
