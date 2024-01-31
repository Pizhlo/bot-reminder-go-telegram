package controller

import (
	"context"
	"fmt"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	tele "gopkg.in/telebot.v3"
)

// Start обрабатывает команду старт: проверяет, зарегистрирован ли пользователь,
// и в зависимости от этого либо запрашивает геолокацию, либо отправляет приветственное сообщение
func (c *Controller) Start(ctx context.Context, telectx tele.Context) error {
	c.logger.Debugf("Controller: handling /start. Checking user...\n")

	if !c.CheckUser(ctx, telectx.Chat().ID) {
		c.logger.Debugf("Controller: user is unknown. Sending location request...\n")
		return c.Location(ctx, telectx)
	}

	c.logger.Debugf("Controller: user is known. Sending start message...\n")

	text := fmt.Sprintf(messages.StartMessage, telectx.Chat().FirstName)

	return telectx.Send(text, tele.RemoveKeyboard, &tele.SendOptions{
		ParseMode: htmlParseMode,
	})
}

// Location запрашивает геолокацию у пользователя
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
