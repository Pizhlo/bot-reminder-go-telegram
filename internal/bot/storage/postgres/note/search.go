package note

import (
	"context"
	"fmt"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

func (db *NoteRepo) SearchByText(ctx context.Context, searchNote model.SearchByText) ([]model.Note, error) {
	var notes []model.Note

	rows, err := db.db.QueryContext(ctx, `select id, text, created from notes.notes where user_id = (select id from users.users where tg_id = $1) and "text" LIKE '%' || $2 || '%'`, searchNote.TgID, searchNote.Text)
	if err != nil {
		return nil, fmt.Errorf("error while searching note by text: %w", err)
	}

	for rows.Next() {
		note := model.Note{}

		err := rows.Scan(&note.ID, &note.Text, &note.Created)
		if err != nil {
			return nil, fmt.Errorf("errpr while scanning note: %w", err)
		}

		notes = append(notes, note)
	}

	if len(notes) == 0 {
		return nil, api_errors.ErrNotesNotFound
	}

	return notes, nil
}
