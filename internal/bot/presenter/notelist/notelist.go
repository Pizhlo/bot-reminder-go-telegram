package notelist

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/note"
)

type NoteList interface {
	UpdateNoteList(context.Context) error
}

type _NoteService interface {
	GetNotes(context.Context, note.SearchParams) ([]*note.Note, error)
}
