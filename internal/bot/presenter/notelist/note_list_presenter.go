package notelist

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view/notelist"
)

type Standard struct {
	service _NoteService
	view    notelist.NoteList
}

var _ NoteList = (*Standard)(nil)

func NewStandard(service _NoteService, view notelist.NoteList) *Standard {
	return &Standard{
		service: service,
		view:    view,
	}
}

func (p *Standard) UpdateNoteList(ctx context.Context) error {
	model, err := p.view.GetViewModel(ctx)
	if err != nil {
		return err
	}

	notes, err := p.service.GetNotes(ctx, *model)
	if err != nil {
		return err
	}

	return p.view.SetNotes(ctx, notes)
}
