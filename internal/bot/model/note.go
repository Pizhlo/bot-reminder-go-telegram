package model

import "time"

type Note struct {
	ID      int // id in DB
	UserID  int // id in db
	Text    string
	Created time.Time
	// HasPhoto bool
}
