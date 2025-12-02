package keyboards

import (
	"fmt"

	tele "gopkg.in/telebot.v4"
)

func NewBuyMenu(id int64) *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	btnText := "Купить 🛒"
	btn := menu.Data(btnText, BuyUniqueCallback, fmt.Sprint(id))
	menu.Inline(menu.Row(btn))

	return menu
}
