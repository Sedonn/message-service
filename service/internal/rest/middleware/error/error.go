package mwerror

import (
	"github.com/gin-gonic/gin"
)

// New создает middleware для глобальной обработки ошибок.
func New() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			c.JSON(-1, gin.H{"error": c.Errors[0]})
		}
	}
}
