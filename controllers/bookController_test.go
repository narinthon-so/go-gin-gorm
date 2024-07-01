package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"go-gin-gorm/config"
	"go-gin-gorm/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// SetupRouter sets up a Gin router for testing
func SetupRouter() *gin.Engine {
	router := gin.Default()
	return router
}

// TestMain sets up the environment for tests
func TestMain(m *testing.M) {
	// Determine the absolute path to the project root
	projectRoot, _ := filepath.Abs("../")
	fmt.Println("Project root:", projectRoot)

	// Path to .env.test file
	envPath := filepath.Join(projectRoot, ".env.test")
	fmt.Println("Loading .env.test file from:", envPath)

	// Debugging: Check if .env.test file exists
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		fmt.Println(".env.test file does not exist")
	} else {
		fmt.Println(".env.test file exists")
	}

	err := godotenv.Load(envPath)
	if err != nil {
		panic("Error loading .env.test file")
	}
	config.Connect()
	config.DB.AutoMigrate(&models.Book{})
	code := m.Run()
	os.Exit(code)
}

// TestGetBooks tests the GetBooks function
func TestGetBooks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := SetupRouter()

	// Seed the database
	book := models.Book{Title: "Test Book", Author: "Test Author"}
	config.DB.Create(&book)

	// Define the route
	router.GET("/api/books", GetBooks)

	// Create a request to send to the above route
	req, _ := http.NewRequest("GET", "/api/books", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	var books []models.Book
	err := json.Unmarshal(w.Body.Bytes(), &books)
	assert.NoError(t, err)
	assert.NotEmpty(t, books)
}

// TestGetBook tests the GetBook function
func TestGetBook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := SetupRouter()

	// Seed the database
	book := models.Book{Title: "Test Book", Author: "Test Author"}
	config.DB.Create(&book)

	// Define the route
	router.GET("/api/books/:id", GetBook)

	// Create a request to send to the above route
	req, _ := http.NewRequest("GET", "/api/books/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	var fetchedBook models.Book
	err := json.Unmarshal(w.Body.Bytes(), &fetchedBook)
	assert.NoError(t, err)
	assert.Equal(t, book.Title, fetchedBook.Title)
	assert.Equal(t, book.Author, fetchedBook.Author)
}

// TestCreateBook tests the CreateBook function
func TestCreateBook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := SetupRouter()

	// Define the route
	router.POST("/api/books", CreateBook)

	// Create a request to send to the above route
	book := models.Book{Title: "New Book", Author: "New Author"}
	jsonValue, _ := json.Marshal(book)
	req, _ := http.NewRequest("POST", "/api/books", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusCreated, w.Code)
	var createdBook models.Book
	err := json.Unmarshal(w.Body.Bytes(), &createdBook)
	assert.NoError(t, err)
	assert.Equal(t, book.Title, createdBook.Title)
	assert.Equal(t, book.Author, createdBook.Author)
}

// TestUpdateBook tests the UpdateBook function
func TestUpdateBook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := SetupRouter()

	// Seed the database
	book := models.Book{Title: "Test Book", Author: "Test Author"}
	config.DB.Create(&book)

	// Define the route
	router.PUT("/api/books/:id", UpdateBook)

	// Create a request to send to the above route
	updatedBook := models.Book{Title: "Updated Book", Author: "Updated Author"}
	jsonValue, _ := json.Marshal(updatedBook)
	req, _ := http.NewRequest("PUT", "/api/books/1", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	var fetchedBook models.Book
	err := json.Unmarshal(w.Body.Bytes(), &fetchedBook)
	assert.NoError(t, err)
	assert.Equal(t, updatedBook.Title, fetchedBook.Title)
	assert.Equal(t, updatedBook.Author, fetchedBook.Author)
}

// TestDeleteBook tests the DeleteBook function
func TestDeleteBook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := SetupRouter()

	// Seed the database
	book := models.Book{Title: "Test Book", Author: "Test Author"}
	config.DB.Create(&book)

	// Define the route
	router.DELETE("/api/books/:id", DeleteBook)

	// Create a request to send to the above route
	req, _ := http.NewRequest("DELETE", "/api/books/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Book deleted", response["message"])
}
