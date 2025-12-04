package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/GeorgeTyupin/labguard/internal/server/config"
	"github.com/GeorgeTyupin/labguard/internal/server/handlers"
	"github.com/go-chi/chi/v5"
)

type ServerApp struct {
	AppName         string
	server          *http.Server
	logger          *slog.Logger
	shutdownTimeout time.Duration
}

func NewServerApp(logger *slog.Logger, cfg *config.Config) *ServerApp {
	appName := "HTTP Server"
	logger = logger.With(slog.String("app", appName))

	handler := registerHandlers()

	server := &http.Server{
		Addr:    cfg.Server.Address,
		Handler: handler,
	}

	application := &ServerApp{
		AppName:         appName,
		server:          server,
		logger:          logger,
		shutdownTimeout: cfg.Server.Timeouts.Shutdown,
	}

	return application
}

func (app *ServerApp) Run() error {
	const op = "server.app.Run"
	logger := app.logger.With(slog.String("op", op))

	logger.Info("Запуск сервера...", slog.String("address", app.server.Addr))

	return app.server.ListenAndServe()
}

func (app *ServerApp) Shutdown() {
	const op = "server.app.Shutdown"
	logger := app.logger.With(slog.String("op", op))
	logger.Info("Завершение сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), app.shutdownTimeout)
	defer cancel()

	if err := app.server.Shutdown(ctx); err != nil {
		warning := fmt.Sprintf("Сервер не успел завершиться за %v.\nЗавершаем принудительно.", app.shutdownTimeout)
		logger.Warn(warning, slog.String("error", err.Error()))

		if err := app.server.Close(); err != nil {
			closeErr := fmt.Sprintf("Ошибка Close: %v", err)
			logger.Error(closeErr, slog.String("error", err.Error()))
		}
	}
}

func registerHandlers() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/users", handlers.UserHandler)
	})

	return r
}
