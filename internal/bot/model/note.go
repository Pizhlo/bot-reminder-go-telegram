package model

import "time"

type Note struct {
	ID      int // id in DB
	TgID    int64
	Text    string
	Created time.Time
	// HasPhoto bool
}

// for searching reminders and notes by text
type SearchByText struct {
	TgID   int64
	UserID int // id in db!
	Text   string
}
type SearchByOneDate struct {
	TgID   int64
	UserID int
	Date   time.Time
}
