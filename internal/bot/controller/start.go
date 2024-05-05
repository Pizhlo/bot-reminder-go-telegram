package controller

import (
	"context"
	"fmt"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// StartCmd отправляет приветственное сообщение и меню на команду /start
func (c *Controller) StartCmd(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Controller: handling /start (or menu btn)\n")

	kb := view.MainMenu()

	text := fmt.Sprintf(messages.StartMessage, telectx.Chat().FirstName)

	return telectx.EditOrSend(text, kb)
}

// MenuCmd обрабатывает команду /menu
func (c *Controller) MenuCmd(ctx context.Context, telectx tele.Context) error {
	return telectx.EditOrSend(messages.MenuMessage, view.MainMenu())
}

// HelpCmd обрабатывает команду /help
func (c *Controller) HelpCmd(ctx context.Context, telectx tele.Context) error {
	sendOpts := &tele.SendOptions{}
	if c.userSrv.CheckUser(ctx, telectx.Chat().ID) {
		sendOpts = &tele.SendOptions{
			ReplyMarkup: view.MainMenu(),
			ParseMode:   htmlParseMode,
		}
	} else {
		sendOpts = &tele.SendOptions{
			ParseMode: htmlParseMode,
		}
	}

	msg := fmt.Sprintf(messages.HelpMessage, telectx.Sender().FirstName)
	return telectx.EditOrSend(msg, sendOpts)
}
