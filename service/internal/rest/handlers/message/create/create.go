package create

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// MessageCreator описывает поведение объекта, который создает новые сообщения.
type MessageCreator interface {
	CreateMessage(ctx context.Context, content string) (uint64, error)
}

type request struct {
	Content string `json:"content" binding:"required,lte=256"`
}

type response struct {
	ID uint64 `json:"id"`
}

// New возвращает новый объект хендлера, который сохраняет данные сообщения.
//
//	@Summary		Создать сообщение
//	@Description	Создание нового сообщения
//	@Tags			messages
//	@Accept			json
//	@Produce		json
//	@Param			message	body		request	true	"Содержимое сообщения"
//	@Success		200		{object}	response
//	@Failure		400		{object}	mwerror.ErrorResponse
//	@Failure		404		{object}	mwerror.ErrorResponse
//	@Failure		500		{object}	mwerror.ErrorResponse
//	@Router			/messages [post]
func New(m MessageCreator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		id, err := m.CreateMessage(c, req.Content)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, response{ID: id})
	}
}
