package note

import (
	"context"
	"fmt"
	"time"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/pkg/errors"

	"github.com/jackc/pgx/v4/pgxpool"
)

type NoteRepo struct {
	*pgxpool.Pool
}

func New(conn *pgxpool.Pool) *NoteRepo {
	return &NoteRepo{conn}
}

func (db *NoteRepo) GetAllNotes(ctx context.Context, id int) ([]model.Note, error) {
	var res []model.Note

	fmt.Println("userid = ", id)

	rows, err := db.Query(ctx, `select * from notes where "user_id" = $1`, id)
	if err != nil {
		return nil, errors.Wrap(err, "error while query")
	}

	for rows.Next() {
		note := model.Note{}
		err := rows.Scan(&note.ID, &note.UserID, &note.Text, &note.Created)
		if err != nil {
			return nil, errors.Wrap(err, "error while scanning")
		}

		res = append(res, note)
	}

	return res, nil
}

func (db *NoteRepo) SaveNote(ctx context.Context, note model.Note) error {
	_, err := db.Exec(ctx, `insert into notes("user_id", text, created) values($1, $2, $3)`, note.UserID, note.Text, time.Now()) // добавить user timezone
	if err != nil {
		return errors.Wrap(err, "error while query")
	}

	return nil
}

func (db *NoteRepo) SearchNotesByText(ctx context.Context, query model.SearchNote) ([]model.Note, error) {
	var notes []model.Note

	rows, err := db.Query(ctx, `select * from notes where user_id = $1 and "text" LIKE '%' || $2 || '%'`, query.UserID, query.Text)
	if err != nil {
		return nil, errors.Wrap(err, "error while making query")
	}

	for rows.Next() {
		note := model.Note{}

		err := rows.Scan(&note.ID, &note.UserID, &note.Text, &note.Created)
		if err != nil {
			return nil, errors.Wrap(err, "error while scanning")
		}

		notes = append(notes, note)
	}

	if len(notes) == 0 {
		return nil, api_errors.ErrNotesNotFound
	}

	return notes, nil
}
