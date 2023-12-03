package note

import (
	"context"
	"errors"
	"time"
)

type Note struct {
	ID        int
	UserID    int
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SearchParams struct {
	UserID int
	Terms  []string
}

var ErrNotFound = errors.New("not found")

type Repo interface {
	Add(ctx context.Context, userID int, text string) (*Note, error)
	//Get(ctx context.Context, id int) (*Note, error)
	//Update(ctx context.Context, id int, updFun func(*Note) (*Note, error)) (*Note, error)
	FindByParams(ctx context.Context, params *SearchParams) ([]*Note, error)
}
