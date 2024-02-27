package view

import tele "gopkg.in/telebot.v3"

var (
	// --------------- заметки --------------

	// inline кнопка для удаления всех заметок
	BtnDeleteAllNotes    = tele.Btn{Text: "❌Удалить все", Unique: "delete_notes"}
	BtnSearchNotesByText = tele.Btn{Text: "🔍Поиск по тексту", Unique: "search_notes_by_text"}
	BtnSearchNotesByDate = tele.Btn{Text: "🔍Поиск по дате", Unique: "search_notes_by_date"}
)

// DeleteAllNotesAndBackToMenu возвращает меню с кнопками:
// удалить все заметки, назад в меню, поиск по тексту, поиск по дате, назад в меню
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

// BackToMenuAndNotesBtn возвращает меню с кнопками: назад в меню, назад в заметки
func BackToMenuAndNotesBtn() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnNotes),
		menu.Row(BtnBackToMenu),
	)

	return menu
}
