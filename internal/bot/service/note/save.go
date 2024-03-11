package note

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

// Save сохраняет заметку пользователя
func (s *NoteService) Save(ctx context.Context, note model.Note) error {
	return s.noteEditor.Save(ctx, note)
}
