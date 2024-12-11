package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
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
	ID               uuid.UUID // ID in DB
	ViewID           int64     // ID показываемый для пользователя
	TgID             int64
	Name, Date, Time string
	TypeString       sql.NullString
	Type             ReminderType
	Created          time.Time
	Job              Job
	Space            *SharedSpace
	Users            []Participant // пользователи, кому необходимо разослать напоминание
}

type Job struct {
	ID      uuid.UUID
	NextRun time.Time
}
