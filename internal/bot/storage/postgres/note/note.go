package note

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/note"
)

type NoteRepo struct {
	db *sql.DB
}

func New(dbURl string) (*NoteRepo, error) {
	db, err := sql.Open("postgres", dbURl)
	if err != nil {
		return nil, fmt.Errorf("connect open a db driver: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("cannot connect to a db: %w", err)
	}
	return &NoteRepo{db}, nil
}

// func (db *NoteRepo) GetAllNotes(ctx context.Context, id int) ([]model.Note, error) {
// 	var res []model.Note

// 	fmt.Println("userid = ", id)

// 	rows, err := db.Query(ctx, `select * from notes where "user_id" = $1`, id)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "error while query")
// 	}

// 	for rows.Next() {
// 		note := model.Note{}
// 		err := rows.Scan(&note.ID, &note.UserID, &note.Text, &note.Created)
// 		if err != nil {
// 			return nil, errors.Wrap(err, "error while scanning")
// 		}

// 		res = append(res, note)
// 	}

// 	return res, nil
// }

// func (db *NoteRepo) SaveNote(ctx context.Context, note model.Note) error {
// 	_, err := db.Exec(ctx, `insert into notes("user_id", text, created) values($1, $2, $3)`, note.UserID, note.Text, time.Now()) // добавить user timezone
// 	if err != nil {
// 		return errors.Wrap(err, "error while query")
// 	}

// 	return nil
// }

type Note struct {
	ID        int
	UserID    int
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (db *NoteRepo) Save(ctx context.Context, note model.Note) error {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("error while creating transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, `insert into notes.notes (user_id, text, created) values((select id from users.users where tg_id=$1), $2, $3) returning id`, note.TgID, note.Text, note.Created)
	if err != nil {
		return fmt.Errorf("error inserting note: %w", err)
	}

	return tx.Commit()
}

func (db *NoteRepo) GetAllByUserID(ctx context.Context, userID int64) ([]model.Note, error) {
	notes := make([]model.Note, 0)

	rows, err := db.db.QueryContext(ctx, `select id, text, created from notes.notes where user_id = (select id from users.users where tg_id = $1)`, userID)
	if err != nil {
		return nil, fmt.Errorf("error while getting all notes from DB by user ID %d: %w", userID, err)
	}
	defer rows.Close()

	for rows.Next() {
		note := model.Note{}

		err := rows.Scan(&note.ID, &note.Text, &note.Created)
		if err != nil {
			return nil, fmt.Errorf("error while scanning note: %w", err)
		}

		notes = append(notes, note)
	}

	if len(notes) == 0 {
		return nil, api_errors.ErrNotesNotFound
	}

	return notes, nil
}

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

type SearchParams struct {
	UserID int
	Terms  []string
}

func (db *NoteRepo) FindByParams(ctx context.Context, params *note.SearchParams) ([]*note.Note, error) {

	return nil, nil
}

// func (db *NoteRepo) SearchNotesByText(ctx context.Context, query model.SearchNote) ([]model.Note, error) {
// 	var notes []model.Note

// 	rows, err := db.Query(ctx, `select * from notes where user_id = $1 and "text" LIKE '%' || $2 || '%'`, query.UserID, query.Text)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "error while making query")
// 	}

// 	for rows.Next() {
// 		note := model.Note{}

// 		err := rows.Scan(&note.ID, &note.UserID, &note.Text, &note.Created)
// 		if err != nil {
// 			return nil, errors.Wrap(err, "error while scanning")
// 		}

// 		notes = append(notes, note)
// 	}

// 	if len(notes) == 0 {
// 		return nil, api_errors.ErrNotesNotFound
// 	}

// 	return notes, nil
// }
