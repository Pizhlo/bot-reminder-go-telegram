package controller

import (
	"context"
	"errors"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	tele "gopkg.in/telebot.v3"
)

// ListNotes возвращает все заметки пользователя и отправляет ему
func (c *Controller) ListNotes(ctx context.Context, telectx tele.Context) error {
	c.logger.Debugf("Controller: handling /notes command.\n")

	message, kb, err := c.noteSrv.GetAll(ctx, telectx.Chat().ID)
	if err != nil {
		if errors.Is(err, api_errors.ErrNotesNotFound) {
			return telectx.Send(messages.NotesNotFoundMessage)
		}

		c.logger.Errorf("Error while handling /notes command. User ID: %d. Error: %+v\n", telectx.Chat().ID, err)

		c.HandleError(telectx, err)

		return err
	}

	c.logger.Debugf("Controller: successfully got all user's notes. Sending message to user...\n")
	return telectx.Send(message, kb)
}

// NextPageNotes обрабатывает кнопку переключения на следующую страницу
func (c *Controller) NextPageNotes(ctx context.Context, telectx tele.Context) error {
	c.logger.Debugf("Controller: handling next notes page command.\n")
	next, kb := c.noteSrv.NextPage(telectx.Chat().ID)

	return telectx.Edit(next, kb)
}

// NextPageNotes обрабатывает кнопку переключения на предыдущую страницу
func (c *Controller) PrevPageNotes(ctx context.Context, telectx tele.Context) error {
	c.logger.Debugf("Controller: handling previous notes page command.\n")
	next, kb := c.noteSrv.PrevPage(telectx.Chat().ID)

	return telectx.Edit(next, kb)
}

// NextPageNotes обрабатывает кнопку переключения на последнюю страницу
func (c *Controller) LastPageNotes(ctx context.Context, telectx tele.Context) error {
	c.logger.Debugf("Controller: handling last notes page command.\n")
	next, kb := c.noteSrv.LastPage(telectx.Chat().ID)

	return telectx.Edit(next, kb)
}

// NextPageNotes обрабатывает кнопку переключения на первую страницу
func (c *Controller) FirstPageNotes(ctx context.Context, telectx tele.Context) error {
	c.logger.Debugf("Controller: handling first notes page command.\n")
	next, kb := c.noteSrv.FirstPage(telectx.Chat().ID)

	return telectx.Edit(next, kb)
}
