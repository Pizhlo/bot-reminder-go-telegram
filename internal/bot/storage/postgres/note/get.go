package note

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

// GetAllByUserID достает из базы все заметки пользователя по ID
func (db *NoteRepo) GetAllByUserID(ctx context.Context, userID int64) ([]model.Note, error) {
	notes := make([]model.Note, 0)

	rows, err := db.db.QueryContext(ctx, `select note_number, text, created from notes.notes 
	join notes.notes_view on notes.notes_view.id = notes.notes.id
	where notes.notes.user_id = (select id from users.users where tg_id = $1) 
	order by created ASC;`, userID)
	if err != nil {
		return nil, fmt.Errorf("error while getting all notes from DB by user ID %d: %w", userID, err)
	}
	defer rows.Close()

	for rows.Next() {
		note := model.Note{}

		err := rows.Scan(&note.ID, &note.Text, &note.Created)
		if err != nil {
			return nil, fmt.Errorf("error while scanning note (get all by user id): %w", err)
		}

		notes = append(notes, note)
	}

	if len(notes) == 0 {
		return nil, api_errors.ErrNotesNotFound
	}

	return notes, nil
}

// GetByID возвращает заметку с переданным ID. Если такой заметки нет, возвращает ErrNotesNotFound
func (db *NoteRepo) GetByID(ctx context.Context, userID int64, noteID int) (*model.Note, error) {
	note := model.Note{}

	row := db.db.QueryRowContext(ctx, `select note_number, text, created from notes.notes 
	join notes.notes_view on notes.notes_view.id = notes.notes.id
	where notes.notes.user_id = (select id from users.users where tg_id = $1)
	and notes.notes_view.note_number = $2
	order by created ASC;`, userID, noteID)

	err := row.Scan(&note.ID, &note.Text, &note.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, api_errors.ErrNotesNotFound
		}
		return nil, fmt.Errorf("error while scanning note (get by id): %w", err)
	}

	return &note, nil
}
