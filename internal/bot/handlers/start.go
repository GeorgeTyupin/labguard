package handlers

import (
	"fmt"
	"log/slog"

	"github.com/GeorgeTyupin/labguard/internal/bot/keyboards"
	"github.com/GeorgeTyupin/labguard/internal/bot/validators"
	tele "gopkg.in/telebot.v4"
)

type RegisterAPIClient interface {
	CheckUserExists(telegramID int64) (bool, error)
	RegisterUser(telegramID int64, name, group string) (string, error)
}

type RegisterState struct {
	Step  int
	Name  string // ФИО пользователя
	Group string // Группа пользователя
}

type StartHandler struct {
	*BaseHandler
	client     RegisterAPIClient
	userStates map[int64]*RegisterState // telegram_id -> stage
}

func NewStartHandler(apiClient RegisterAPIClient, logger *slog.Logger) *StartHandler {
	baseHandler := NewBaseHandler(logger)

	handler := &StartHandler{
		BaseHandler: baseHandler,
		client:      apiClient,
		userStates:  make(map[int64]*RegisterState),
	}

	return handler
}

func (h *StartHandler) Handle(c tele.Context) error {
	const op = "start.Handle"
	logger := h.logger.With(slog.String("op", op))

	telegramID := c.Sender().ID

	// Начинаем процесс регистрации
	exists, err := h.client.CheckUserExists(telegramID)
	if err != nil {
		logger.Warn("нет метода проверки зарегистрированного пользователя", slog.String("error", err.Error()))
		return c.Send("❌ Ошибка при проверке регистрации")
	}

	if exists {
		return c.Send("Вы уже зарегистрированы! Используйте /my для просмотра токена")
	}

	h.mu.Lock()
	h.userStates[telegramID] = &RegisterState{Step: 1}
	h.mu.Unlock()

	text := `Привет! 👋

Здесь вы можете купить готовые лабораторные работы и курсовые с полным исходным кодом.

После покупки получите:
✅ Рабочий код с комментариями
✅ Доступ к GitHub репозиторию
✅ Персональную лицензию на использование

Для начала давайте зарегистрируемся!


📝 Напишите своё ФИО:`

	return c.Send(text)
}

func (h *StartHandler) HandleMessage(c tele.Context) error {
	const op = "start.HandleMessage"
	logger := h.logger.With(slog.String("op", op))

	telegramID := c.Sender().ID

	h.mu.RLock()
	state, ok := h.userStates[telegramID]
	if !ok {
		return nil // Не в процессе регистрации
	}
	h.mu.RUnlock()

	switch state.Step {
	case 1:
		// Сохраняем ФИО в состояние
		state.Name = c.Text()
		if err := validators.ValidateName(state.Name); err != nil {
			return c.Send(fmt.Sprintf("Неверный формат ФИО : %s.\n\nВы ввели %s.\nВведите ФИО еще раз:", err.Error(), state.Name))
		}

		state.Step = 2
		return c.Send("👥 Теперь введите группу:")

	case 2:
		// Сохраняем группу в состояние
		state.Group = c.Text()
		if err := validators.ValidateGroup(state.Group); err != nil {
			return c.Send(fmt.Sprintf("Неверный формат группу: %s.\n\nВы ввели: %s.\nВведите группу еще раз", err.Error(), state.Group))
		}

		state.Step = 3

		menu := keyboards.NewYesNoMenu()

		return c.Send(fmt.Sprintf("ФИО: %s\nГруппа: %s\n\nВсё верно?", state.Name, state.Group), menu)

	case 3:
		check := c.Text()
		// Проверяем кнопку, которую нажал пользователь
		switch check {
		case keyboards.YesText:
			// Регистрируем пользователя с сохранёнными данными
			token, err := h.client.RegisterUser(telegramID, state.Name, state.Group)
			if err != nil {
				logger.Warn("нет метода регистрации", slog.String("error", err.Error()))
				return c.Send(fmt.Sprintf("❌ Произошла внутренняя ошибка. Попробуйте %s ещё раз позже.", StartEndpoint),
					h.sendOptions[msgTypeError],
				)
			}

			// Удаляем состояние после успешной регистрации
			h.mu.Lock()
			delete(h.userStates, telegramID)
			h.mu.Unlock()

			successMsg := fmt.Sprintf(
				"✅ Регистрация завершена!\n\n"+
					"👤 ФИО: %s\n"+
					"👥 Группа: %s\n"+
					"🔑 Токен: ```%s```\n\n"+
					"📋 Доступные команды:\n"+
					"/catalog — список доступных продуктов\n"+
					"/my — мои покупки и токен\n"+
					"/devices — сброс устройства",
				state.Name, state.Group, token,
			)
			return c.Send(successMsg, h.sendOptions[msgTypeSuccess])

		case keyboards.NoText:
			// Сбрасываем регистрацию
			h.mu.Lock()
			delete(h.userStates, telegramID)
			h.mu.Unlock()

			return c.Send(
				fmt.Sprintf("Регистрация отменена. Введите %s для повторной попытки.", StartEndpoint),
				h.sendOptions[msgTypeSuccess],
			)

		default:
			h.mu.Lock()
			delete(h.userStates, telegramID)
			h.mu.Unlock()

			return c.Send(
				fmt.Sprintf("Сделан неверный выбор. Введите %s для повторной попытки.", StartEndpoint),
				h.sendOptions[msgTypeSuccess],
			)
		}
	}

	return nil
}
