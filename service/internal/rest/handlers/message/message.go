package messagerest

import (
	"github.com/gin-gonic/gin"

	"github.com/sedonn/message-service/internal/rest/handlers/message/create"
	"github.com/sedonn/message-service/internal/rest/handlers/message/get"
)

// Messenger описывает поведение объекта, который обеспечивает бизнес-логику работы с сообщениями.
type Messenger interface {
	get.MessageGetter
	create.MessageCreator
}

// Handler это корневой хендлер сообщений.
type Handler struct {
	messenger Messenger
}

// New создает новый корневой хендлер сообщений.
func New(m Messenger) *Handler {
	return &Handler{
		messenger: m,
	}
}

// BindTo привязывает хендлер к определенной группе маршрутов.
func (h *Handler) BindTo(router *gin.RouterGroup) {
	message := router.Group("/messages")
	{
		message.GET("/", get.New(h.messenger))
		message.POST("/", create.New(h.messenger))
	}
}
