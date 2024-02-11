package view

import tele "gopkg.in/telebot.v3"

var (
	// -------------- главное меню --------------

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

	// --------------- профиль --------------

	// inline кнопка просмотра информации по подписке
	BtnSubscription = selector.Data("🖊Подписка", "subscription")

	// --------------- заметки --------------

	// inline кнопка для удаления всех заметок
	BtnDeleteAllNotes    = selector.Data("❌Удалить все", "delete_notes")
	BtnSearchNotesByText = selector.Data("🔍Поиск по тексту", "search_notes_by_text")
	BtnSearchNotesByDate = selector.Data("🔍Поиск по дате", "search_notes_by_text")
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

// DeleteAllNotesAndBackToMenu возвращает меню с двумя кнопками: удалить все заметки и назад в меню
func DeleteAllNotesAndBackToMenu() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnDeleteAllNotes),
		menu.Row(BtnSearchNotesByText, BtnSearchNotesByDate),
		menu.Row(BtnBackToMenu),
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
