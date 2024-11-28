package model

import (
	"time"
)

// совместный доступ для нескольких пользователей к заметкам и напоминаниям
type SharedSpace struct {
	ID           int
	ViewID       int // номер для отображения
	Name         string
	Created      time.Time
	Participants []User
	Creator      User
	Notes        []Note
	Reminders    []Reminder
}

// структура, связывающая пространство и пользователя (многие-ко-многим)
type LinkSpaceUser struct {
	ID      int
	SpaceID int
	UserID  int
}
