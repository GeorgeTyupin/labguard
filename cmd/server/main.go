package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/GeorgeTyupin/labguard/internal/server/app"
	"github.com/go-chi/chi/v5"
)

func main() {
	handler := chi.NewRouter()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	application := app.NewServerApp(handler, logger, ":8080", 10)

	signalCh := make(chan os.Signal, 2)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	errCh := make(chan error, 1)

	go func() {
		if err := application.Run(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case sig := <-signalCh:
		logger.Info("Получен сигнал завершения", slog.String("signal", sig.String()))
	case err := <-errCh:
		logger.Error("Ошибка запуска сервера", slog.String("error", err.Error()))
	}

	application.Shutdown()
}
