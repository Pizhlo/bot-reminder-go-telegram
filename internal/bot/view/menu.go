package view

import tele "gopkg.in/telebot.v3"

var (
	//mainMenu = &tele.ReplyMarkup{}

	// inline кнопка для просмотра профиля
	BtnProfile = selector.Data("👤Профиль", "profile")
	// inline кнопка для просмотра настроек
	BtnSettings = selector.Data("⚙️Настройки", "settings")

	// inline кнопка просмотра заметок
	BtnNotes = selector.Data("📝Заметки", "notes")
	// inline кнопка просмотра напоминаний
	BtnReminders = selector.Data("⏰Напоминания", "reminders")

	// inline кнопка для возвращения в меню
	BtnBackToMenu = selector.Data("⬅️Меню", "menu")
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

// NotesAndMenuBtns возвращает меню с двумя кнопками: Заметки и назад в меню
func NotesAndMenuBtns() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnNotes),
		menu.Row(BtnBackToMenu),
	)

	return menu
}
