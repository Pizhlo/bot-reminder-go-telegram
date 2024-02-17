package controller

import (
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"gopkg.in/telebot.v3"
)

// SendReminder отправляет пользователю напоминание в указанное время
func (c *Controller) SendReminder(ctx telebot.Context, reminder model.Reminder) error {
	c.logger.Debugf("Sending reminder to: %d\n", ctx.Chat().ID)

	err := ctx.EditOrSend(view.ReminderMessage(reminder), &telebot.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.DeleteReminderBtn(reminder),
	})

	if err != nil {
		c.HandleError(ctx, err)
	}

	return nil
}
