package model

type Note struct {
	ID     int // id in DB
	UserID int // id in db
	Text   string
	// HasPhoto bool
}
