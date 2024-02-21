package controller

import (
	"context"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"gopkg.in/telebot.v3"
)

// MinutesDuration принимает от пользователя количество минут, в которые нужно присылать уведомления
func (c *Controller) MinutesDuration(ctx context.Context, telectx telebot.Context) error {
	err := c.reminderSrv.ProcessMinutes(telectx.Chat().ID, telectx.Message().Text)
	if err != nil {
		return telectx.EditOrSend(messages.InvalidMinutesMessage, view.BackToReminderMenuBtns())
	}

	return c.saveReminder(ctx, telectx)
}

// HoursDuration принимает от пользователя количество часов, в которые нужно присылать уведомления
func (c *Controller) HoursDuration(ctx context.Context, telectx telebot.Context) error {
	err := c.reminderSrv.ProcessHours(telectx.Chat().ID, telectx.Message().Text)
	if err != nil {
		return telectx.EditOrSend(messages.InvalidHoursMessage, view.BackToReminderMenuBtns())
	}

	return c.saveReminder(ctx, telectx)
}

// DaysDuration принимает от пользователя количество дней, в которые нужно присылать уведомления
func (c *Controller) DaysDuration(ctx context.Context, telectx telebot.Context) error {
	err := c.reminderSrv.ProcessDaysDuration(telectx.Chat().ID, telectx.Message().Text)
	if err != nil {
		return err
	}

	err = c.reminderSrv.SaveDate(telectx.Chat().ID, telectx.Message().Text)
	if err != nil {
		return err
	}

	return telectx.EditOrSend(messages.ReminderTimeMessage, view.BackToReminderMenuBtns())
}
