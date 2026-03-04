package util

import "github.com/gin-gonic/gin"

// ErrorResponse sends a standardized JSON error.
func ErrorResponse(c *gin.Context, status int, code string, message string) {
	c.JSON(status, gin.H{
		"code":    code,
		"message": message,
	})
}
