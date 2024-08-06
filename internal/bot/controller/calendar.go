package controller

import (
	"context"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	tele "gopkg.in/telebot.v3"
)

// PrevMonth обрабатывает кнопку "предыдущий месяц" в календаре
func (c *Controller) PrevMonth(ctx context.Context, telectx tele.Context) error {
	// если используется календарь напоминаний
	if c.reminderCalendar[telectx.Chat().ID] {
		return telectx.Edit(c.reminderSrv.PrevMonth(telectx.Chat().ID))
	}

	return telectx.Edit(c.noteSrv.PrevMonth(telectx.Chat().ID))
}

// NextMonth обрабатывает кнопку "следующий месяц" в календаре
func (c *Controller) NextMonth(ctx context.Context, telectx tele.Context) error {
	if c.reminderCalendar[telectx.Chat().ID] {
		return telectx.Edit(c.reminderSrv.NextMonth(telectx.Chat().ID))
	}

	return telectx.Edit(c.noteSrv.NextMonth(telectx.Chat().ID))
}

// PrevYear обрабатывает кнопку "предыдущий год" в календаре
func (c *Controller) PrevYear(ctx context.Context, telectx tele.Context) error {
	if c.reminderCalendar[telectx.Chat().ID] {
		return telectx.Edit(c.reminderSrv.PrevYear(telectx.Chat().ID))
	}

	return telectx.Edit(c.noteSrv.PrevYear(telectx.Chat().ID))
}

// NextYear обрабатывает кнопку "следующий год" в календаре
func (c *Controller) NextYear(ctx context.Context, telectx tele.Context) error {
	if c.reminderCalendar[telectx.Chat().ID] {
		return telectx.Edit(c.reminderSrv.NextYear(telectx.Chat().ID))
	}

	return telectx.Edit(c.noteSrv.NextYear(telectx.Chat().ID))
}

func (c *Controller) SetReminderCalendar(userID int64) {
	c.reminderCalendar[userID] = true
}

func (c *Controller) SetNoteCalendar(userID int64) {
	c.noteCalendar[userID] = true
}

func (c *Controller) ResetCalendars(userID int64) {
	c.noteCalendar[userID] = false
	c.reminderCalendar[userID] = false
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
	t, err := c.userSrv.GetLocation(ctx, telectx.Chat().ID)
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

// SetupReminderCalendar устанавливает месяц и год в календаре на текущие для напоминаний
func (c *Controller) SetupReminderCalendar(ctx context.Context, telectx tele.Context) {
	c.reminderSrv.SetupCalendar(telectx.Chat().ID)
}

// SetupNoteCalendar устанавливает месяц и год в календаре на текущие для заметок
func (c *Controller) SetupNoteCalendar(ctx context.Context, telectx tele.Context) {
	c.noteSrv.SetupCalendar(telectx.Chat().ID)
}

// ListMonths обрабатывает нажатие на кнопку с названием месяца, открывая список всех месяцев
func (c *Controller) ListMonths(ctx context.Context, telectx tele.Context) error {
	return telectx.Edit("text", view.ListMonthsKb())
}
