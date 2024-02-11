package controller

import (
	"context"
	"fmt"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	tele "gopkg.in/telebot.v3"
)

// StartCmd отправляет приветственное сообщение и меню
func (c *Controller) StartCmd(ctx context.Context, telectx tele.Context) error {
	c.logger.Debugf("Controller: handling /start (or menu btn)\n")

	kb := view.MainMenu()

	text := fmt.Sprintf(messages.StartMessage, telectx.Chat().FirstName)

	return telectx.EditOrSend(text, kb)
}
