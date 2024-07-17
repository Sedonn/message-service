package message

import (
	"log/slog"

	"context"

	"github.com/sedonn/message-service/internal/domain/models"
	messagerest "github.com/sedonn/message-service/internal/rest/handlers/message"
)

// MessageProvider описывает поведение объекта, который обеспечивает получение данных сообщений.
type MessageProvider interface {
	// Messages возвращает данные о всех сообщениях.
	Messages(ctx context.Context, limit uint, offset uint) ([]models.Message, error)
	// ProcessedMessages возвращает данные только обработанных сообщений.
	ProcessedMessages(ctx context.Context, limit uint, offset uint) ([]models.Message, error)
	// UnprocessedMessages возвращает данные только необработанных сообщений.
	UnprocessedMessages(ctx context.Context, limit uint, offset uint) ([]models.Message, error)
}

// MessageSaver описывает поведение объекта, который обеспечивает сохранение данных сообщений.
type MessageSaver interface {
	// SaveMessage сохраняет данные нового сообщения.
	SaveMessage(ctx context.Context, m models.Message) (uint64, error)
}

// Message предоставляет бизнес-логику работы с сообщениями.
type Message struct {
	log             *slog.Logger
	messageProvider MessageProvider
	messageSaver    MessageSaver
}

var _ messagerest.Messenger = (*Message)(nil)

// New создает новый сервис для работы с сообщениями.
func New(log *slog.Logger, mp MessageProvider, ms MessageSaver) *Message {
	return &Message{
		log:             log,
		messageProvider: mp,
		messageSaver:    ms,
	}
}

// GetMessages получает все сообщения.
func (m *Message) GetMessages(ctx context.Context, pageID uint) ([]models.Message, error) {
	const (
		op       = "message.GetMessages"
		pageSize = 10
	)
	log := m.log.With(slog.String("op", op))

	log.Info("attempt to get messages", slog.Int("page_size", pageSize))

	messages, err := m.messageProvider.Messages(ctx, pageSize, pageID*pageSize)
	if err != nil {
		log.Error("failed to get messages", slog.String("err", err.Error()))

		return nil, err
	}

	log.Info("success to get messages", slog.Int("messages_count", len(messages)))

	return messages, nil
}

// GetProcessedMessages получает только обработанные сообщения.
func (m *Message) GetProcessedMessages(ctx context.Context, pageID uint) ([]models.Message, error) {
	const (
		op       = "message.GetProcessedMessages"
		pageSize = 10
	)
	log := m.log.With(slog.String("op", op))

	log.Info("attempt to get processed messages", slog.Int("page_size", pageSize))

	messages, err := m.messageProvider.ProcessedMessages(ctx, pageSize, pageID*pageSize)
	if err != nil {
		log.Error("failed to get processed messages", slog.String("err", err.Error()))

		return nil, err
	}

	log.Info("success to get processed messages", slog.Int("messages_count", len(messages)))

	return messages, nil
}

// GetUnprocessedMessages получает только необработанные сообщения.
func (m *Message) GetUnprocessedMessages(ctx context.Context, pageID uint) ([]models.Message, error) {
	const (
		op       = "message.GetUnprocessedMessages"
		pageSize = 10
	)
	log := m.log.With(slog.String("op", op))

	log.Info("attempt to get unprocessed messages", slog.Int("page_size", pageSize))

	messages, err := m.messageProvider.UnprocessedMessages(ctx, pageSize, pageID*pageSize)
	if err != nil {
		log.Error("failed to get unprocessed messages", slog.String("err", err.Error()))

		return nil, err
	}

	log.Info("success to get unprocessed messages", slog.Int("messages_count", len(messages)))

	return messages, nil
}

// CreateMessage создает новое сообщение.
func (m *Message) CreateMessage(ctx context.Context, content string) (uint64, error) {
	const op = "message.GetMessages"
	log := m.log.With(slog.String("op", op))

	log.Info("attempt to create message", slog.Int("message_size", len(content)))

	id, err := m.messageSaver.SaveMessage(ctx, models.Message{Content: content})
	if err != nil {
		log.Info("failed to create message", slog.String("err", err.Error()))

		return 0, err
	}

	log.Info("success to create message", slog.Uint64("message_id", id))

	return id, nil
}
