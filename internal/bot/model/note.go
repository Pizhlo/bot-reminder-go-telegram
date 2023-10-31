package model

import "time"

type Note struct {
	ID      int // id in DB
	UserID  int // id in db
	TgID    int64
	Text    string
	Created time.Time
	// HasPhoto bool
}

// for searching notes by text
type SearchNote struct {
	TgID   int64
	UserID int // id in db!
	Text   string
}
