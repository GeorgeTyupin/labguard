package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type ServerApp struct {
	AppName         string
	server          *http.Server
	logger          *slog.Logger
	shutdownTimeout time.Duration
}

func NewServerApp(handler http.Handler, logger *slog.Logger, port string, shutdownTimeout int) *ServerApp {
	appName := "HTTP Server"
	logger = logger.With(slog.String("app", appName))

	server := &http.Server{
		Addr:    port,
		Handler: handler,
	}

	application := &ServerApp{
		AppName:         appName,
		server:          server,
		logger:          logger,
		shutdownTimeout: time.Duration(shutdownTimeout) * time.Second,
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
