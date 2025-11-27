package app

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/GeorgeTyupin/labguard/internal/bot/config"
	"github.com/GeorgeTyupin/labguard/internal/bot/handlers"
	"github.com/GeorgeTyupin/labguard/internal/bot/middleware/loggers"
	"github.com/GeorgeTyupin/labguard/internal/bot/services/api"
	tele "gopkg.in/telebot.v4"
)

const (
	confPath = "configs/bot/bot.yaml"
	envPath  = "configs/bot/bot.env"
)

type BotApp struct {
	AppName string
	Bot     *tele.Bot
	Config  *config.Config
	Logger  *slog.Logger
}

func NewBot(logger *slog.Logger) (*BotApp, error) {
	appName := "Телеграмм бот"

	token, err := config.GetBotToken(envPath)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить токен из переменных окружения в приложении приложение %s, возникла ошибка %w", appName, err)
	}

	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("не удалось сконфигурировать приложение %s, возникла ошибка %w", appName, err)
	}

	bot.Use(loggers.MessageLogger(logger))

	cfg, err := config.Load(confPath)
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить конфиг в приложении приложение %s, возникла ошибка %w", appName, err)
	}

	application := &BotApp{
		Bot:     bot,
		AppName: appName,
		Config:  cfg,
		Logger:  logger,
	}

	application.registerHandlers()

	return application, nil
}

func (app *BotApp) registerHandlers() {
	// TODO Сделать регистрацию всех handlers, после их реализации
	startHandler := handlers.NewStartHandler(api.NewHttpClient(), app.Logger)

	app.Bot.Handle("/start", startHandler.Handle)
	app.Bot.Handle(tele.OnText, startHandler.HandleMessage)
}
