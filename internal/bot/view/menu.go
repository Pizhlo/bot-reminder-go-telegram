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

	// --------------- напоминания --------------

	// inline кнопка для возвращения к выбору типа напоминания
	BtnBackToReminderType = selector.Data("⬅️К напоминаниям", "reminder_type")

	// inline кнопка создания напоминания
	BtnCreateReminder = selector.Data("📝Создать напоминание", "create_reminder")

	// тип напоминания: несколько раз в день
	BtnSeveralTimesDayReminder = selector.Data("Несколько раз в день", "several_times_day")

	// тип напоминания: ежедневно
	BtnEveryDayReminder = selector.Data("Ежедневно", "everyday")

	// тип напоминания: Раз в неделю
	BtnEveryWeekReminder = selector.Data("Раз в неделю", "every_week")

	// тип напоминания: Раз в несколько дней
	BtnSeveralDaysReminder = selector.Data("Раз в несколько дней", "once_several_days")

	// тип напоминания: Раз в месяц
	BtnOnceMonthReminder = selector.Data("Раз в месяц", "once_month")

	// тип напоминания: Раз в год
	BtnOnceYear = selector.Data("Раз в год", "once_year")

	// тип напоминания: Один раз
	BtnOnce = selector.Data("Выбрать дату", "date")
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

// RemindersAndMenuBtns возвращает меню с двумя кнопками: Напоминания и назад в меню
func RemindersAndMenuBtns() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnReminders),
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

// CreateReminderAndBackToMenu возвращает кнопку создания напоминания и назад в меню
func CreateReminderAndBackToMenu() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnCreateReminder),
		menu.Row(BtnBackToMenu),
	)

	return menu
}

// ReminderTypes возвращает меню с типами напоминаний
func ReminderTypes() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnSeveralTimesDayReminder, BtnEveryDayReminder),
		menu.Row(BtnEveryWeekReminder, BtnSeveralDaysReminder),
		menu.Row(BtnOnceMonthReminder, BtnOnceYear),
		menu.Row(BtnOnce),
		menu.Row(BtnBackToMenu),
	)

	return menu
}

// BackToReminderMenuBtns возвращает меню с кнопками: Назад к напоминаниям, в меню
func BackToReminderMenuBtns() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnBackToMenu, BtnBackToReminderType),
	)

	return menu
}
