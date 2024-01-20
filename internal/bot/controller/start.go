package controller

import (
	"context"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	tele "gopkg.in/telebot.v3"
)

func (c *Controller) Start(ctx context.Context, telectx tele.Context) error {
	if !c.CheckUser(ctx, telectx.Chat().ID) {
		return c.Location(ctx, telectx)
	}
	return telectx.Send(messages.StartMessage, tele.RemoveKeyboard)
}

func (c *Controller) Location(ctx context.Context, telectx tele.Context) error {
	locMenu := &tele.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: true}

	locBtn := locMenu.Location("Отправить геолокацию")
	rejectBtn := locMenu.Text("Отказаться")

	locMenu.Reply(
		locMenu.Row(locBtn),
		locMenu.Row(rejectBtn),
	)

	return telectx.Send(messages.StartMessageLocation, locMenu)
}
