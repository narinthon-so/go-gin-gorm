package routes

import (
	"go-gin-gorm/controllers"
	"go-gin-gorm/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	api.Use(middleware.Authenticate()) // Apply JWT middleware to protected routes
	{
		api.GET("/books", controllers.GetBooks)
		api.GET("/books/:id", controllers.GetBook)
		api.POST("/books", controllers.CreateBook)
		api.PUT("/books/:id", controllers.UpdateBook)
		api.DELETE("/books/:id", controllers.DeleteBook)
	}
	auth := router.Group("/auth")
	{
		auth.POST("/login", controllers.Login)
	}
}
