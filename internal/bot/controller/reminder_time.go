package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"gopkg.in/telebot.v3"
)

// ReminderTime обрабатывает время напоминания, сохраняет и создает отложенный вызов
func (c *Controller) ReminderTime(ctx context.Context, telectx telebot.Context) error {
	// проверяем время на валидность и сохраняем если проверка прошла успешно
	err := c.reminderSrv.ParseTime(telectx.Chat().ID, telectx.Message().Text)
	if err != nil {
		return err
	}

	// если дата сегодняшняя - валидируем время, что оно не прошло
	r, err := c.reminderSrv.GetFromMemory(telectx.Chat().ID)
	if err != nil {
		return err
	}

	if r.Type == model.DateType {
		loc, err := c.userSrv.GetLocation(ctx, telectx.Chat().ID)
		if err != nil {
			return err
		}

		todayWithTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), 0, 0, time.Local)

		userDateWithTime, err := time.Parse("02.01.2006 15:04", fmt.Sprintf("%s %s", r.Date, telectx.Message().Text))
		if err != nil {
			return err
		}

		userDateTimezone := time.Date(userDateWithTime.Year(), userDateWithTime.Month(), userDateWithTime.Day(), userDateWithTime.Hour(),
			userDateWithTime.Minute(), userDateWithTime.Second(), 0, loc)

		todayYear, todayMonth, todayDay := todayWithTime.Date()
		userYear, userMonth, userDay := userDateTimezone.Date()

		todayDate := time.Date(todayYear, todayMonth, todayDay, 0, 0, 0, 0, loc)
		userDate := time.Date(userYear, userMonth, userDay, 0, 0, 0, 0, loc)

		// проверяем, сегодняшняя ли дата
		if todayDate.Equal(userDate) {
			err = c.reminderSrv.ValidateTime(loc, userDateTimezone)
			if err != nil {
				return err
			}
		}

	}

	err = c.reminderSrv.SaveTime(telectx.Chat().ID, telectx.Message().Text)
	if err != nil {
		return err
	}

	// сохраняем напоминание
	return c.saveReminder(ctx, telectx)
}
