package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Note struct {
	ID          uuid.UUID // id in DB
	ViewID      int
	TgID        int64
	Text        string
	Created     time.Time
	LastEditSql sql.NullTime
	LastEdit    time.Time
	// HasPhoto bool
}

// for searching reminders and notes by text
type SearchByText struct {
	TgID   int64
	UserID int // id in db!
	Text   string
}

// для поиска заметок по одной дате
type SearchByOneDate struct {
	TgID   int64
	UserID int
	Date   time.Time
}

// для поиска заметок по двум датам
type SearchByTwoDates struct {
	TgID                  int64
	UserID                int
	FirstDate, SecondDate time.Time
}
