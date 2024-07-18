package app

import (
	"log/slog"

	restapp "github.com/sedonn/message-service/internal/app/rest"
	"github.com/sedonn/message-service/internal/config"
	"github.com/sedonn/message-service/internal/event/kafka/consumer"
	"github.com/sedonn/message-service/internal/event/kafka/producer"
	"github.com/sedonn/message-service/internal/repository/postgresql"
	"github.com/sedonn/message-service/internal/services/message"
)

// App это микросервис сообщений.
type App struct {
	RESTApp       *restapp.App
	EventProducer *producer.Producer
	EventConsumer *consumer.Consumer
}

// New создает новый микросервис сообщений.
func New(log *slog.Logger, cfg *config.Config) *App {
	const op = "app.New"

	repository, err := postgresql.New(cfg)
	if err != nil {
		panic(err)
	}
	log.Info("database connected", slog.String("op", op), slog.String("database", cfg.DB.Database))

	producer, err := producer.New(&cfg.Kafka)
	if err != nil {
		panic(err)
	}

	messageService := message.New(log, repository, repository, repository, producer)

	consumer, err := consumer.New(log, &cfg.Kafka, messageService)
	if err != nil {
		panic(err)
	}

	restApp := restapp.New(log, cfg.REST.Port, messageService)

	return &App{
		RESTApp:       restApp,
		EventProducer: producer,
		EventConsumer: consumer,
	}
}
