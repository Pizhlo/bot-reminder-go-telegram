package controller

import (
	"context"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"gopkg.in/telebot.v3"
)

// MinutesDuration принимает от пользователя количество минут, в которые нужно присылать уведомления
func (c *Controller) MinutesDuration(ctx context.Context, telectx telebot.Context) error {
	err := c.reminderSrv.ProcessMinutes(telectx.Chat().ID, telectx.Message().Text)
	if err != nil {
		return telectx.EditOrSend(messages.InvalidMinutesMessage)
	}

	return c.saveReminder(ctx, telectx)
}

// HoursDuration принимает от пользователя количество часов, в которые нужно присылать уведомления
func (c *Controller) HoursDuration(ctx context.Context, telectx telebot.Context) error {
	err := c.reminderSrv.ProcessHours(telectx.Chat().ID, telectx.Message().Text)
	if err != nil {
		return telectx.EditOrSend(messages.InvalidHoursMessage)
	}

	return c.saveReminder(ctx, telectx)
}
