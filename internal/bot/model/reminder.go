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
)

type Reminder struct {
	ID               int
	TgID             int64
	Name, Date, Time string
	Type             ReminderType
	Created          time.Time
}
