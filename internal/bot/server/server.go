package server

import (
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/calendar"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	tele "gopkg.in/telebot.v3"
)

type Server struct {
	NoteEditor     noteEditor
	ReminderEditor reminderEditor
	UserEditor     userEditor
	Calendar       *calendar.Calendar
	Logger         *logger.Logger
	Bot            *tele.Bot
}

type noteEditor interface{}
type reminderEditor interface{}
type userEditor interface{}

func New(noteEditor noteEditor, reminderEditor reminderEditor, userEditor userEditor, calendar *calendar.Calendar, logger *logger.Logger, bot *tele.Bot) *Server {
	return &Server{noteEditor, reminderEditor, userEditor, calendar, logger, bot}
}
