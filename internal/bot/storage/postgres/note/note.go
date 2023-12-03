package note

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

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

func (db *NoteRepo) Add(ctx context.Context, userID int, text string) (*note.Note, error) {
	return nil, nil
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
