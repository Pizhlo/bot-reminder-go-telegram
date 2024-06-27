package note

import (
	"context"
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

// EditNote обновляет заметку с переданным viewID, заменяя старый текст заметки новым
func (n *NoteService) EditNote(ctx context.Context, userID int64, viewID int64, text string, loc *time.Location) error {
	editNote := model.EditNote{
		TgID:    userID,
		ViewID:  viewID,
		Text:    text,
		Timetag: time.Now().In(loc),
	}

	return n.noteEditor.UpdateNote(ctx, editNote)
}
