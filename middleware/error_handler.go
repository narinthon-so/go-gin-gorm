package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Only proceed if there are errors to handle
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			var statusCode int

			// Determine the status code based on the type of error
			switch err.Type {
			case gin.ErrorTypeBind:
				statusCode = http.StatusBadRequest
			case gin.ErrorTypePublic:
				statusCode = http.StatusBadRequest
			default:
				statusCode = http.StatusInternalServerError
			}

			// Respond with JSON error message
			c.JSON(statusCode, gin.H{
				"error": err.Error(),
			})
		}
	}
}
