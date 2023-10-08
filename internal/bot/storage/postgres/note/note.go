package note

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type NoteRepo struct {
	*pgxpool.Pool
}

func New(conn *pgxpool.Pool) *NoteRepo {
	return &NoteRepo{conn}
}
