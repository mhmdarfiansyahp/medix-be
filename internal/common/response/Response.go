package response

import (
	"github.com/gin-gonic/gin"
)

// Success standardizes success response to {status, message, data}
func Success(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, gin.H{
		"status":  "success",
		"message": message,
		"data":    data,
	})
}

// Error standardizes error response to {status, message, data}
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"status":  "error",
		"message": message,
		"data":    nil,
	})
}