package note

import (
	"context"
	"database/sql"
	"fmt"
)

// DeleteAllByUserID удаляет все заметки пользователя по user ID
func (db *NoteRepo) DeleteAllByUserID(ctx context.Context, userID int64) error {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("error while creating transaction: %w", err)
	}

	_, err = tx.Exec(`delete from notes.notes where user_id = (select id from users.users where tg_id = $1)`, userID)
	if err != nil {
		return fmt.Errorf("error while deleting all notes by user ID: %w", err)
	}

	return tx.Commit()
}

// DeleteNoteByID удаляет одну заметку. Для удаления необходим ID заметки и пользователя
func (db *NoteRepo) DeleteNoteByViewID(ctx context.Context, userID int64, noteID int) error {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("error while creating transaction: %w", err)
	}

	_, err = tx.Exec(`delete from notes.notes where user_id = (select id from users.users where tg_id = $1) and id = (select id from notes.notes_view where notes.notes_view.note_number = $2 and user_id = (select id from users.users where tg_id = $1))`, userID, noteID)
	if err != nil {
		return fmt.Errorf("error while deleting note by user ID: %w", err)
	}

	return tx.Commit()
}
