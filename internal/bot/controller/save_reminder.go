package controller

import (
	"context"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"gopkg.in/telebot.v3"
)

// SaveReminder сохраняет напоминание: достает часовой пояс, чтобы установить поле created, и передает в reminder service
func (c *Controller) SaveReminder(ctx context.Context, telectx telebot.Context) error {
	// достаем часовой пояс пользователя, чтобы установить поле created
	tz, err := c.userSrv.GetTimezone(ctx, telectx.Chat().ID)
	if err != nil {
		return err
	}

	err = c.reminderSrv.Save(ctx, telectx.Chat().ID, tz)
	if err != nil {
		return err
	}

	return telectx.EditOrSend(messages.SuccessCreationMessage, view.BackToMenuBtn())
}
