package handlers

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/GeorgeTyupin/labguard/internal/bot/keyboards"
	"github.com/GeorgeTyupin/labguard/internal/bot/models"
	tele "gopkg.in/telebot.v4"
)

type MyAPIClient interface {
	CheckUserExists(telegramID int64) (bool, error)
	GetProducts(telegramID int64) ([]*models.Product, error)
}

type MyHandler struct {
	*BaseProductsHandler
	client MyAPIClient
}

func NewMyHandler(apiClient MyAPIClient, logger *slog.Logger, cache ProductsCache) *MyHandler {
	baseHandler := NewBaseProductsHandler(logger, cache, true)

	handler := &MyHandler{
		BaseProductsHandler: baseHandler,
		client:              apiClient,
	}

	return handler
}

func (h *MyHandler) Handle(c tele.Context) error {
	const op = "my.Handle"
	logger := h.logger.With(slog.String("op", op))

	telegramID := c.Sender().ID

	// Проверяем регистрацию пользователя
	_, err := h.client.CheckUserExists(telegramID)
	if err != nil {
		logger.Warn("нет метода проверки зарегистрированного пользователя", slog.String("error", err.Error()))
		return c.Send("❌ Ошибка при проверке регистрации")
	}

	// TODO: расскоментировать после реализации
	// if !exists {
	// 	return c.Send(fmt.Sprintf("Вы еще не зарегистрированы! Используйте %s для регистрации", StartEndpoint))
	// }

	var products []*models.Product

	products, err = h.Cache.Get(telegramID) // Сначала пробуем получить из кеша
	if err != nil {
		products, err = h.client.GetProducts(telegramID)
		if err != nil {
			logger.Warn("нет метода получения списка продуктов", slog.String("error", err.Error()))
			return c.Send("❌ Ошибка при попытке получить список продуктов")
		}
	}

	h.Cache.Set(telegramID, products)

	productsMenu := keyboards.NewProductsMenu(products, h.purchased)

	return c.Send("Список ваших продуктов:\n", productsMenu)
}

func (h *MyHandler) HandleCallbacks(c tele.Context) error {
	const op = "my.HandleCallbacks"
	logger := h.logger.With(slog.String("op", op))

	defer c.Respond()

	// Проверяем, что это callback для продуктов
	if c.Callback().Unique != keyboards.MyUniqueCallback {
		logger.Warn("Unique не совпадает с my", slog.String("unique", c.Callback().Unique))
		return nil
	}

	// Извлекаем индекс продукта
	productIdx, err := strconv.Atoi(c.Callback().Data)
	if err != nil {
		logger.Error(
			"Не удалось конвертировать индекс продукта из строки в число",
			slog.String("data", c.Callback().Data),
		)
		return c.Send(fmt.Sprintf("❌ Возникла внутренняя ошибка. Попробуйте ввести %s еще раз", MyEndpoint))
	}

	telegramID := c.Sender().ID

	products, err := h.Cache.Get(telegramID)
	if err != nil || productIdx < 0 || productIdx >= len(products) {
		logger.Info("Ошибка получения элемента из кеша", slog.String("error", err.Error()))
		return c.Send(fmt.Sprintf("❌ Продукт не найден. Попробуйте вызвать %s еще раз", MyEndpoint))
	}
	product := products[productIdx]

	logger.Info("Успешно получили продукт через callback", slog.Any("product", product))

	message := fmt.Sprintf(
		"*📦 %s*\n\n"+
			"_%s_\n\n"+
			"🔗 [GitHub](%s)\n",
		product.Name,
		product.Description,
		product.Link,
	)

	return c.Send(message, h.sendOptions[msgTypeSuccess])
}
