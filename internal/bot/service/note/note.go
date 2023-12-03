package note

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/note"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/account"
)

type Service interface {
	GetNotes(ctx context.Context, query note.SearchParams) ([]*note.Note, error)
	AddNote(ctx context.Context, userID int, text string) (*note.Note, error)
	//UpdateNote(ctx context.Context, id int, text string) (*note.Note, error)
}

type Standard struct {
	users account.Service
	notes note.Repo
}

var _ Service = (*Standard)(nil)

func NewStandard(users account.Service, notes note.Repo) *Standard {
	return &Standard{
		users: users,
		notes: notes,
	}
}

func (p *Standard) AddNote(ctx context.Context, userID int, text string) (*note.Note, error) {
	u, err := p.users.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("cannot add a note: %w", err)
	}

	loc, err := time.LoadLocation(u.Timezone.Name)
	if err != nil {
		return nil, fmt.Errorf("error while setting timezone for created field in note table: %w", err)
	}

	created := time.Now().In(loc)

	n, err := p.notes.Add(ctx, userID, text, created)
	if err != nil {
		return nil, fmt.Errorf("cannot add a note: %w", err)
	}

	return n, nil
}

func (p *Standard) GetNotes(ctx context.Context, query note.SearchParams) ([]*note.Note, error) {
	_, err := p.users.GetUser(ctx, query.UserID)
	if errors.Is(err, account.ErrUserNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot get notes: %w", err)
	}

	nn, err := p.notes.FindByParams(ctx, &query)
	if errors.Is(err, note.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot get notes: %w", err)
	}

	return nn, nil
}

// func (p *Standard) UpdateNote(ctx context.Context, id int, text string) (*note.Note, error) {
// 	n, err := p.notes.Update(ctx, id, func(n *note.Note) (*note.Note, error) {
// 		return &note.Note{
// 			Text: text,
// 		}, nil
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot update a %d note: %w", id, err)
// 	}

// 	return n, nil
// }
