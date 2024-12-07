package controller

import (
	"context"
	"fmt"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

// SendReminder отправляет пользователю напоминание в указанное время
func (c *Controller) SendReminder(ctx context.Context, reminder *model.Reminder) error {
	logrus.Debugf("Sending reminder to: %d. Reminder: %+v\n", reminder.TgID, reminder)

	msg, err := view.ReminderMessage(reminder)
	if err != nil {
		return err
	}

	c.bot.Handle(&view.BtnDeleteReminder, func(telectx telebot.Context) error {
		if err := c.ProcessDeleteReminder(ctx, telectx, reminder); err != nil {
			c.HandleError(telectx, err, "reminder_works")
		}

		return nil
	})

	// _, err = c.bot.Send(&telebot.Chat{ID: reminder.TgID}, msg, &telebot.SendOptions{
	// 	ParseMode:   htmlParseMode,
	// 	ReplyMarkup: view.ReminderWorkMenu(),
	// })

	err = c.sendMessage(reminder.TgID, msg, view.ReminderWorkMenu())

	if err != nil {
		logrus.Errorf("error while sending reminder to user: %v", err)
		msg := fmt.Sprintf(messages.ErrorSendingReminderMessage, err)
		c.bot.Send(&telebot.Chat{ID: c.channelID}, msg, &telebot.SendOptions{
			ParseMode: htmlParseMode,
		})
	}

	return nil
}

// ProcessDeleteReminder обрабатывает кнопку "Удалить" у сработавшего напоминания
func (c *Controller) ProcessDeleteReminder(ctx context.Context, telectx telebot.Context, r *model.Reminder) error {
	err := c.reminderSrv.DeleteReminder(ctx, r)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf(messages.ReminderDeletedMessage, r.Name)

	return telectx.Edit(msg, &telebot.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.BackToRemindersAndMenu(),
	})
}
