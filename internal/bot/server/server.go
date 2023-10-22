package server

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/calendar"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	tele "gopkg.in/telebot.v3"
)

type Server struct {
	NoteEditor          noteEditor
	ReminderEditor      reminderEditor
	UserEditor          userEditor
	TimezoneCacheEditor timezoneCacheEditor
	UserCacheEditor     userCacheEditor
	Calendar            *calendar.Calendar
	Logger              *logger.Logger
	Bot                 *tele.Bot
}

// db
type noteEditor interface {
	CreateNote(ctx context.Context, note model.Note) error
}
type reminderEditor interface{}
type userEditor interface {
	SaveUser(telegramID int64) (int, error)
}

// cache
type timezoneCacheEditor interface {
	GetUserTimezone(id int64) (model.UserTimezone, error)
	SaveUserTimezone(id int64, tz model.UserTimezone)
}
type userCacheEditor interface {
	GetUser(tgID int64) (model.User, error)
	SaveUser(id int, tgID int64)
}

func New(noteEditor noteEditor, reminderEditor reminderEditor, userEditor userEditor, tzCache timezoneCacheEditor, userCacheEditor userCacheEditor, calendar *calendar.Calendar, logger *logger.Logger, bot *tele.Bot) *Server {
	return &Server{noteEditor, reminderEditor, userEditor, tzCache, userCacheEditor, calendar, logger, bot}
}
