package note

import (
	"context"

	messageeditor "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/handler/message_editor"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	tele "gopkg.in/telebot.v3"
)

type notesHandler struct {
	srv           noteServer
	messageEditor messageEditor
}

type messageEditor interface {
	MakeNotesMessage() []string
}

type noteServer interface {
	GetUserID(ctx context.Context, tgID int64) (int, error)
	GetAllNotes(ctx context.Context, tgID int64) ([]model.Note, error)
}

func NewNotesHandler(srv noteServer) *notesHandler {
	return &notesHandler{srv, nil}
}

func (h *notesHandler) Handle(ctx tele.Context) error {
	c, cancel := context.WithCancel(context.TODO())
	defer cancel()

	notes, err := h.srv.GetAllNotes(c, ctx.Chat().ID)
	if err != nil {
		return err
	}

	messageEditor := messageeditor.New(notes)
	h.messageEditor = messageEditor
	messages := h.messageEditor.MakeNotesMessage()

	for _, msg := range messages {
		err := ctx.Send(msg)
		if err != nil {
			return err
		}
	}

	return nil
}
