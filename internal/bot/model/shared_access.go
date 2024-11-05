package model

import "time"

// совместный доступ для нескольких пользователей к заметкам и напоминаниям
type SharedSpace struct {
	ID           int
	Name         string
	Created      time.Time
	Participants []User
}

// структура, связывающая пространство и пользователя (многие-ко-многим)
type LinkSpaceUser struct {
	ID      int
	SpaceID int
	UserID  int
}
