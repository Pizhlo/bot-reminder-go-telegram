package controller

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	tele "gopkg.in/telebot.v3"
)

// SearchNoteByTwoDates обрабатывает кнопку "поиск по диапазону дат". Отправляет пользователю клавиатуру с календарем
func (c *Controller) SearchNoteByTwoDates(ctx context.Context, telectx tele.Context) error {
	return telectx.EditOrSend(messages.SearchByTwoDatesFirstDateMessage, c.noteSrv.Calendar(telectx.Chat().ID))
}

// SearchNoteByTwoDatesFirstDate производит поиск заметок по двум датам. Обрабатывает первую дату
func (c *Controller) SearchNoteByTwoDatesFirstDate(ctx context.Context, telectx tele.Context) error {
	day, err := strconv.Atoi(telectx.Callback().Unique)
	if err != nil {
		return fmt.Errorf("error while converting string %s to type int: %w", telectx.Callback().Unique, err)
	}

	loc, err := c.userSrv.GetLocation(ctx, telectx.Chat().ID)
	if err != nil {
		return fmt.Errorf("error while getting user's timezone: %w", err)
	}

	firstDate := time.Date(c.noteSrv.CurYear(telectx.Chat().ID), c.noteSrv.CurMonth(telectx.Chat().ID), day, 0, 0, 0, 0, loc)

	today := time.Now()

	if firstDate.After(today) {
		return api_errors.ErrFirstDayFuture
	}

	c.noteSrv.SaveFirstDate(telectx.Chat().ID, firstDate)

	return telectx.EditOrSend(messages.SearchByTwoDatesSecondDateMessage, c.noteSrv.Calendar(telectx.Chat().ID))
}

// SearchNoteByTwoDatesSecondDate производит поиск заметок по двум датам. Обрабатывает вторую дату
func (c *Controller) SearchNoteByTwoDatesSecondDate(ctx context.Context, telectx tele.Context) error {
	day, err := strconv.Atoi(telectx.Callback().Unique)
	if err != nil {
		return fmt.Errorf("error while converting string %s to type int: %w", telectx.Callback().Unique, err)
	}

	loc, err := c.userSrv.GetLocation(ctx, telectx.Chat().ID)
	if err != nil {
		return fmt.Errorf("error while getting user's timezone: %w", err)
	}

	secondDate := time.Date(c.noteSrv.CurYear(telectx.Chat().ID), c.noteSrv.CurMonth(telectx.Chat().ID), day, 0, 0, 0, 0, loc)

	err = c.noteSrv.ValidateSearchDate(telectx.Chat().ID, secondDate)
	if err != nil {
		return err
	}

	err = c.noteSrv.SaveSecondDate(telectx.Chat().ID, secondDate)
	if err != nil {
		return err
	}

	search, err := c.noteSrv.GetSearchNote(telectx.Chat().ID)
	if err != nil {
		return err
	}

	msg, kb, err := c.noteSrv.SearchByTwoDates(ctx, telectx.Chat().ID)
	if err != nil {
		if errors.Is(err, api_errors.ErrNotesNotFound) {
			msg := fmt.Sprintf(messages.NoNotesFoundByTwoDatesMessage, search.FirstDate.Format("02.01.2006"), search.SecondDate.Format("02.01.2006"))
			return telectx.EditOrSend(msg, view.BackToMenuAndNotesBtn())
		}

		return err
	}

	return telectx.EditOrSend(msg, &tele.SendOptions{
		ReplyMarkup: kb,
		ParseMode:   htmlParseMode,
	})
}

func (c *Controller) SecondDateBeforeFirst(ctx context.Context, telectx tele.Context) error {
	return telectx.EditOrSend(messages.FirstDateBeforeSecondMessage, &tele.SendOptions{
		ReplyMarkup: c.noteSrv.Calendar(telectx.Chat().ID),
		ParseMode:   htmlParseMode,
	})
}

func (c *Controller) SecondDateInFuture(ctx context.Context, telectx tele.Context) error {
	return telectx.EditOrSend(messages.SecondDateInFutureMessage, &tele.SendOptions{
		ReplyMarkup: c.noteSrv.Calendar(telectx.Chat().ID),
		ParseMode:   htmlParseMode,
	})
}

func (c *Controller) FirstDateInFuture(ctx context.Context, telectx tele.Context) error {
	return telectx.EditOrSend(messages.FirstDateInFutureMessage, &tele.SendOptions{
		ReplyMarkup: c.noteSrv.Calendar(telectx.Chat().ID),
		ParseMode:   htmlParseMode,
	})
}
