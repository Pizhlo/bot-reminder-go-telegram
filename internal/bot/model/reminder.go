package model

import (
	"time"
)

type ReminderType string

const (
	// несколько раз в день
	SeveralTimesDayType ReminderType = "several_times_day"
	// ежедневно
	EverydayType ReminderType = "everyday"
	// раз в неделю
	EveryWeekType ReminderType = "everyweek"
	// раз в несколько дней
	SeveralDaysType ReminderType = "several_days"
	// раз в месяц
	OnceMonthType ReminderType = "once_month"
	// раз в год
	OnceYearType ReminderType = "once_year"
	// один раз - в указанную дату
	DateType ReminderType = "date"
	// прислать напоминание сегодня
	Today ReminderType = "today"
	// прислать завтра
	Tomorrow ReminderType = "tomorrow"
)

type Reminder struct {
	ID               int64
	TgID             int64
	Name, Date, Time string
	Type             ReminderType
	Created          time.Time
}
