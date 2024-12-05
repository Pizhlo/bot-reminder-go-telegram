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
	Participants []Participant
	Creator      Participant
	Notes        []Note
	Reminders    []Reminder
}

type Participant struct {
	User
	State string // состояние (добавлен, отклонил, ожидание)
}

// состояния участников пространства
const (
	// состояние, при котором участник еще не ответил на приглашение
	PendingState = "pending"
	// состояние, когда участник успешно добавлен в пространство
	AddedState = "added"
	// состояние, когда участник отклонил приглашение
	RejectedState = "rejected"
)

// структура, связывающая пространство и пользователя (многие-ко-многим)
type LinkSpaceUser struct {
	ID      int
	SpaceID int
	UserID  int
}
