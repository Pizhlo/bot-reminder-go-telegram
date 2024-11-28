package model

import (
	"database/sql"
	"strconv"
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
	SpaceID     int
	// HasPhoto bool
}

const createdFieldFormat = "02.01.2006 15:04:05"

// Fields реализует интерфейс baseView.record. Возвращает значения всех полей
func (n Note) Fields() map[string]string {
	res := make(map[string]string)

	res["ID"] = n.ID.String()
	res["ViewID"] = strconv.Itoa(n.ViewID)
	res["TgID"] = strconv.Itoa(int(n.TgID))
	res["Text"] = n.Text
	res["Created"] = n.Created.Format(createdFieldFormat)

	if n.LastEditSql.Valid {
		res["Edited"] = n.LastEditSql.Time.GoString()

		if n.SpaceID != 0 {
			res["SpaceID"] = strconv.Itoa(n.SpaceID)
		}
	} else {
		if n.SpaceID != 0 {
			res["SpaceID"] = strconv.Itoa(n.SpaceID)
		}
	}

	return res
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

// для редактирования заметки
type EditNote struct {
	TgID    int64
	ViewID  int64
	Text    string
	Timetag time.Time
}
