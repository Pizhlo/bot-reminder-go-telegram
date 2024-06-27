package controller

import (
	"context"
	"fmt"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	tele "gopkg.in/telebot.v3"
)

// AskNoteText запрашивает новый текст существующей заметки
func (c *Controller) AskNoteText(ctx context.Context, telectx tele.Context) error {
	return telectx.EditOrSend(messages.AskNewNoteTextMessage)
}

func (c *Controller) UpdateNote(ctx context.Context, telectx tele.Context, noteID int) error {
	msg := fmt.Sprintf("отредактирована заметка: %+v", noteID)
	return telectx.EditOrSend(msg)
}
