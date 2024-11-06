package model

import (
	"time"

	user "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
)

// совместный доступ для нескольких пользователей к заметкам и напоминаниям
type SharedSpace struct {
	ID           int
	ViewID       int // номер для отображения
	Name         string
	Created      time.Time
	Participants []user.User
	Creator      user.User
	Notes        []Note
	Reminders    []Reminder
}

// структура, связывающая пространство и пользователя (многие-ко-многим)
type LinkSpaceUser struct {
	ID      int
	SpaceID int
	UserID  int
}
