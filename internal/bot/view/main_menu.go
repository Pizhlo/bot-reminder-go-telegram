package view

import tele "gopkg.in/telebot.v3"

var (
	// -------------- главное меню --------------

	// inline кнопка для просмотра часового пояса
	BtnTimezone = tele.Btn{Text: "🌍Часовой пояс", Unique: "timezone"}

	// inline кнопка просмотра заметок
	BtnNotes = tele.Btn{Text: "📝Заметки", Unique: "notes"}
	// inline кнопка просмотра напоминаний
	BtnReminders = tele.Btn{Text: "⏰Напоминания", Unique: "reminders"}

	// inline кнопка для доступа к совместным пространствам
	BtnSharedAccess = tele.Btn{Text: "🫂Совместный доступ", Unique: "shared_access"}

	BtnBugReport = tele.Btn{Text: "⚠️Сообщить о баге", Unique: "bag_report"}

	// inline кнопка для возвращения в меню
	BtnBackToMenu = tele.Btn{Text: "⬅️Меню", Unique: "menu"}

	// inline кнопка чтобы создать еще одно напоминание
	BtnCreateOneElseReminder = tele.Btn{Text: "🔄Создать еще", Unique: "create_reminder"}

	// --------------- часовой пояс --------------

	// inline кнопка изменения часового пояса
	BtnEditTimezone = tele.Btn{Text: "✏️Изменить", Unique: "edit_timezone"}
)

// BackToMenuBtn возвращает кнопку возврата в меню
func BackToMenuBtn() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnBackToMenu),
	)

	return menu
}

// BackToMenuAndCreateOneElse возвращает меню с кнопками: создать еще, назад в меню
func BackToMenuAndCreateOneElse() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnReminders),
		menu.Row(BtnCreateOneElseReminder),
		menu.Row(BtnBackToMenu),
	)

	return menu
}

// MainMenu возвращает главное меню.
// Кнопки: Профиль, Настройки, Заметки, Напоминания, Сообщить о баге
func MainMenu() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnNotes, BtnReminders),
		menu.Row(BtnSharedAccess),
		menu.Row(BtnTimezone),
		menu.Row(BtnBugReport),
	)

	return menu
}

// TimezoneMenu возвращает меню раздела Часовой пояс
func TimezoneMenu() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnEditTimezone),
		menu.Row(BtnBackToMenu),
	)

	return menu
}
