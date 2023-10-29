package note

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/server"
	tele "gopkg.in/telebot.v3"
)

type SaveNoteHandler struct {
	srv *server.Server
}

func NewSaveNoteHandler(srv *server.Server) *SaveNoteHandler {
	return &SaveNoteHandler{srv}
}

func (h *SaveNoteHandler) Handle(ctx tele.Context) error {
	c, cancel := context.WithCancel(context.TODO()) // тот ли контекст?
	defer cancel()

	id, err := h.srv.GetUserID(c, ctx.Chat().ID)
	if err != nil {
		return err
	}

	note := model.Note{
		UserID: id,
		Text:   ctx.Text(),
	}

	return h.srv.SaveNote(note)
}
