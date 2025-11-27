package handlers

import (
	"fmt"
	"log/slog"

	"github.com/GeorgeTyupin/labguard/internal/bot/keyboards"
	tele "gopkg.in/telebot.v4"
)

type RegisterAPIClient interface {
	// TODO Добавить методы после того как реализуется сам клиент
	CheckUserExists(telegramID int64) (bool, error)
	RegisterUser(telegramID int64, name, group string) (string, error)
}

type RegisterState struct {
	Step  int
	Name  string // ФИО пользователя
	Group string // Группа пользователя
}

type StartHandler struct {
	client     RegisterAPIClient
	userStates map[int64]*RegisterState // telegram_id -> stage
	logger     *slog.Logger
}

func NewStartHandler(apiClient RegisterAPIClient, logger *slog.Logger) *StartHandler {
	return &StartHandler{
		client:     apiClient,
		userStates: make(map[int64]*RegisterState),
		logger:     logger,
	}
}

func (sh *StartHandler) Handle(c tele.Context) error {
	const op = "start.Handle"
	logger := sh.logger.With(slog.String("op", op))

	telegramID := c.Sender().ID

	// Начинаем процесс регистрации
	exists, err := sh.client.CheckUserExists(telegramID)
	if err != nil {
		logger.Warn("нет метода проверки зарегистрированного пользователя", slog.String("error", err.Error()))
		return c.Send("❌ Ошибка при проверке регистрации")
	}

	if exists {
		return c.Send("Ты уже зарегистрирован/зарегистрирована! Используй /my для просмотра токена")
	}

	sh.userStates[telegramID] = &RegisterState{Step: 1}
	text := `Привет! 👋

Здесь ты можешь купить готовые лабораторные работы и курсовые с полным исходным кодом.

После покупки получишь:
✅ Рабочий код с комментариями
✅ Доступ к GitHub репозиторию
✅ Персональную лицензию на использование

Для начала давай зарегистрируемся!


📝 Напиши своё ФИО:`

	return c.Send(text)
}

func (sh *StartHandler) HandleMessage(c tele.Context) error {
	const op = "start.HandleMessage"
	logger := sh.logger.With(slog.String("op", op))

	telegramID := c.Sender().ID
	state, ok := sh.userStates[telegramID]
	if !ok {
		return nil // Не в процессе регистрации
	}

	switch state.Step {
	case 1:
		// Сохраняем ФИО в состояние
		state.Name = c.Text()
		state.Step = 2
		return c.Send("👥 Теперь введи группу:")

	case 2:
		// Сохраняем группу в состояние
		state.Group = c.Text()
		state.Step = 3

		menu := keyboards.NewYesNoMenu()

		return c.Send(fmt.Sprintf("ФИО: %s\nГруппа: %s\n\nВсё верно?", state.Name, state.Group), menu)

	case 3:
		check := c.Text()
		// Проверяем кнопку, которую нажал пользователь
		switch check {
		case keyboards.YesText:
			// Регистрируем пользователя с сохранёнными данными
			token, err := sh.client.RegisterUser(telegramID, state.Name, state.Group)
			if err != nil {
				logger.Warn("нет метода регистрации", slog.String("error", err.Error()))
				return c.Send("❌ Произошла внутренняя ошибка. Попробуй /start ещё раз позже.")
			}

			// Удаляем состояние после успешной регистрации
			delete(sh.userStates, telegramID)

			return c.Send(fmt.Sprintf("✅ Регистрация завершена!\n\n👤 ФИО: %s\n👥 Группа: %s\n🔑 Токен: ```%s```", state.Name, state.Group, token), &tele.SendOptions{ParseMode: tele.ModeMarkdown})
		case keyboards.NoText:
			// Сбрасываем регистрацию
			delete(sh.userStates, telegramID)
			return c.Send("Регистрация отменена. Введи /start для повторной попытки")
		default:
			delete(sh.userStates, telegramID)
			return c.Send("Сделан неверный выбор. Введи /start для повторной попытки")
		}
	}

	return nil
}
