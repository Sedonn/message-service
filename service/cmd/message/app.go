package main

import (
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

	cfg := config.MustLoad()

	log := logger.New(cfg.Env)
	log.Info("logger initialized", slog.String("op", op), slog.String("env", cfg.Env))

	application := app.New(log, cfg)
	go application.RESTApp.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	application.RESTApp.Stop()
	log.Info("REST-API server is shut down", slog.String("op", op))
}
