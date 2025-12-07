package app

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/GeorgeTyupin/labguard/internal/bot/config"
	"github.com/GeorgeTyupin/labguard/internal/bot/handlers"
	"github.com/GeorgeTyupin/labguard/internal/bot/keyboards"
	"github.com/GeorgeTyupin/labguard/internal/bot/middleware/loggers"
	"github.com/GeorgeTyupin/labguard/internal/bot/models"
	"github.com/GeorgeTyupin/labguard/internal/bot/services/api"
	"github.com/GeorgeTyupin/labguard/pkg/cache"
	tele "gopkg.in/telebot.v4"
)

type BotApp struct {
	AppName string
	Bot     *tele.Bot
	Config  *config.Config
	Logger  *slog.Logger
	cleanup []func()
}

func NewBot(logger *slog.Logger, cfg *config.Config) (*BotApp, error) {
	appName := "Телеграмм бот"
	logger = logger.With(slog.String("app", appName))

	pref := tele.Settings{
		Token:  cfg.BotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("не удалось сконфигурировать приложение %s, возникла ошибка %w", appName, err)
	}

	bot.Use(loggers.MessageLogger(logger))

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
	apiClient := api.NewHttpClient()

	// Приложение для регистрации
	startHandler := handlers.NewStartHandler(apiClient, app.Logger)
	app.Bot.Handle(handlers.StartEndpoint, startHandler.Handle)
	app.Bot.Handle(tele.OnText, startHandler.HandleMessage)

	productCache := cache.NewCacheWithTTL[int64, []*models.Product](time.Duration(10 * time.Minute)) // Кеш неоплаченных продуктов

	// Приложение для получения списка доступных продуктов
	catalogHandler := handlers.NewCatalogHandler(apiClient, app.Logger, productCache)
	app.cleanup = append(app.cleanup, func() {
		catalogHandler.Cache.Stop()
	})
	app.Bot.Handle(handlers.CatalogEndpoint, catalogHandler.Handle)
	productBtn := &tele.Btn{Unique: keyboards.CatalogUniqueCallback}
	app.Bot.Handle(productBtn, catalogHandler.HandleCatalogCallbacks)
	buyBtn := &tele.Btn{Unique: keyboards.BuyUniqueCallback}
	app.Bot.Handle(buyBtn, catalogHandler.HandleBuyCallbacks)

	// Приложение для получения списка купленных продуктов
	myProductsCache := cache.NewCacheWithTTL[int64, []*models.Product](time.Duration(10 * time.Minute))
	myHandler := handlers.NewMyHandler(apiClient, app.Logger, myProductsCache)
	app.cleanup = append(app.cleanup, func() {
		myHandler.Cache.Stop()
	})
	app.Bot.Handle(handlers.MyEndpoint, myHandler.Handle)
	myProductBtn := &tele.Btn{Unique: keyboards.MyUniqueCallback}
	app.Bot.Handle(myProductBtn, myHandler.HandleCallbacks)
}

func (app *BotApp) Shutdown() {
	for _, cleanFunc := range app.cleanup {
		cleanFunc()
	}

	app.Bot.Stop()
}
