package controller

import (
	"context"
	"fmt"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"gopkg.in/telebot.v3"
)

const (
	minutesDate = "minutes"
	hoursDate   = "hours"
)

// EverydayReminder обрабатывает кнопку "ежедневное напоминание"
func (c *Controller) EverydayReminder(ctx context.Context, telectx telebot.Context) error {
	// сохраняем тип напоминания - "everyday"
	err := c.reminderSrv.SaveType(telectx.Chat().ID, model.EverydayType)
	if err != nil {
		return err
	}

	// если сообщение прислал пользователь - это новое название напоминания
	if !telectx.Message().Sender.IsBot {
		// сохраняем новое название напоминания, если пользователь прислал повторно
		c.reminderSrv.SaveName(telectx.Chat().ID, telectx.Message().Text)
	}

	// сохраняем в качестве даты также строку "everyday"
	err = c.reminderSrv.SaveDate(telectx.Chat().ID, string(model.EverydayType))
	if err != nil {
		return err
	}

	return telectx.EditOrSend(messages.ReminderTimeMessage, &telebot.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.BackToReminderMenuBtns(),
	})
}

// EverydayReminder обрабатывает кнопку "несколько раз в день"
func (c *Controller) SeveralTimesADayReminder(ctx context.Context, telectx telebot.Context) error {
	// сохраняем тип напоминания - "several times a day"
	err := c.reminderSrv.SaveType(telectx.Chat().ID, model.SeveralTimesDayType)
	if err != nil {
		return err
	}

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
	err := c.reminderSrv.SaveDate(telectx.Chat().ID, minutesDate)
	if err != nil {
		return err
	}

	return telectx.EditOrSend(messages.MinutesDurationMessage, view.BackToReminderMenuBtns())
}

// OnceInMinutes обрабатывает кнопку "раз в несколько часов"
func (c *Controller) OnceInHours(ctx context.Context, telectx telebot.Context) error {
	err := c.reminderSrv.SaveDate(telectx.Chat().ID, hoursDate)
	if err != nil {
		return err
	}

	return telectx.EditOrSend(messages.HoursDurationMessage, view.BackToReminderMenuBtns())
}

// OnceInMinutes обрабатывает кнопку "раз в неделю"
func (c *Controller) EveryWeek(ctx context.Context, telectx telebot.Context) error {
	// сохраняем тип напоминания - "everyweek"
	err := c.reminderSrv.SaveType(telectx.Chat().ID, model.EveryWeekType)
	if err != nil {
		return err
	}

	var msg string

	// если сообщение прислал пользователь - это новое название напоминания
	if !telectx.Message().Sender.IsBot {
		// сохраняем новое название напоминания, если пользователь прислал повторно
		c.reminderSrv.SaveName(telectx.Chat().ID, telectx.Message().Text)
	}

	r, err := c.reminderSrv.GetFromMemory(telectx.Chat().ID)
	if err != nil {
		msg = messages.ChooseWeekDayMessage
	}

	msg = fmt.Sprintf(messages.ChooseWeekDayMessage, r.Name, "раз в неделю")

	return telectx.EditOrSend(msg, &telebot.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.WeekMenu(),
	})
}

// WeekDay обрабатывает выбор пользователя определенного дня недели - понедельник, вторник, и т.д.
// Это приходит с нажатий на клавиатуру view.WeekMenu()
func (c *Controller) WeekDay(ctx context.Context, telectx telebot.Context) error {
	err := c.reminderSrv.SaveDate(telectx.Chat().ID, telectx.Callback().Unique)
	if err != nil {
		return err
	}

	return telectx.EditOrSend(messages.ReminderTimeMessage, view.BackToReminderMenuBtns())
}

// SeveralDays обрабатывает кнопку "раз в несколько дней"
func (c *Controller) SeveralDays(ctx context.Context, telectx telebot.Context) error {
	// сохраняем тип напоминания - "several_days"
	err := c.reminderSrv.SaveType(telectx.Chat().ID, model.SeveralDaysType)
	if err != nil {
		return err
	}

	var msg string

	// если сообщение прислал пользователь - это новое название напоминания
	if !telectx.Message().Sender.IsBot {
		// сохраняем новое название напоминания, если пользователь прислал повторно
		c.reminderSrv.SaveName(telectx.Chat().ID, telectx.Message().Text)
	}

	r, err := c.reminderSrv.GetFromMemory(telectx.Chat().ID)
	if err != nil {
		msg = messages.DaysDurationMessage
	}

	msg = fmt.Sprintf(messages.DaysDurationMessage, r.Name, "раз в несколько дней")

	return telectx.EditOrSend(msg, &telebot.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.BackToReminderMenuBtns(),
	})
}

// Month обрабатывает кнопку "раз в месяц"
func (c *Controller) Month(ctx context.Context, telectx telebot.Context) error {
	// сохраняем тип напоминания - "once_month"
	err := c.reminderSrv.SaveType(telectx.Chat().ID, model.OnceMonthType)
	if err != nil {
		return err
	}

	var msg string

	// если сообщение прислал пользователь - это новое название напоминания
	if !telectx.Message().Sender.IsBot {
		// сохраняем новое название напоминания, если пользователь прислал повторно
		c.reminderSrv.SaveName(telectx.Chat().ID, telectx.Message().Text)
	}

	r, err := c.reminderSrv.GetFromMemory(telectx.Chat().ID)
	if err != nil {
		msg = messages.MonthDayMessage
	}

	msg = fmt.Sprintf(messages.MonthDayMessage, r.Name, "раз в месяц")

	return telectx.EditOrSend(msg, &telebot.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.BackToReminderMenuBtns(),
	})
}

// Year обрабатывает кнопку "раз в год"
func (c *Controller) Year(ctx context.Context, telectx telebot.Context) error {
	// сохраняем тип напоминания - "once_year"
	err := c.reminderSrv.SaveType(telectx.Chat().ID, model.OnceYearType)
	if err != nil {
		return err
	}

	var msg string

	// если сообщение прислал пользователь - это новое название напоминания
	if !telectx.Message().Sender.IsBot {
		// сохраняем новое название напоминания, если пользователь прислал повторно
		c.reminderSrv.SaveName(telectx.Chat().ID, telectx.Message().Text)
	}

	r, err := c.reminderSrv.GetFromMemory(telectx.Chat().ID)
	if err != nil {
		msg = messages.CalendarMessage
	}

	msg = fmt.Sprintf(messages.CalendarMessage, r.Name, "раз в год")

	kb := c.reminderSrv.Calendar(telectx.Chat().ID)

	btns := c.reminderSrv.DaysBtns(telectx.Chat().ID)

	for _, btn := range btns {
		c.bot.Handle(&btn, func(telectx telebot.Context) error {
			return c.ProcessDate(ctx, telectx)
		})
	}

	return telectx.EditOrSend(msg, &telebot.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: kb,
	})
}

func (c *Controller) DaysBtns(ctx context.Context, telectx telebot.Context) []telebot.Btn {
	return c.reminderSrv.DaysBtns(telectx.Chat().ID)
}
