package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/GeorgeTyupin/labguard/internal/server/app"
	"github.com/GeorgeTyupin/labguard/internal/server/config"
	"github.com/GeorgeTyupin/labguard/internal/server/repository/postgres"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg := config.MustLoad(logger)

	db := postgres.MustDBPoolInit(logger, cfg.PostgresConfig)
	defer db.Close()

	application := app.NewServerApp(logger, cfg, db)

	signalCh := make(chan os.Signal, 2)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	errCh := make(chan error, 2)

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
