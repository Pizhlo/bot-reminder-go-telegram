package view

import (
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	tele "gopkg.in/telebot.v3"
)

var (
	// --------------- напоминания --------------

	// inline кнопка для удаления сработавшего напоминания
	BtnDeleteReminder = selector.Data("❌Удалить", "")

	// inline кнопка создания напоминания
	BtnCreateReminder = selector.Data("📝Создать напоминание", "create_reminder")

	// inline кнопка удаления всех напоминаний
	BtnDeleteAllReminders = selector.Data("❌Удалить все", "delete_reminders")

	// --------------- типы --------------

	// inline кнопка для возвращения к выбору типа напоминания
	BtnBackToReminderType = selector.Data("⬅️К напоминаниям", "reminder_type")

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

// CreateReminderAndBackToMenu возвращает кнопку создания напоминания, удалить все и назад в меню
func CreateReminderDeleteAndBackToMenu() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnCreateReminder),
		menu.Row(BtnDeleteAllReminders),
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

// RemindersAndMenuBtns возвращает меню с двумя кнопками: Напоминания и назад в меню
func RemindersAndMenuBtns() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnReminders),
		menu.Row(BtnBackToMenu),
	)

	return menu
}

// DeleteReminderBtn возвращает меню с кнопкой Удалить.
// Используется для сообщений, когда срабатывает напоминание.
func DeleteReminderBtn(reminder model.Reminder) *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnDeleteReminder),
	)

	return menu
}
