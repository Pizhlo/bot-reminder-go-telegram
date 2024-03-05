package view

import tele "gopkg.in/telebot.v3"

// LocationMenu отправляет меню с двумя кнопками: Отправить гео, Отказаться
func LocationMenu() *tele.ReplyMarkup {
	locMenu := &tele.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: true}

	locBtn := locMenu.Location("Отправить геолокацию")
	rejectBtn := locMenu.Text("Отказаться")

	locMenu.Reply(
		locMenu.Row(locBtn),
		locMenu.Row(rejectBtn),
	)

	return locMenu
}
