package model

import "time"

type Reminder struct {
	ID                     int
	TgID                   int64
	Text, Type, Date, Time string
	Created                time.Time
}
