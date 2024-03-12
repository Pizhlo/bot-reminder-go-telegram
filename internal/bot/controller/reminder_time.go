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

		todayWithTime := time.Now().In(loc)

		userDateWithTime, err := time.Parse("02.01.2006 15:04", fmt.Sprintf("%s %s", r.Date, telectx.Message().Text))
		if err != nil {
			return err
		}

		todayYear, todayMonth, todayDay := todayWithTime.Date()
		userYear, userMonth, userDay := userDateWithTime.Date()

		todayDate := time.Date(todayYear, todayMonth, todayDay, 0, 0, 0, 0, loc)
		userDate := time.Date(userYear, userMonth, userDay, 0, 0, 0, 0, loc)

		// проверяем, сегодняшняя ли дата
		if todayDate.Equal(userDate) {
			err = c.reminderSrv.ValidateTime(loc, userDateWithTime)
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
