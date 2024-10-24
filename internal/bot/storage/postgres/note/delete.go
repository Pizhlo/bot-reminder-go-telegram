package note

import (
	"context"
	"errors"
	"fmt"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/elastic"
	"github.com/google/uuid"
)

// DeleteAllByUserID удаляет все заметки пользователя по user ID
func (db *NoteRepo) DeleteAllByUserID(ctx context.Context, userID int64) error {
	tx, err := db.tx(ctx)
	if err != nil {
		return fmt.Errorf("error while creating transaction: %w", err)
	}

	_, err = tx.Exec(`delete from notes.notes where user_id = (select id from users.users where tg_id = $1)`, userID)
	if err != nil {
		return fmt.Errorf("error while deleting all notes by user ID: %w", err)
	}

	data := elastic.Data{
		Index: elastic.NoteIndex,
		Model: elastic.Note{
			TgID: userID,
		},
	}

	err = db.elasticClient.DeleteAllByUserID(ctx, data)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("error while deleting all notes from elastic: %+v", err)
	}

	return tx.Commit()
}

// DeleteNoteByID удаляет одну заметку. Для удаления необходим ID заметки
func (db *NoteRepo) DeleteByID(ctx context.Context, noteID uuid.UUID) error {
	search := elastic.Data{
		Index: elastic.NoteIndex,
		Model: elastic.Note{
			ID: noteID,
		},
	}

	// сначала необходимо узнать id заметки в эластике для удаления
	ids, err := db.elasticClient.SearchByID(ctx, search)
	if err != nil {
		if errors.Is(err, api_errors.ErrRecordsNotFound) {
			return api_errors.ErrNotesNotFound
		}
		return err
	}

	tx, err := db.tx(ctx)
	if err != nil {
		return fmt.Errorf("error while creating transaction: %w", err)
	}

	_, err = tx.Exec(`delete from notes.notes where id = $1`, noteID)
	if err != nil {
		return fmt.Errorf("error while deleting note by user ID: %w", err)
	}

	search.Model = elastic.Note{ElasticID: ids[0]}
	err = db.elasticClient.Delete(ctx, search)
	if err != nil {
		//  откатываем изменения в случае ошибки
		_ = tx.Rollback()
		return fmt.Errorf("error while deleting note from elastic: %+v", err)
	}

	return tx.Commit()
}
