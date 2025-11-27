package keyboards

import (
	tele "gopkg.in/telebot.v4"
)

const (
	YesText = "✅ Да"
	NoText  = "❌ Нет"
)

func NewYesNoMenu() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}

	yesBtn := menu.Text(YesText)
	noBtn := menu.Text(NoText)

	menu.Reply(
		menu.Row(yesBtn, noBtn),
	)

	return menu
}
