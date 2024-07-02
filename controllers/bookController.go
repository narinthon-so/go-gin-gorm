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
		config.Logger.Error("Error fetching books: ", err)
		c.Error(err)
		return
	}
	config.Logger.Info("Fetched all books")
	c.JSON(http.StatusOK, books)
}

func GetBook(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var book models.Book
	if result := config.DB.First(&book, id); result.Error != nil {
		config.Logger.Error("Error fetching book: ", result.Error)
		c.Error(result.Error).SetType(gin.ErrorTypePublic).SetMeta("Book not found")
		return
	}
	config.Logger.Infof("Fetched book with ID %d", id)
	c.JSON(http.StatusOK, book)
}

func CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		config.Logger.Error("Error binding JSON: ", err)
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	if err := config.DB.Create(&book).Error; err != nil {
		config.Logger.Error("Error creating book: ", err)
		c.Error(err)
		return
	}
	config.Logger.Info("Created new book: ", book)
	c.JSON(http.StatusCreated, book)
}

func UpdateBook(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var book models.Book
	if result := config.DB.First(&book, id); result.Error != nil {
		config.Logger.Error("Error fetching book: ", result.Error)
		c.Error(result.Error).SetType(gin.ErrorTypePublic).SetMeta("Book not found")
		return
	}

	if err := c.ShouldBindJSON(&book); err != nil {
		config.Logger.Error("Error binding JSON: ", err)
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := config.DB.Save(&book).Error; err != nil {
		config.Logger.Error("Error updating book: ", err)
		c.Error(err)
		return
	}
	config.Logger.Infof("Updated book with ID %d", id)
	c.JSON(http.StatusOK, book)
}

func DeleteBook(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var book models.Book
	if result := config.DB.First(&book, id); result.Error != nil {
		config.Logger.Error("Error fetching book: ", result.Error)
		c.Error(result.Error).SetType(gin.ErrorTypePublic).SetMeta("Book not found")
		return
	}

	if err := config.DB.Delete(&book).Error; err != nil {
		config.Logger.Error("Error deleting book: ", err)
		c.Error(err)
		return
	}
	config.Logger.Infof("Deleted book with ID %d", id)
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
}
