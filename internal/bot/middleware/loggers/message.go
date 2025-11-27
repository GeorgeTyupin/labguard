package loggers

import (
	"log/slog"
	"time"

	tele "gopkg.in/telebot.v4"
)

func MessageLogger(logger *slog.Logger) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			logger.Debug("Обработка сообщения")
			start := time.Now()

			err := next(c)
			duration := time.Since(start)
			if err != nil {
				logger.Error("Ошибка обработки",
					slog.Int64("telegram_id", c.Sender().ID),
					slog.String("error", err.Error()),
					slog.Duration("duration", duration),
				)
			} else {
				logger.Info("Сообщение обработано",
					slog.Int64("telegram_id", c.Sender().ID),
					slog.String("user_message", c.Text()),
					slog.Duration("duration", duration),
				)
			}

			return err
		}
	}
}
