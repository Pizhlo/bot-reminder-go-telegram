package controller

import (
	"context"
	"fmt"
	"time"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	gocron "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/scheduler"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
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
	nextRun, err := c.createReminder(ctx, telectx)
	if err != nil {
		return err
	}

	// сохраняем задачу в базе
	err = c.reminderSrv.SaveJobID(ctx, nextRun.JobID, telectx.Chat().ID)
	if err != nil {
		return err
	}

	r, err := c.reminderSrv.GetFromMemory(telectx.Chat().ID)
	if err != nil {
		return err
	}

	layout := "02.01.2006 15:04:05"

	nextRunMsg := view.ProcessTypeAndDate(r.Type, r.Date, r.Time)

	msg := fmt.Sprintf(messages.SuccessCreationMessage, r.Name, nextRunMsg, nextRun.NextRun.Format(layout))

	return telectx.EditOrSend(msg, &telebot.SendOptions{
		ReplyMarkup: view.BackToMenuBtn(),
		ParseMode:   htmlParseMode,
	})
}

func (c *Controller) createReminder(ctx context.Context, telectx telebot.Context) (gocron.NextRun, error) {
	r, err := c.reminderSrv.GetFromMemory(telectx.Chat().ID)
	if err != nil {
		return gocron.NextRun{}, err
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
		return gocron.NextRun{}, fmt.Errorf("unknown type of reminder: %s", r.Type)
	}
}
