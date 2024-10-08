package get

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sedonn/message-service/internal/domain/models"
)

// MessageCreator описывает поведение объекта, который извлекает и фильтрует данные сообщений.
type MessageGetter interface {
	// GetMessages получает все сообщения.
	GetMessages(ctx context.Context, pageID uint) ([]models.Message, error)

	// GetProcessedMessages получает только обработанные сообщения.
	GetProcessedMessages(ctx context.Context, pageID uint) ([]models.Message, error)

	// GetUnprocessedMessages получает только необработанные сообщения.
	GetUnprocessedMessages(ctx context.Context, pageID uint) ([]models.Message, error)
}

type request struct {
	PageID    *uint `form:"page,default=0" binding:"numeric,gte=0"`
	Processed *bool `form:"processed" binding:"omitempty,boolean"`
}

type response []models.Message

// New возвращает новый хендлер, который извлекает данные сообщений.
//
//	@Summary		Получить сообщения
//	@Description	Получение сообщений.
//	@Tags			messages
//	@Accept			json
//	@Produce		json
//	@Param			page		query		uint	false	"Номер страницы. Если пуст - 0"
//	@Param			processed	query		bool	false	"Статус - обработано. Если пусто - выводит все сообщения"
//	@Success		200			{array}		models.Message
//	@Failure		400			{object}	mwerror.ErrorResponse
//	@Failure		404			{object}	mwerror.ErrorResponse
//	@Failure		500			{object}	mwerror.ErrorResponse
//	@Router			/messages [get]
func New(m MessageGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req request
		if err := c.ShouldBindQuery(&req); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		var (
			messages []models.Message
			err      error
		)
		switch {
		case req.Processed == nil:
			messages, err = m.GetMessages(c, *req.PageID)
		case *req.Processed:
			messages, err = m.GetProcessedMessages(c, *req.PageID)
		case !*req.Processed:
			messages, err = m.GetUnprocessedMessages(c, *req.PageID)
		}

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, response(messages))
	}
}
