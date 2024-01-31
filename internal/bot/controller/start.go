package controller

import (
	"context"
	"fmt"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	tele "gopkg.in/telebot.v3"
)

func (c *Controller) Start(ctx context.Context, telectx tele.Context) error {
	c.logger.Debugf("Handling /start. Checking user...\n")

	if !c.CheckUser(ctx, telectx.Chat().ID) {
		c.logger.Debugf("User is unknown. Sending location request...\n")
		return c.Location(ctx, telectx)
	}

	c.logger.Debugf("User is known. Sending start message...\n")

	text := fmt.Sprintf(messages.StartMessage, telectx.Chat().FirstName)

	return telectx.Send(text, tele.RemoveKeyboard, &tele.SendOptions{
		ParseMode: htmlParseMode,
	})
}

func (c *Controller) Location(ctx context.Context, telectx tele.Context) error {
	locMenu := &tele.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: true}

	locBtn := locMenu.Location("Отправить геолокацию")
	rejectBtn := locMenu.Text("Отказаться")

	locMenu.Reply(
		locMenu.Row(locBtn),
		locMenu.Row(rejectBtn),
	)

	txt := fmt.Sprintf(messages.StartMessageLocation, telectx.Chat().FirstName)

	return telectx.Send(txt, locMenu)
}
