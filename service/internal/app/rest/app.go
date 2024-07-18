package restapp

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/sedonn/message-service/internal/pkg/logger"
	messagerest "github.com/sedonn/message-service/internal/rest/handlers/message"
	mwerror "github.com/sedonn/message-service/internal/rest/middleware/error"
)

// App это REST-сервер.
type App struct {
	log        *slog.Logger
	httpServer *http.Server
	port       int
}

// New создает новый REST-сервер.
func New(log *slog.Logger, port int, m messagerest.Messenger) *App {
	router := gin.Default()

	router.Use(mwerror.New())

	v1 := router.Group("/v1")
	{
		messagerest.New(m).BindTo(v1)
	}

	srv := &http.Server{
		Addr:    net.JoinHostPort("", strconv.Itoa(port)),
		Handler: router.Handler(),
	}

	return &App{
		log:        log,
		httpServer: srv,
		port:       port,
	}
}

// MustRun запускает REST-API сервер. Паникует при ошибке.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run запускает REST-API сервер.
func (a *App) Run() error {
	const op = "restapp.Run"
	log := a.log.With(slog.String("op", op))

	log.Info("starting REST-API server", slog.String("address", a.httpServer.Addr))
	if err := a.httpServer.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}

	return nil
}

// Stop останавливает REST-API сервер.
func (a *App) Stop() {
	const op = "restapp.Stop"
	log := a.log.With(slog.String("op", op), slog.String("address", a.httpServer.Addr))

	log.Info("shutting down REST-API server")
	if err := a.httpServer.Shutdown(context.Background()); err != nil {
		log.Error("failed to shut down REST-API server", logger.StringError(err))
	}

	log.Info("REST-API server is shut down")
}
