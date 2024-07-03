package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"go-gin-gorm/config"
	"go-gin-gorm/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/api/books", GetBooks)
	router.GET("/api/books/:id", GetBook)
	router.POST("/api/books", CreateBook)
	router.PUT("/api/books/:id", UpdateBook)
	router.DELETE("/api/books/:id", DeleteBook)
	return router
}

func generateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "admin",
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte("test_secret"))
}

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env.test")
	if err != nil {
		panic("Error loading .env.test file: " + err.Error())
	}

	config.Connect()
	config.DB.AutoMigrate(&models.Book{})

	code := m.Run()

	config.DB.Migrator().DropTable(&models.Book{})

	os.Exit(code)
}

func TestGetBooks(t *testing.T) {
	r := SetupRouter()
	token, _ := generateToken()
	req, _ := http.NewRequest("GET", "/api/books", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestCreateBook(t *testing.T) {
	r := SetupRouter()
	token, _ := generateToken()
	book := models.Book{
		Title:  "Test Book",
		Author: "Test Author",
	}
	jsonValue, _ := json.Marshal(book)
	req, _ := http.NewRequest("POST", "/api/books", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

func TestCreateBookInvalid(t *testing.T) {
	r := SetupRouter()
	token, _ := generateToken()
	book := map[string]string{
		"author": "Test Author",
	}
	jsonValue, _ := json.Marshal(book)
	req, _ := http.NewRequest("POST", "/api/books", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestGetBookNotFound(t *testing.T) {
	r := SetupRouter()
	token, _ := generateToken()
	req, _ := http.NewRequest("GET", "/api/books/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestUpdateBook(t *testing.T) {
	r := SetupRouter()
	token, _ := generateToken()
	book := models.Book{
		Title:  "Updated Book",
		Author: "Updated Author",
	}
	jsonValue, _ := json.Marshal(book)
	req, _ := http.NewRequest("PUT", "/api/books/1", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestDeleteBook(t *testing.T) {
	r := SetupRouter()
	token, _ := generateToken()
	req, _ := http.NewRequest("DELETE", "/api/books/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}
