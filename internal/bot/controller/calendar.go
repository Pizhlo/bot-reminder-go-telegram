package controller

import (
	"context"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
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

func (c *Controller) ProcessDate(ctx context.Context, telectx tele.Context) error {
	r, err := c.reminderSrv.GetFromMemory(telectx.Chat().ID)
	if err != nil {
		return err
	}

	switch r.Type {
	case model.OnceYearType:
		err := c.reminderSrv.SaveCalendarDate(telectx.Chat().ID, telectx.Callback().Unique)
		if err != nil {
			return err
		}
	}

	return telectx.EditOrSend(messages.ReminderTimeMessage, view.BackToMenuBtn())
}
