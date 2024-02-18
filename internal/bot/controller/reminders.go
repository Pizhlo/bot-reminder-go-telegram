package controller

import (
	"context"
	"errors"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	tele "gopkg.in/telebot.v3"
)

// Reminders показывает пользователю все его напоминания
func (c *Controller) Reminders(ctx context.Context, telectx tele.Context) error {
	msg, kb, err := c.reminderSrv.GetAll(ctx, telectx.Chat().ID)
	if err != nil {
		if errors.Is(err, api_errors.ErrRemindersNotFound) {
			return telectx.EditOrSend(messages.NoRemindersMessage, view.CreateReminderAndBackToMenu())
		}
		return err
	}

	return telectx.EditOrSend(msg, kb)
}
