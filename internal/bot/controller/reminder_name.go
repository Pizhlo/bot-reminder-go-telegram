package controller

import (
	"context"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"gopkg.in/telebot.v3"
)

// ReminderName обрабатывает название напоминания
func (c *Controller) ReminderName(ctx context.Context, telectx telebot.Context) error {
	c.reminderSrv.SaveReminderName(telectx.Chat().ID, telectx.Message().Text)

	return telectx.EditOrSend(messages.TypeOfReminderMessage, view.ReminderTypes())
}
