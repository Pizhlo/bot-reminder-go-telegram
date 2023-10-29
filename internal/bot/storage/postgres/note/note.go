package note

import (
	"context"
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"

	"github.com/jackc/pgx/v4/pgxpool"
)

type NoteRepo struct {
	*pgxpool.Pool
}

func New(conn *pgxpool.Pool) *NoteRepo {
	return &NoteRepo{conn}
}

func (db *NoteRepo) SaveNote(ctx context.Context, note model.Note) error {
	_, err := db.Exec(ctx, `insert into notes("user_id", text, created) values($1, $2, $3)`, note.UserID, note.Text, time.Now()) // добавить user timezone
	if err != nil {
		return err
	}

	return nil
}
