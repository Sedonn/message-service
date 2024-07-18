package message

import (
	"log/slog"

	"context"

	"github.com/sedonn/message-service/internal/domain/events"
	"github.com/sedonn/message-service/internal/domain/models"
	"github.com/sedonn/message-service/internal/event/kafka/consumer"
	"github.com/sedonn/message-service/internal/pkg/logger"
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

// MessageUpdater описывает поведение объекта, который обеспечивает обновление данных сообщений.
type MessageUpdater interface {
	// UpdateMessage сохраняет данные существующего сообщения.
	UpdateMessage(ctx context.Context, m models.Message) (models.Message, error)
}

// MessageEventProducer описывает поведение объекта, который обеспечивает отправку событий связанных с данными сообщений.
type MessageEventProducer interface {
	// NotifyStartProcessingMessage создает событие старта обработки сообщения.
	NotifyStartProcessingMessage(e events.StartProcessingMessage) error
}

// Message предоставляет бизнес-логику работы с сообщениями.
type Message struct {
	log                  *slog.Logger
	messageProvider      MessageProvider
	messageSaver         MessageSaver
	messageUpdater       MessageUpdater
	messageEventProducer MessageEventProducer
}

var _ messagerest.Messenger = (*Message)(nil)
var _ consumer.MessageEventSubscriber = (*Message)(nil)

// New создает новый сервис для работы с сообщениями.
func New(log *slog.Logger, mp MessageProvider, ms MessageSaver, mu MessageUpdater, mep MessageEventProducer) *Message {
	return &Message{
		log:                  log,
		messageProvider:      mp,
		messageSaver:         ms,
		messageUpdater:       mu,
		messageEventProducer: mep,
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
		log.Error("failed to get messages", logger.StringError(err))

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
		log.Error("failed to get processed messages", logger.StringError(err))

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
		log.Error("failed to get unprocessed messages", logger.StringError(err))

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
		log.Error("failed to create message", logger.StringError(err))

		return 0, err
	}

	log = log.With(slog.Uint64("message_id", id))
	log.Info("success to create message")

	e := events.StartProcessingMessage{
		ID:      id,
		Content: content,
	}

	if err := m.messageEventProducer.NotifyStartProcessingMessage(e); err != nil {
		log.Error("failed to send message for processing", logger.StringError(err))

		return 0, err
	}

	log.Info("success to send message for processing")

	return id, nil
}

// OnMessageProcessed implements consumer.MessageEventConsumer.
func (m *Message) OnMessageProcessed(ctx context.Context, e events.CompleteProcessingMessage) {
	const op = "message.OnMessageProcessed"
	log := m.log.With(slog.String("op", op), slog.Uint64("message_id", e.ID))

	log.Info("attempt to update processed message")
	_, err := m.messageUpdater.UpdateMessage(ctx, models.Message{
		ID:          e.ID,
		ProcessedAt: &e.ProcessedAt,
	})

	if err != nil {
		log.Error("failed to update processed message", logger.StringError(err))
	}

	log.Info("success to update processed message")
}
