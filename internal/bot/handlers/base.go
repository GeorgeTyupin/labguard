package handlers

import (
	"log/slog"
	"sync"

	"github.com/GeorgeTyupin/labguard/internal/bot/models"
	tele "gopkg.in/telebot.v4"
)

const (
	msgTypeSuccess = "success"
	msgTypeError   = "error"

	StartEndpoint   = "/start"
	MyEndpoint      = "/my"
	CatalogEndpoint = "/catalog"
)

type BaseHandler struct {
	sendOptions map[string]*tele.SendOptions
	logger      *slog.Logger
	mu          sync.RWMutex
}

func NewBaseHandler(logger *slog.Logger) *BaseHandler {
	handler := &BaseHandler{
		sendOptions: make(map[string]*tele.SendOptions),
		logger:      logger,
	}

	handler.setSendOptions()

	return handler
}

func (h *BaseHandler) setSendOptions() {
	opt := make(map[string]*tele.SendOptions)
	opt[msgTypeSuccess] = &tele.SendOptions{
		ParseMode:   tele.ModeMarkdown,
		ReplyMarkup: &tele.ReplyMarkup{RemoveKeyboard: true},
	}

	opt[msgTypeError] = &tele.SendOptions{
		ReplyMarkup: &tele.ReplyMarkup{RemoveKeyboard: true},
	}

	h.sendOptions = opt
}

type ProductsCache interface {
	Get(int64) ([]*models.Product, error)
	Set(int64, []*models.Product)
	Delete(int64)
	Stop()
}

type BaseProductsHandler struct {
	*BaseHandler
	Cache     ProductsCache
	purchased bool
}

func NewBaseProductsHandler(logger *slog.Logger, cache ProductsCache, purchased bool) *BaseProductsHandler {
	baseHandler := NewBaseHandler(logger)

	productsHandler := &BaseProductsHandler{
		BaseHandler: baseHandler,
		Cache:       cache,
		purchased:   purchased,
	}

	return productsHandler
}
