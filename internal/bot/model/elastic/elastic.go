package elastic

import "github.com/google/uuid"

const (
	NoteIndex     = "notes"
	ReminderIndex = "reminders"
)

// Структура для сохранения и поиска в ElasticSearch
type Data struct {
	Model interface{} // Note / Reminder
	Index string
}

// Структура для хранения и поиска по заметкам
type Note struct {
	ID   uuid.UUID // id из базы
	TgID int64
	Text string
}

// Структура для хранения и поиска по напоминаниям
type Reminder struct {
	ID   uuid.UUID // id из базы
	TgID int64
	Text string
}
