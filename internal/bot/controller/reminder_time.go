package controller

import (
	"context"
	"time"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"gopkg.in/telebot.v3"
)

// ReminderTime обрабатывает время напоминания
func (c *Controller) ReminderTime(ctx context.Context, telectx telebot.Context) error {
	// проверяем время на валидность и сохраняем если проверка прошла успешно
	err := c.reminderSrv.ProcessTime(telectx.Chat().ID, telectx.Message().Text)
	if err != nil {
		switch err.(type) {
		case *time.ParseError:
			return telectx.EditOrSend(messages.InvalidTimeMessage, view.BackToMenuBtn())
		default:
			return err
		}
	}

	// сохраняем напоминание, потому что назначение времени - последний этап
	return c.SaveReminder(ctx, telectx)
}
