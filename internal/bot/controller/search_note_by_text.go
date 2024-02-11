package controller

import (
	"context"
	"errors"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	tele "gopkg.in/telebot.v3"
)

// AskTextForSearch запрашивает у пользователя текст заметок для поиска
func (c *Controller) AskTextForSearch(ctx context.Context, telectx tele.Context) error {
	return telectx.EditOrSend(messages.SearchNotesByTextMessage, view.BackToMenuBtn())
}

// SearchNoteByText производит поиск заметок по тексту
func (c *Controller) SearchNoteByText(ctx context.Context, telectx tele.Context) error {
	searchNote := model.SearchNoteByText{
		TgID: telectx.Chat().ID,
		Text: telectx.Message().Text,
	}

	message, kb, err := c.noteSrv.SearchByText(ctx, searchNote)
	if err != nil {
		if errors.Is(err, api_errors.ErrNotesNotFound) {
			return telectx.EditOrSend(messages.NoNotesFoundMessage, view.BackToMenuBtn())
		}

		return err
	}

	return telectx.EditOrSend(message, kb)
}
