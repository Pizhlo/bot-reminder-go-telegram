package note

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

// EditNote обновляет заметку с переданным viewID, заменяя старый текст заметки новым
func (n *NoteService) EditNote(ctx context.Context, userID int64, viewID int64, text string) error {
	editNOte := model.EditNote{
		TgID:   userID,
		ViewID: viewID,
		Text:   text,
	}

	return n.noteEditor.UpdateNote(ctx, editNOte)
}
