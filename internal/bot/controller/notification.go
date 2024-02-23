package controller

import (
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"gopkg.in/telebot.v3"
)

// SendReminder отправляет пользователю напоминание в указанное время
func (c *Controller) SendReminder(ctx telebot.Context, reminder model.Reminder) error {
	c.logger.Debugf("Sending reminder to: %d\n", ctx.Chat().ID)

	msg, err := view.ReminderMessage(reminder)
	if err != nil {
		return err
	}

	err = ctx.EditOrSend(msg, &telebot.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.DeleteReminderBtn(reminder),
	})

	if err != nil {
		return err
	}

	return nil
}
