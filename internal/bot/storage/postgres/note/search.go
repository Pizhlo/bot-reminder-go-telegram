package note

import (
	"context"
	"fmt"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/elastic"
	"github.com/jmoiron/sqlx"
)

func (db *NoteRepo) SearchByText(ctx context.Context, searchNote model.SearchByText) ([]model.Note, error) {
	search := elastic.Data{
		Index: elastic.NoteIndex,
		Model: elastic.Note{
			TgID: searchNote.TgID,
			Text: searchNote.Text,
		},
	}

	ids, err := db.elasticClient.Search(ctx, search)
	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return nil, api_errors.ErrNotesNotFound
	}

	var notes []model.Note

	q, args, err := sqlx.In(`select note_number, text, created, last_edit from notes.notes
	join notes.notes_view on notes.notes_view.id = notes.notes.id
	where notes.notes.id IN(?);`, ids)
	if err != nil {
		return nil, fmt.Errorf("error while creating query while searching notes: %+v", err)
	}

	q = sqlx.Rebind(sqlx.DOLLAR, q)
	rows, err := db.db.Query(q, args...)
	if err != nil {
		return nil, fmt.Errorf("error while searching notes by text: %w", err)
	}

	for rows.Next() {
		note := model.Note{}

		err := rows.Scan(&note.ViewID, &note.Text, &note.Created, &note.LastEditSql)
		if err != nil {
			return nil, fmt.Errorf("error while scanning note (search by text): %w", err)
		}

		notes = append(notes, note)
	}

	if len(notes) == 0 {
		return nil, api_errors.ErrNotesNotFound
	}

	return notes, nil
}

func (db *NoteRepo) SearchByOneDate(ctx context.Context, searchNote model.SearchByOneDate) ([]model.Note, error) {
	var notes []model.Note

	rows, err := db.db.QueryContext(ctx, `select note_number, text, created from notes.notes 
	join notes.notes_view on notes.notes_view.id = notes.notes.id
	where notes.notes.user_id = (select id from users.users where tg_id = $1)
	and created >= $2::date AND created < ($2::date + '1 day'::interval);`, searchNote.TgID, searchNote.Date)
	if err != nil {
		return nil, fmt.Errorf("error while searching notes by one date: %w", err)
	}

	for rows.Next() {
		note := model.Note{}

		err := rows.Scan(&note.ViewID, &note.Text, &note.Created)
		if err != nil {
			return nil, fmt.Errorf("error while scanning note (search by one date): %w", err)
		}

		notes = append(notes, note)
	}

	if len(notes) == 0 {
		return nil, api_errors.ErrNotesNotFound
	}

	return notes, nil
}

func (db *NoteRepo) SearchByTwoDates(ctx context.Context, searchNote *model.SearchByTwoDates) ([]model.Note, error) {
	var notes []model.Note

	rows, err := db.db.QueryContext(ctx, `select note_number, text, created from notes.notes 
	join notes.notes_view on notes.notes_view.id = notes.notes.id
	where notes.notes.user_id = (select id from users.users where tg_id = $1) and created >= $2::date
	AND created <= $3::date;`, searchNote.TgID, searchNote.FirstDate, searchNote.SecondDate)
	if err != nil {
		return nil, fmt.Errorf("error while searching notes by two dates: %w", err)
	}

	for rows.Next() {
		note := model.Note{}

		err := rows.Scan(&note.ViewID, &note.Text, &note.Created)
		if err != nil {
			return nil, fmt.Errorf("error while scanning note (search by two dates): %w", err)
		}

		notes = append(notes, note)
	}

	if len(notes) == 0 {
		return nil, api_errors.ErrNotesNotFound
	}

	return notes, nil
}
