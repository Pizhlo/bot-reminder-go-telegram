package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/gocron"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/google/uuid"
	"gopkg.in/telebot.v3"
)

// saveReminder сохраняет напоминание
func (c *Controller) saveReminder(ctx context.Context, telectx telebot.Context) error {
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

	err = c.reminderSrv.Save(ctx, telectx.Chat().ID)
	if err != nil {
		return err
	}

	// создаем отложенный вызов
	jobID, err := c.createReminder(ctx, telectx)
	if err != nil {
		return err
	}

	// сохраняем задачу в базе
	err = c.reminderSrv.SaveJobID(ctx, jobID, telectx.Chat().ID)
	if err != nil {
		return err
	}

	r, err := c.reminderSrv.GetFromMemory(telectx.Chat().ID)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf(messages.SuccessCreationMessage, r.Name)

	return telectx.EditOrSend(msg, view.BackToMenuBtn())
}

func (c *Controller) createReminder(ctx context.Context, telectx telebot.Context) (uuid.UUID, error) {
	r, err := c.reminderSrv.GetFromMemory(telectx.Chat().ID)
	if err != nil {
		return uuid.UUID{}, err
	}

	params := gocron.FuncParams{
		Ctx:      telectx,
		Reminder: *r,
	}

	switch r.Type {
	case model.EverydayType:
		return c.scheduler.CreateEverydayJob(r.Time, c.SendReminder, params)
	case model.SeveralTimesDayType:
		return c.scheduler.CreateMinutesReminder(r.Time, c.SendReminder, params)
	default:
		return uuid.UUID{}, fmt.Errorf("unknown type of reminder: %s", r.Type)
	}
}
