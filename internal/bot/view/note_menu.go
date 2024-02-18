package view

import tele "gopkg.in/telebot.v3"

var (
	// --------------- заметки --------------

	// inline кнопка для удаления всех заметок
	BtnDeleteAllNotes    = selector.Data("❌Удалить все", "delete_notes")
	BtnSearchNotesByText = selector.Data("🔍Поиск по тексту", "search_notes_by_text")
	BtnSearchNotesByDate = selector.Data("🔍Поиск по дате", "search_notes_by_text")
)

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

// NotesAndMenuBtns возвращает меню с двумя кнопками: Заметки и назад в меню
func NotesAndMenuBtns() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnNotes),
		menu.Row(BtnBackToMenu),
	)

	return menu
}
