package keyboards

import (
	"fmt"

	"github.com/GeorgeTyupin/labguard/internal/bot/models"
	tele "gopkg.in/telebot.v4"
)

func NewProductsMenu(products []*models.Product) *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	productsBtnList := make([]tele.Row, 0, len(products))

	for i, product := range products {
		if !product.Purchased {
			btnText := fmt.Sprintf("%s за %.0f₽", product.Name, product.Price)
			btn := menu.Data(btnText, "product", fmt.Sprint(i))

			productsBtnList = append(productsBtnList, menu.Row(btn))
		}
	}

	menu.Inline(productsBtnList...)

	return menu
}
