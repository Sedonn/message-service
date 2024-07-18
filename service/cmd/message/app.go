package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/sedonn/message-service/internal/app"
	"github.com/sedonn/message-service/internal/config"
	"github.com/sedonn/message-service/internal/pkg/logger"
)

func main() {
	const op = "message.main"
	ctx, cancel := context.WithCancel(context.Background())

	cfg := config.MustLoad()

	log := logger.New(cfg.Env)
	log.Info("logger initialized", slog.String("op", op), slog.String("env", cfg.Env))

	application := app.New(log, cfg)
	application.EventConsumer.MustRun(ctx)
	go application.RESTApp.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	cancel()
	application.RESTApp.Stop()
	if err := application.EventProducer.Stop(); err != nil {
		log.Error("failed to close event producer", logger.StringError(err))
	}

	application.EventConsumer.Stop()
}
