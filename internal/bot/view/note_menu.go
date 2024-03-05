package view

import tele "gopkg.in/telebot.v3"

var (
	// --------------- заметки --------------

	// inline кнопка для удаления всех заметок
	BtnDeleteAllNotes = tele.Btn{Text: "❌Удалить все", Unique: "delete_notes"}
	// inline кнопка для поиска по тексту
	BtnSearchNotesByText = tele.Btn{Text: "🔍Поиск по тексту", Unique: "search_notes_by_text"}
	// inline кнопка для поиска по дате
	BtnSearchNotesByDate = tele.Btn{Text: "🔍Поиск по дате", Unique: "search_notes_by_date"}

	// --------------- поиск по дате --------------

	// inline кнопка для возвращения к выбору типа напоминания
	BtnBackToDateType = tele.Btn{Text: "⬅️К выбору", Unique: "search_notes_by_date"}
	// inline кнопка для поиска по одной дате
	BtnSearchByOneDate = tele.Btn{Text: "По одной дате", Unique: "search_by_one_date"}
	// inline кнопка для поиска по двум датам
	BtnSearchByTwoDate = tele.Btn{Text: "По диапазону дат", Unique: "search_by_two_dates"}
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

// SearchByDateBtn возвращает меню с кнопками: поиск по одной дате, поиск по диапазону дат, назад в меню, назад в заметки
func SearchByDateBtn() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnSearchByOneDate, BtnSearchByTwoDate),
		menu.Row(BtnNotes),
		menu.Row(BtnBackToMenu),
	)

	return menu
}
