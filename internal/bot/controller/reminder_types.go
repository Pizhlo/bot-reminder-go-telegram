package controller

import (
	"context"
	"fmt"

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
		ReplyMarkup: view.BackToReminderMenuBtns(),
	})
}

// EverydayReminder обрабатывает кнопку "несколько раз в день"
func (c *Controller) SeveralTimesADayReminder(ctx context.Context, telectx telebot.Context) error {
	// сохраняем тип напоминания - "several times a day"
	c.reminderSrv.SaveType(telectx.Chat().ID, model.SeveralTimesDayType)

	var msg string

	// если сообщение прислал пользователь - это новое название напоминания
	if !telectx.Message().Sender.IsBot {
		// сохраняем новое название напоминания, если пользователь прислал повторно
		c.reminderSrv.SaveName(telectx.Chat().ID, telectx.Message().Text)
	}

	r, err := c.reminderSrv.GetFromMemory(telectx.Chat().ID)
	if err != nil {
		msg = messages.ChooseMinutesOrHoursWithoutNameMessage
	}

	msg = fmt.Sprintf(messages.ChooseMinutesOrHoursMessage, r.Name, "несколько раз в день")

	return telectx.EditOrSend(msg, &telebot.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.SeveralTimesBtns(),
	})
}

// OnceInMinutes обрабатывает кнопку "раз в несколько минут"
func (c *Controller) OnceInMinutes(ctx context.Context, telectx telebot.Context) error {
	c.reminderSrv.SaveDate(telectx.Chat().ID, "minutes")

	return telectx.EditOrSend(messages.MinutesDurationMessage)
}

// OnceInMinutes обрабатывает кнопку "раз в несколько часов"
func (c *Controller) OnceInHours(ctx context.Context, telectx telebot.Context) error {
	c.reminderSrv.SaveDate(telectx.Chat().ID, "hours")

	return telectx.EditOrSend(messages.HoursDurationMessage)
}
