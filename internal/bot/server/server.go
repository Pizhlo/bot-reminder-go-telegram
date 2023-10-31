package server

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

// communicates with db
type Server struct {
	noteEditor          noteEditor
	reminderEditor      reminderEditor
	userEditor          userEditor
	timezoneCacheEditor timezoneCacheEditor
	userCacheEditor     userCacheEditor
}

// db
type noteEditor interface {
	SaveNote(ctx context.Context, note model.Note) error
	GetAllNotes(ctx context.Context, id int) ([]model.Note, error)
	SearchNotesByText(ctx context.Context, query model.SearchNote) ([]model.Note, error)
}

type reminderEditor interface{}
type userEditor interface {
	GetUser(ctx context.Context, tgID int64) (model.User, error)
	SaveUser(ctx context.Context, telegramID int64) (int, error)
}

// cache
type timezoneCacheEditor interface {
	GetUserTimezone(id int64) (model.UserTimezone, error)
	SaveUserTimezone(id int64, tz model.UserTimezone)
}
type userCacheEditor interface {
	GetUser(ctx context.Context, tgID int64) (model.User, error)
	SaveUser(ctx context.Context, id int, tgID int64)
}

func New(noteEditor noteEditor, reminderEditor reminderEditor, userEditor userEditor, tzCache timezoneCacheEditor, userCacheEditor userCacheEditor) *Server {
	return &Server{noteEditor, reminderEditor, userEditor, tzCache, userCacheEditor}
}
