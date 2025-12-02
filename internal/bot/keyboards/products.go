package keyboards

import (
	"fmt"

	"github.com/GeorgeTyupin/labguard/internal/bot/models"
	tele "gopkg.in/telebot.v4"
)

const (
	MyUniqueCallback      = "my"
	CatalogUniqueCallback = "catalog"
)

func NewProductsMenu(products []*models.Product, purchased bool) *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	productsBtnList := make([]tele.Row, 0, len(products))

	for i, product := range products {
		if product.Purchased == purchased {
			continue
		}

		var btn tele.Btn
		if purchased {
			btnText := fmt.Sprintf("%s. Куплено ✅", product.Name)
			btn = menu.Data(btnText, MyUniqueCallback, fmt.Sprint(i))
		} else {
			btnText := fmt.Sprintf("%s за %.0f₽", product.Name, product.Price)
			btn = menu.Data(btnText, CatalogUniqueCallback, fmt.Sprint(i))

		}
		productsBtnList = append(productsBtnList, menu.Row(btn))
	}

	menu.Inline(productsBtnList...)

	return menu
}
