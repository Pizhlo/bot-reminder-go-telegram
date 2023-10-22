package server

import (
	"context"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	tele "gopkg.in/telebot.v3"
)

func (s *Server) CreateNote(ctx tele.Context) error {
	user, err := s.UserCacheEditor.GetUser(ctx.Chat().ID)
	if err != nil {
		return err
	}
	note := model.Note{UserID: user.ID, Text: ctx.Text()}
	c, cancel := context.WithCancel(context.TODO()) // тот ли контекст?
	defer cancel()

	if err := s.NoteEditor.CreateNote(c, note); err != nil {
		return err
	}
	return ctx.Send(messages.SuccessfullyCreatedNoteMessage)
}
