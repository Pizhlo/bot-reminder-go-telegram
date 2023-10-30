package note

import (
	"context"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/server"
	tele "gopkg.in/telebot.v3"
)

type NotesHandler struct {
	srv *server.Server
}

func NewNotesHandler(srv *server.Server) *NotesHandler {
	return &NotesHandler{srv}
}

func (h *NotesHandler) Handle(ctx tele.Context) error {
	c, cancel := context.WithCancel(context.TODO())
	defer cancel()

	id, err := h.srv.GetUserID(c, ctx.Chat().ID)
	if err != nil {
		return err
	}

	notes, err := h.srv.GetAllNotes(c, id)
	if err != nil {
		return err
	}

	messages := h.makeMessage(notes)

	for _, msg := range messages {
		err := ctx.Send(msg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *NotesHandler) makeMessage(notes []model.Note) []string {
	sumLen := 0
	var message string
	var messages []string

	for _, note := range notes {
		message += fmt.Sprintf("%s - создано %s\n\nДля удаления нажмите /del%d\n\n", note.Text, note.Created.Format("02.01.2006 в 15:04"), note.ID)
		sumLen += len([]rune(message))

		if sumLen >= 3000 { //  telegram: message is too long (400) - length must not be more than 4096
			fmt.Println(sumLen)
			messages = append(messages, message)
			sumLen = 0
			message = ""
		}
	}

	return messages
}
