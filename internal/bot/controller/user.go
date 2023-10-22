package controller

import (
	"errors"

	api_err "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	tele "gopkg.in/telebot.v3"
)

func (c *Controller) StartMsg(ctx tele.Context, err error) error {
	if err != nil {
		if errors.Is(err, api_err.ErrUserNotFound) {
			locationMenu.Reply(
				locationMenu.Row(locationBtn),
				locationMenu.Row(rejectBtn),
			)
			return ctx.Send(messages.StartMessageLocation, locationMenu)
		}
		return err
	}
	return ctx.Send(messages.StartMessage)
}

func (c *Controller) saveUser(telegramID int64) error {
	id, err := c.srv.UserEditor.SaveUser(telegramID)
	if err != nil {
		return err
	}

	c.srv.UserCacheEditor.SaveUser(id, telegramID)

	return nil
}
