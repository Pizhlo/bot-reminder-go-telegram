package elastic

import "github.com/google/uuid"

// Структура для хранения и поиска по напоминаниям
type Reminder struct {
	ID   uuid.UUID // id из базы
	TgID int64
	Text string
}
