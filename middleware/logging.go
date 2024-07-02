package middleware

import (
	"bytes"
	"go-gin-gorm/config"
	"go-gin-gorm/models"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LoggerToDB() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Capture request body
		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				requestBody = string(bodyBytes)
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				config.Logger.Debug("Captured Request Body: ", requestBody)
			} else {
				config.Logger.Error("Failed to read request body: ", err)
			}
		}

		// Use custom response writer
		bw := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bw

		// Process request
		c.Next()

		latency := time.Since(start).Milliseconds()
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		userAgent := c.Request.UserAgent()
		responseBody := bw.body.String()

		log := models.RequestLog{
			Method:       method,
			Path:         path,
			StatusCode:   statusCode,
			Latency:      latency,
			ClientIP:     clientIP,
			UserAgent:    userAgent,
			RequestBody:  requestBody,
			ResponseBody: responseBody,
		}

		config.Logger.Debug("Request Log to be saved: ", log)

		if err := config.DB.Create(&log).Error; err != nil {
			config.Logger.Error("Failed to log request: ", err)
		}
	}
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		stop := time.Now()
		latency := stop.Sub(start)

		// Get the status we sent
		status := c.Writer.Status()

		entry := config.Logger.WithFields(map[string]interface{}{
			"status":     status,
			"latency":    latency,
			"ip":         c.ClientIP(),
			"method":     c.Request.Method,
			"path":       path,
			"query":      raw,
			"user-agent": c.Request.UserAgent(),
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.String())
		} else {
			entry.Info("Request completed")
		}
	}
}
