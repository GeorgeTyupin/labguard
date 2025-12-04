package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/GeorgeTyupin/labguard/internal/server/app"
	"github.com/GeorgeTyupin/labguard/internal/server/config"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg, err := config.LoadConf()
	if err != nil {
		logger.Error("Ошибка загрузки конфига", slog.String("error", err.Error()))
		os.Exit(1)
	}

	application := app.NewServerApp(logger, cfg)

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
