package controller

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	tele "gopkg.in/telebot.v3"
)

// SearchNoteByOnedate обрабатывает кнопку "поиск по одной дате". Отправляет пользователю клавиатуру с календарем
func (c *Controller) SearchNoteByOnedate(ctx context.Context, telectx tele.Context) error {
	return telectx.EditOrSend(messages.SearchOneDateMessage, c.noteSrv.Calendar(telectx.Chat().ID))
}

// SearchNoteByOnedate производит поиск заметок по выбранной дате
func (c *Controller) SearchNoteBySelectedDate(ctx context.Context, telectx tele.Context) error {
	day, err := strconv.Atoi(telectx.Callback().Unique)
	if err != nil {
		return fmt.Errorf("error while converting string %s to type int: %w", telectx.Callback().Unique, err)
	}

	month := c.noteSrv.CurMonth(telectx.Chat().ID)
	year := c.noteSrv.CurYear(telectx.Chat().ID)

	tz, err := c.userSrv.GetTimezone(ctx, telectx.Chat().ID)
	if err != nil {
		return fmt.Errorf("error while getting user's timezone: %w", err)
	}

	loc, err := time.LoadLocation(tz.Name)
	if err != nil {
		return fmt.Errorf("error while loading user's timezone: %w", err)
	}

	search := model.SearchByOneDate{
		TgID: telectx.Chat().ID,
		Date: time.Date(year, month, day, 0, 0, 0, 0, loc),
	}

	notes, kb, err := c.noteSrv.SearchOneDate(ctx, search)
	if err != nil {
		if errors.Is(err, api_errors.ErrNotesNotFound) {
			msg := fmt.Sprintf(messages.NoNotesFoundByDateMessage, search.Date.Format("02.01.2006"))
			return telectx.EditOrSend(msg, view.BackToMenuAndNotesBtn())
		}

		return err
	}

	return telectx.EditOrSend(notes, &tele.SendOptions{
		ReplyMarkup: kb,
		ParseMode:   htmlParseMode,
	})
}
