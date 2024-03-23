package view

import (
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	tele "gopkg.in/telebot.v3"
)

var (
	// --------------- напоминания --------------

	// inline кнопка для удаления сработавшего напоминания
	//BtnDeleteReminder = tele.Btn{Text:"❌Удалить", Unique:"")

	// inline кнопка создания напоминания
	BtnCreateReminder = tele.Btn{Text: "📝Создать напоминание", Unique: "create_reminder"}

	// inline кнопка удаления всех напоминаний
	BtnDeleteAllReminders = tele.Btn{Text: "❌Удалить все", Unique: "delete_reminders"}

	// --------------- типы --------------

	// inline кнопка для возвращения к выбору типа напоминания
	BtnBackToReminderType = tele.Btn{Text: "⬅️К выбору", Unique: "reminder_type"}

	// тип напоминания: несколько раз в день
	BtnSeveralTimesDayReminder = tele.Btn{Text: "Несколько раз в день", Unique: "several_times_day"}

	// тип напоминания: ежедневно
	BtnEveryDayReminder = tele.Btn{Text: "Ежедневно", Unique: "everyday"}

	// тип напоминания: Раз в неделю
	BtnEveryWeekReminder = tele.Btn{Text: "Раз в неделю", Unique: "every_week"}

	// тип напоминания: Раз в несколько дней
	BtnSeveralDaysReminder = tele.Btn{Text: "Раз в несколько дней", Unique: "once_several_days"}

	// тип напоминания: Раз в месяц
	BtnOnceMonthReminder = tele.Btn{Text: "Раз в месяц", Unique: "once_month"}

	// тип напоминания: Раз в год
	BtnOnceYear = tele.Btn{Text: "Раз в год", Unique: "once_year"}

	// тип напоминания: Один раз
	BtnOnce = tele.Btn{Text: "Выбрать дату", Unique: "date"}
	// тип напоминания: Сегодня
	BtnToday = tele.Btn{Text: "Сегодня", Unique: "today"}

	// тип напоминания: Завтра
	BtnTomorrow = tele.Btn{Text: "Завтра", Unique: "tomorrow"}
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
		menu.Row(BtnToday, BtnTomorrow),
		menu.Row(BtnSeveralTimesDayReminder, BtnEveryDayReminder),
		menu.Row(BtnEveryWeekReminder, BtnSeveralDaysReminder),
		menu.Row(BtnOnceMonthReminder, BtnOnceYear),
		menu.Row(BtnOnce),
		menu.Row(BtnBackToMenu),
	)

	return menu
}

// BackToRemindersAndMenu возвращает меню с кнопками: назад к напоминаниями, назад в меню
func BackToRemindersAndMenu() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnReminders),
		menu.Row(BtnBackToMenu),
	)

	return menu
}

// BackToReminderMenuBtns возвращает меню с кнопками: Назад к типам напоминаний, в меню
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

// DeleteReminderBtn возвращает кнопку Удалить.
// Используется для сообщений, когда срабатывает напоминание.
func DeleteReminderBtn(reminder model.Reminder) tele.Btn {
	unique := fmt.Sprintf("%d-%d", reminder.ID, reminder.TgID)

	return tele.Btn{Text: "❌Удалить", Unique: unique, Data: unique}
}

// DeleteReminderMenu возвращает меню с кнопкой Удалить.
// Используется для сообщений, когда срабатывает напоминание.
func DeleteReminderMenu(reminder model.Reminder) *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	unique := fmt.Sprintf("%d.%d", reminder.ID, reminder.TgID)

	btn := menu.Data("❌Удалить", unique)

	menu.Inline(
		menu.Row(btn),
	)

	return menu
}

var (
	// --------------- несколько раз в день --------------

	// тип напоминания: Раз в несколько минут
	BtnMinutesReminder = tele.Btn{Text: "Раз в несколько минут", Unique: "minutes"}

	// тип напоминания: Раз в несколько часов
	BtnHoursReminder = tele.Btn{Text: "Раз в несколько часов", Unique: "hours"}
)

// SeveralTimesBtns возвращает меню с двумя кнопками: раз в несколько минут, раз в несколько часов
func SeveralTimesBtns() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnMinutesReminder, BtnHoursReminder),
		menu.Row(BtnBackToMenu, BtnBackToReminderType),
	)

	return menu
}

var (
	// --------------- раз в неделю --------------

	// тип напоминания: Каждый понедельник
	MondayBtn = tele.Btn{Text: "Понедельник", Unique: "monday"}

	// тип напоминания: Каждый вторник
	TuesdayBtn = tele.Btn{Text: "Вторник", Unique: "tuesday"}

	// тип напоминания: Каждую среду
	WednesdayBtn = tele.Btn{Text: "Среда", Unique: "wednesday"}

	// тип напоминания: Каждый четверг
	ThursdayBtn = tele.Btn{Text: "Четверг", Unique: "thursday"}

	// тип напоминания: Каждую пятницу
	FridayBtn = tele.Btn{Text: "Пятница", Unique: "friday"}

	// тип напоминания: Каждую субботу
	SaturdayBtn = tele.Btn{Text: "Суббота", Unique: "saturday"}

	// тип напоминания: Каждое воскресенье
	SundayBtn = tele.Btn{Text: "Воскресенье", Unique: "sunday"}
)

// WeekMenu возвращает меню днями недели
func WeekMenu() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(MondayBtn, TuesdayBtn),
		menu.Row(WednesdayBtn, ThursdayBtn),
		menu.Row(FridayBtn, SaturdayBtn),
		menu.Row(SundayBtn),
		menu.Row(BtnBackToMenu, BtnBackToReminderType),
	)

	return menu
}
