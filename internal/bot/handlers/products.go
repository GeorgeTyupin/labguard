package handlers

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/GeorgeTyupin/labguard/internal/bot/keyboards"
	"github.com/GeorgeTyupin/labguard/internal/bot/models"
	tele "gopkg.in/telebot.v4"
)

type ProductsAPIClient interface {
	CheckUserExists(telegramID int64) (bool, error)
	GetProducts(telegramID int64) ([]*models.Product, error)
}

type ICache interface {
	Get(int64) ([]*models.Product, error)
	Set(int64, []*models.Product)
	Delete(int64)
	Stop()
}

type ProductsHandler struct {
	base         *BaseHandler
	client       ProductsAPIClient
	UserProducts ICache
}

func NewProductsHandler(apiClient ProductsAPIClient, logger *slog.Logger, cache ICache) *ProductsHandler {
	baseHandler := NewBaseHandler(logger)

	handler := &ProductsHandler{
		base:         baseHandler,
		client:       apiClient,
		UserProducts: cache,
	}

	return handler
}

func (h *ProductsHandler) Handle(c tele.Context) error {
	const op = "products.Handle"
	logger := h.base.logger.With(slog.String("op", op))

	telegramID := c.Sender().ID

	// Проверяем регистрацию пользователя
	_, err := h.client.CheckUserExists(telegramID)
	if err != nil {
		logger.Warn("нет метода проверки зарегистрированного пользователя", slog.String("error", err.Error()))
		return c.Send("❌ Ошибка при проверке регистрации")
	}

	// TODO: расскоментировать после реализации
	// if !exists {
	// 	return c.Send("Вы еще не зарегистрированы! Используйте /start для регистрации")
	// }

	products, err := h.client.GetProducts(telegramID)
	if err != nil {
		logger.Warn("нет метода получения списка продуктов", slog.String("error", err.Error()))
		return c.Send("❌ Ошибка при попытке получить список продуктов")
	}

	h.UserProducts.Set(telegramID, products)

	productsMenu := keyboards.NewProductsMenu(products)

	return c.Send("Список продуктов:\n", productsMenu)
}

func (h *ProductsHandler) HandleCallbacks(c tele.Context) error {
	const op = "products.HandleCallbacks"
	logger := h.base.logger.With(slog.String("op", op))

	defer c.Respond()

	// Проверяем, что это callback для продуктов
	if c.Callback().Unique != "product" {
		logger.Warn("Unique не совпадает с product", slog.String("unique", c.Callback().Unique))
		return nil
	}

	// Извлекаем индекс продукта
	productIdx, err := strconv.Atoi(c.Callback().Data)
	if err != nil {
		logger.Error(
			"Не удалось конвертировать индекс продукта из строки в число",
			slog.String("data", c.Callback().Data),
		)
		return c.Send("❌ Возникла внутренняя ошибка. Попробуйте ввести /products еще раз")
	}

	telegramID := c.Sender().ID

	products, err := h.UserProducts.Get(telegramID)
	if err != nil || productIdx < 0 || productIdx >= len(products) {
		logger.Info("Ошибка получения элемента из кеша", slog.String("error", err.Error()))
		return c.Send("❌ Продукт не найден. Попробуйте вызвать /products еще раз")
	}
	product := products[productIdx]

	logger.Info("Успешно получили продукт через callback", slog.Any("product", product))

	message := fmt.Sprintf(
		"*📦 %s*\n\n"+
			"_%s_\n\n"+
			"💰 *Цена:* %.0f₽\n",
		product.Name,
		product.Description,
		product.Price,
	)

	return c.Send(message, h.base.sendOptions[msgTypeSuccess])
}
