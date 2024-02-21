package random

import (
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

// Notes генерирует необходимое количество заметок
func Notes(n int) []model.Note {
	var notes []model.Note

	for i := 0; i < n; i++ {
		note := Note()
		notes = append(notes, note)
	}

	return notes
}

// Note генерирует одну заметку, заполненную рандомными данными
func Note() model.Note {
	return model.Note{
		ID:      Int(0, 10),
		TgID:    1,
		Text:    String(10),
		Created: time.Now(),
	}
}
