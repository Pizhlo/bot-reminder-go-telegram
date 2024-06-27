package note

import (
	"context"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

func (db *NoteRepo) UpdateNote(ctx context.Context, note model.EditNote) error {
	tx, err := db.tx(ctx)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `update notes.notes set text = $1, last_edit = $4
	where id = (select id from notes.notes_view where note_number = $3 and user_id = (select id from users.users where tg_id = $2))`,
		note.Text, note.TgID, note.ViewID, note.Timetag)
	if err != nil {
		return fmt.Errorf("error while updating note: %+v", err)
	}

	return tx.Commit()
}
