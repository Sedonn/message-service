package mwerror

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse хранит данные ошибки для API ответа.
type ErrorResponse struct {
	Err string `json:"error"`
}

// New создает middleware для глобальной обработки ошибок.
func New() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		if c.Writer.Status() == http.StatusInternalServerError {
			c.JSON(-1, ErrorResponse{Err: "internal server error"})
			return
		}

		c.JSON(-1, ErrorResponse{Err: c.Errors[0].Error()})
	}
}
