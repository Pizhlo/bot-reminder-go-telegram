package note

import (
	"context"
	"errors"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	messageeditor "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/handler/message_editor"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	tele "gopkg.in/telebot.v3"
)

type searchByTextHandler struct {
	srv           notesEditor
	messageEditor messageEditor
}

type notesEditor interface { // server
	GetUserID(ctx context.Context, tgID int64) (int, error)
	SearchNotesByText(ctx context.Context, query model.SearchNote) ([]model.Note, error)
}

func NewSearchByTextHandler(srv notesEditor) *searchByTextHandler {
	return &searchByTextHandler{srv: srv}
}

func (h *searchByTextHandler) Handle(ctx tele.Context) error {
	c, cancel := context.WithCancel(context.TODO())
	defer cancel()

	query := model.SearchNote{
		TgID: ctx.Chat().ID,
		Text: ctx.Text(),
	}

	notes, err := h.srv.SearchNotesByText(c, query)
	if err != nil {
		if errors.Is(err, api_errors.ErrNotesNotFound) {
			return ctx.Send(messages.NotesNotFoundMessage)
		}
		return err
	}

	messageEditor := messageeditor.New(notes)
	h.messageEditor = messageEditor
	notesMessages := h.messageEditor.MakeNotesMessage()

	err = ctx.Send(messages.FoundNotesMessage)
	if err != nil {
		return err
	}

	for _, msg := range notesMessages {
		err := ctx.Send(msg)
		if err != nil {
			return err
		}
	}

	return nil
}
