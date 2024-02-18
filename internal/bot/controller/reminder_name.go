package controller

import (
	"context"
	"fmt"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"gopkg.in/telebot.v3"
)

// ReminderName обрабатывает название напоминания
func (c *Controller) ReminderName(ctx context.Context, telectx telebot.Context) error {
	c.reminderSrv.SaveName(telectx.Chat().ID, telectx.Message().Text)

	txt := fmt.Sprintf(messages.TypeOfReminderMessage, telectx.Message().Text)

	return telectx.EditOrSend(txt, &telebot.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.ReminderTypes(),
	})
}
