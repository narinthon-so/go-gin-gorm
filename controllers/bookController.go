package controllers

import (
	"net/http"
	"strconv"

	"go-gin-gorm/config"
	"go-gin-gorm/models"

	"github.com/gin-gonic/gin"
)

func GetBooks(c *gin.Context) {
	var books []models.Book
	if err := config.DB.Find(&books).Error; err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, books)
}

func GetBook(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var book models.Book
	if result := config.DB.First(&book, id); result.Error != nil {
		c.Error(result.Error).SetType(gin.ErrorTypePublic).SetMeta("Book not found")
		return
	}
	c.JSON(http.StatusOK, book)
}

func CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	if err := config.DB.Create(&book).Error; err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, book)
}

func UpdateBook(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var book models.Book
	if result := config.DB.First(&book, id); result.Error != nil {
		c.Error(result.Error).SetType(gin.ErrorTypePublic).SetMeta("Book not found")
		return
	}

	if err := c.ShouldBindJSON(&book); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := config.DB.Save(&book).Error; err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, book)
}

func DeleteBook(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var book models.Book
	if result := config.DB.First(&book, id); result.Error != nil {
		c.Error(result.Error).SetType(gin.ErrorTypePublic).SetMeta("Book not found")
		return
	}

	if err := config.DB.Delete(&book).Error; err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
}
