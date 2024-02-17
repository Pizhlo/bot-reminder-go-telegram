package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/gocron"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"gopkg.in/telebot.v3"
)

// ReminderTime обрабатывает время напоминания, сохраняет и создает отложенный вызов
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

	// достаем часовой пояс пользователя, чтобы установить поле created
	tz, err := c.userSrv.GetTimezone(ctx, telectx.Chat().ID)
	if err != nil {
		return err
	}

	userTz, err := time.LoadLocation(tz.Name)
	if err != nil {
		return fmt.Errorf("error while setting timezone (for setting 'created' field): %w", err)
	}

	// сохраняем поле created
	c.reminderSrv.SaveCreatedField(telectx.Chat().ID, userTz)

	// сохраняем напоминание, потому что назначение времени - последний этап
	err = c.reminderSrv.Save(ctx, telectx.Chat().ID)
	if err != nil {
		return err
	}

	r, err := c.reminderSrv.GetFromMemory(telectx.Chat().ID)
	if err != nil {
		return err
	}

	params := gocron.FuncParams{
		Ctx:      telectx,
		Reminder: *r,
	}

	// создаем отложенный вызов
	jobID, err := c.scheduler.CreateEverydayJob(telectx.Message().Text, c.SendReminder, params)
	if err != nil {
		return err
	}

	// сохраняем задачу в базе
	err = c.reminderSrv.SaveJobID(ctx, jobID, telectx.Chat().ID)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf(messages.SuccessCreationMessage, r.Name)

	return telectx.EditOrSend(msg, view.BackToMenuBtn())
}
