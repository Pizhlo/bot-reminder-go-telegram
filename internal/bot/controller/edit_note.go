package controller

import (
	"context"
	"fmt"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	tele "gopkg.in/telebot.v3"
)

// AskNoteText запрашивает новый текст существующей заметки
func (c *Controller) AskNoteText(ctx context.Context, telectx tele.Context) error {
	return telectx.EditOrSend(messages.AskNewNoteTextMessage, view.BackToMenuBtn())
}

func (c *Controller) UpdateNote(ctx context.Context, telectx tele.Context, viewID int) error {
	err := c.noteSrv.EditNote(ctx, telectx.Chat().ID, int64(viewID), telectx.Message().Text)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf(messages.EditNoteSuccessMessage, viewID)
	return telectx.EditOrSend(msg, view.BackToMenuAndNotesBtn())
}
