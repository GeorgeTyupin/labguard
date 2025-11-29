package handlers

import (
	"log/slog"
	"sync"

	tele "gopkg.in/telebot.v4"
)

const (
	msgTypeSuccess = "success"
	msgTypeError   = "error"
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
