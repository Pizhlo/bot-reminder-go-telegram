package note

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

// Save сохраняет заметку в базе данных. Для сохранения требуется: ID пользователя, содержимое заметки, дата создания
func (db *NoteRepo) Save(ctx context.Context, note model.Note) error {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("error while creating transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, `insert into notes.notes (user_id, text, created) values((select id from users.users where tg_id=$1), $2, $3)`, note.TgID, note.Text, note.Created)
	if err != nil {
		return fmt.Errorf("error inserting note: %w", err)
	}

	return tx.Commit()
}
