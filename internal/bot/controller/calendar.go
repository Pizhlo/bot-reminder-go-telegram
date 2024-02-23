package controller

import (
	"context"

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
