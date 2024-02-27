package view

import tele "gopkg.in/telebot.v3"

var (
	// -------------- главное меню --------------

	// inline кнопка для просмотра профиля
	BtnProfile = tele.Btn{Text: "👤Профиль", Unique: "profile"}

	// inline кнопка для просмотра настроек
	BtnSettings = tele.Btn{Text: "⚙️Настройки", Unique: "settings"}

	// inline кнопка просмотра заметок
	BtnNotes = tele.Btn{Text: "📝Заметки", Unique: "notes"}
	// inline кнопка просмотра напоминаний
	BtnReminders = tele.Btn{Text: "⏰Напоминания", Unique: "reminders"}

	// inline кнопка для возвращения в меню
	BtnBackToMenu = tele.Btn{Text: "⬅️Меню", Unique: "menu"}

	// --------------- профиль --------------

	// inline кнопка просмотра информации по подписке
	BtnSubscription = tele.Btn{Text: "🖊Подписка", Unique: "subscription"}
)

// BackToMenuBtn возвращает кнопку возврата в меню
func BackToMenuBtn() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnBackToMenu),
	)

	return menu
}

// MainMenu возвращает главное меню.
// Кнопки: Профиль, Настройки, Заметки, Напоминания
func MainMenu() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnProfile, BtnSettings),
		menu.Row(BtnNotes, BtnReminders),
	)

	return menu
}

// ProfileMenu возвращает меню раздела Профиль
func ProfileMenu() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnSubscription),
		menu.Row(BtnBackToMenu),
	)

	return menu
}
