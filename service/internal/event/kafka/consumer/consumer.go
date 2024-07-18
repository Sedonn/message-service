package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/IBM/sarama"

	"github.com/sedonn/message-service/internal/config"
	"github.com/sedonn/message-service/internal/domain/events"
	"github.com/sedonn/message-service/internal/pkg/logger"
)

// MessageEventSubscriber описывает поведение объекта, который выполняет события связанные с сообщениямиЛ.
type MessageEventSubscriber interface {
	// OnMessageProcessed вызывается при завершении обработки сообщения.
	OnMessageProcessed(ctx context.Context, e events.CompleteProcessingMessage)
}

// Consumer получает сообщения из kafka.
type Consumer struct {
	log                  *slog.Logger
	cfg                  *config.KafkaConfig
	client               sarama.ConsumerGroup
	wg                   *sync.WaitGroup
	ready                chan bool
	messageEventConsumer MessageEventSubscriber
}

var _ sarama.ConsumerGroupHandler = (*Consumer)(nil)

// New создает нового Consumer.
func New(log *slog.Logger, cfg *config.KafkaConfig, mec MessageEventSubscriber) (*Consumer, error) {
	const group = "message-service"

	client, err := sarama.NewConsumerGroup(strings.Split(cfg.Brokers, ","), group, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize kafka consumer: %w", err)
	}

	return &Consumer{
		log:                  log,
		cfg:                  cfg,
		client:               client,
		wg:                   &sync.WaitGroup{},
		ready:                make(chan bool),
		messageEventConsumer: mec,
	}, nil
}

// MustRun запускает Consumer.
func (c *Consumer) MustRun(ctx context.Context) {
	const op = "consumer.MustRun"
	log := c.log.With(slog.String("op", op))

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			if err := c.client.Consume(ctx, []string{c.cfg.Topics.ProcessedMessages}, c); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				panic("consumer client error" + err.Error())
			}

			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-c.ready
	log.Info("kafka consumer start working")
}

// MustRun останавливает Consumer.
func (c *Consumer) Stop() {
	const op = "consumer.Stop"
	log := c.log.With(slog.String("op", op))

	log.Info("closing kafka consumer")
	c.wg.Wait()
	if err := c.client.Close(); err != nil {
		log.Error("failed to close kafka consumer", logger.StringError(err))
	}

	log.Info("kafka consumer closed")
}

// Setup реализует метод ConsumerGroupHandler.Setup.
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

// ConsumeClaim реализует метод sarama.ConsumerGroupHandler.
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	const op = "consumer.ConsumeClaim"
	log := c.log.With(slog.String("op", op))

	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			log.Info("received new message",
				slog.String("message_key", string(msg.Key)),
				slog.String("topic", msg.Topic),
			)

			switch msg.Topic {
			case c.cfg.Topics.ProcessedMessages:
				c.consumeMessageProcessedEvent(session.Context(), msg)
			}
		case <-session.Context().Done():
			return nil
		}
	}
}

// Cleanup реализует метод sarama.ConsumerGroupHandler.
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// consumeMessageProcessedEvent передает полученное событие о завершении обработки сообщения в подписчика.
func (c *Consumer) consumeMessageProcessedEvent(ctx context.Context, msg *sarama.ConsumerMessage) {
	const op = "consumer.consumeMessageProcessedEvent"
	log := c.log.With(slog.String("op", op), slog.String("message_key", string(msg.Key)))

	var e events.CompleteProcessingMessage
	if err := json.Unmarshal(msg.Value, &e); err != nil {
		log.Error("failed to unmarshal message value", logger.StringError(err))
		return
	}

	e.ProcessedAt = time.Now()

	c.messageEventConsumer.OnMessageProcessed(ctx, e)
}
