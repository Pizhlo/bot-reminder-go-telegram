package controller

import (
	"context"
	"time"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	tele "gopkg.in/telebot.v3"
)

// PrevMonth обрабатывает кнопку "предыдущий месяц" в календаре
func (c *Controller) PrevMonth(ctx context.Context, telectx tele.Context) error {
	return telectx.Edit(c.reminderSrv.PrevMonth(telectx.Chat().ID))
}

// NextMonth обрабатывает кнопку "следующий месяц" в календаре
func (c *Controller) NextMonth(ctx context.Context, telectx tele.Context) error {
	return telectx.Edit(c.reminderSrv.NextMonth(telectx.Chat().ID))
}

// PrevYear обрабатывает кнопку "предыдущий год" в календаре
func (c *Controller) PrevYear(ctx context.Context, telectx tele.Context) error {
	return telectx.Edit(c.reminderSrv.PrevYear(telectx.Chat().ID))
}

// NextYear обрабатывает кнопку "следующий год" в календаре
func (c *Controller) NextYear(ctx context.Context, telectx tele.Context) error {
	return telectx.Edit(c.reminderSrv.NextYear(telectx.Chat().ID))
}

// SaveCalendarDate сохраняет дату, указанную в календаре, без валидации (для напоминаний раз в год)
func (c *Controller) SaveCalendarDate(ctx context.Context, telectx tele.Context) error {
	err := c.reminderSrv.SaveCalendarDate(telectx.Chat().ID, telectx.Callback().Unique)
	if err != nil {
		return err
	}

	return telectx.EditOrSend(messages.ReminderTimeMessage, view.BackToMenuBtn())
}

// ProcessDate валидирует дату и сохраняет
func (c *Controller) ProcessDate(ctx context.Context, telectx tele.Context) error {
	userTz, err := c.userSrv.GetTimezone(ctx, telectx.Chat().ID)
	if err != nil {
		return err
	}

	t, err := time.LoadLocation(userTz.Name)
	if err != nil {
		return err
	}

	err = c.reminderSrv.ValidateDate(telectx.Chat().ID, telectx.Callback().Unique, t)
	if err != nil {
		return err
	}

	return c.SaveCalendarDate(ctx, telectx)
}

// InvalidDate обрабатывает неверную дату
func (c *Controller) InvalidDate(ctx context.Context, telectx tele.Context) error {
	return telectx.EditOrSend(messages.InvalidDateMessage, c.reminderSrv.Calendar(telectx.Chat().ID))
}

// SetupCalendar устанавливает месяц и год в календаре на текущие
func (c *Controller) SetupCalendar(ctx context.Context, telectx tele.Context) {
	c.reminderSrv.SetupCalendar(telectx.Chat().ID)
}
