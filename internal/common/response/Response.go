package response

import (
	"github.com/gin-gonic/gin"
)

// Success menyeragamkan response sukses menjadi {status, message, data}
func Success(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, gin.H{
		"status":  "success",
		"message": message,
		"data":    data,
	})
}

// Error menyeragamkan response error menjadi {status, message, data}
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"status":  "error",
		"message": message,
		"data":    nil,
	})
}