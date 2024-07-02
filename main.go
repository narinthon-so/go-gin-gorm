package main

import (
	"go-gin-gorm/config"
	"go-gin-gorm/middleware"
	"go-gin-gorm/models"
	"go-gin-gorm/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	config.SetupLogger()
	config.Connect()
	config.DB.AutoMigrate(&models.Book{}, &models.RequestLog{}) // Include RequestLog in AutoMigrate

	r := gin.Default()
	r.Use(middleware.LoggerToDB())   // Apply the logging to DB middleware
	r.Use(middleware.Logger())       // Apply the logging middleware
	r.Use(middleware.ErrorHandler()) // Apply the custom error handler middleware
	routes.SetupRoutes(r)
	r.Run()
}
