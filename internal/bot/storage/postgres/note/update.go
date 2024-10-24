package note

import (
	"context"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/elastic"
	"github.com/google/uuid"
)

func (db *NoteRepo) UpdateNote(ctx context.Context, note model.EditNote) error {
	tx, err := db.tx(ctx)
	if err != nil {
		return err
	}

	var id uuid.UUID
	err = tx.QueryRowContext(ctx, `update notes.notes set text = $1, last_edit = $4
	where id = (select id from notes.notes_view where note_number = $3 and user_id = (select id from users.users where tg_id = $2)) returning id`,
		note.Text, note.TgID, note.ViewID, note.Timetag).Scan(&id)
	if err != nil {
		return fmt.Errorf("error while updating note: %+v", err)
	}

	data := elastic.Data{
		Index: elastic.NoteIndex,
		Model: &elastic.Note{
			ID:   id,
			TgID: note.TgID,
			Text: note.Text,
		},
	}

	err = db.elasticClient.Update(ctx, data)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("error while updating note in elastic: %+v", err)
	}

	return tx.Commit()
}
