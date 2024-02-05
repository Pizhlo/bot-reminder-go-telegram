package note

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

// Save сохраняет заметку пользователя
func (s *NoteService) Save(ctx context.Context, note model.Note) error {
	s.logger.Debugf("Note service: saving user's note. Model: %+v\n", note)

	return s.noteEditor.Save(ctx, note)
}
