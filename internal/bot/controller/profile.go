package controller

import (
	"context"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	tele "gopkg.in/telebot.v3"
)

func (c *Controller) Profile(ctx context.Context, telectx tele.Context) error {
	return telectx.EditOrSend(messages.ProfileMessage, view.ProfileMenu())
}
