package note

import (
	"context"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	tele "gopkg.in/telebot.v3"
)

type saveNoteHandler struct {
	srv notes
}

type notes interface {
	GetUserID(ctx context.Context, tgID int64) (int, error)
	SaveNote(ctx context.Context, note model.Note) error
}

func NewSaveNoteHandler(srv notes) *saveNoteHandler {
	return &saveNoteHandler{srv}
}

func (h *saveNoteHandler) Handle(ctx tele.Context) error {
	c, cancel := context.WithCancel(context.TODO()) // тот ли контекст?
	defer cancel()

	note := model.Note{
		TgID: ctx.Chat().ID,
		Text: ctx.Text(),
	}

	if err := h.srv.SaveNote(c, note); err != nil {
		return err
	}

	return ctx.Send(messages.SuccessfullyCreatedNoteMessage)
}
