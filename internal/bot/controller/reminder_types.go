package controller

import (
	"context"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"gopkg.in/telebot.v3"
)

// EverydayReminder обрабатывает кнопку "ежедневное напоминание"
func (c *Controller) EverydayReminder(ctx context.Context, telectx telebot.Context) error {
	// сохраняем тип напоминания - "everyday"
	c.reminderSrv.SaveType(telectx.Chat().ID, model.EverydayType)

	// сохраняем в качестве даты также строку "everyday"
	c.reminderSrv.SaveDate(telectx.Chat().ID, string(model.EverydayType))

	return telectx.EditOrSend(messages.ReminderTimeMessage, &telebot.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.BackToMenuBtn(),
	})
}
