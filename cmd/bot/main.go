package main

import (
	"log/slog"
	"os"
	"os/signal"

	"github.com/GeorgeTyupin/labguard/internal/bot/app"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app, err := app.NewBot(logger)
	if err != nil {
		logger.Error("Не удалось создать приложение бота", slog.String("error", err.Error()))
		return
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	go app.Bot.Start()
	logger.Info("Бот успешно запустился")

	<-ch
	logger.Info("Бот остановлен")
}
