package producer

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/IBM/sarama"
	"github.com/google/uuid"

	"github.com/sedonn/message-service/internal/config"
	"github.com/sedonn/message-service/internal/domain/events"
	"github.com/sedonn/message-service/internal/services/message"
)

// Producer отправляет сообщения в Kafka.
type Producer struct {
	cfg *config.KafkaConfig
	sp  sarama.SyncProducer
}

var _ message.MessageEventProducer = (*Producer)(nil)

// New создает нового Producer.
func New(cfg *config.KafkaConfig) (*Producer, error) {
	sp, err := sarama.NewSyncProducer(strings.Split(cfg.Brokers, ","), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize kafka producer: %w", err)
	}

	return &Producer{
		cfg: cfg,
		sp:  sp,
	}, nil
}

// Stop закрывает подключение Producer.
func (p *Producer) Stop() error { return p.sp.Close() }

// NotifyStartProcessingMessage создает событие старта обработки сообщения.
func (p *Producer) NotifyStartProcessingMessage(e events.StartProcessingMessage) error {
	return p.sendMessage(p.cfg.Topics.ProcessingMessages, e)
}

// sendMessage обертка для отправки событий в Kafka.
func (p *Producer) sendMessage(topic string, payload any) error {
	requestID := uuid.New().String()

	pBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal event payload: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(requestID),
		Value: sarama.ByteEncoder(pBytes),
	}

	if _, _, err := p.sp.SendMessage(msg); err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	return nil
}
