package keyboards

import (
	tele "gopkg.in/telebot.v4"
)

func NewYesNoMenu() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}

	yesBtn := menu.Text("✅ Да")
	noBtn := menu.Text("❌ Нет")

	menu.Reply(
		menu.Row(yesBtn, noBtn),
	)

	return menu
}
