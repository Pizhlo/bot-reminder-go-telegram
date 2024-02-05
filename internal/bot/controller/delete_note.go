package controller

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	tele "gopkg.in/telebot.v3"
)

// DeleteNoteByID удаляет одну заметку пользователя по ID заметки
func (c *Controller) DeleteNoteByID(ctx context.Context, telectx tele.Context) error {
	text := telectx.Message().Text

	c.logger.Debugf("Controller: handling %s command\n", text)

	// достаем айди заметки из сообщения (напр. /del6 -> 6)
	noteIDString, found := strings.CutPrefix(text, "/del")
	if !found {
		err := fmt.Errorf("error in controller.DeleteNoteByID(): not found suffix `/del` in message text: %s", text)
		c.HandleError(telectx, err)
		return err
	}

	// приводим к типу int
	noteID, err := strconv.Atoi(noteIDString)
	if err != nil {
		err := fmt.Errorf("error while convertion string %s to int while handling command %s: %w", noteIDString, text, err)
		c.HandleError(telectx, err)
		return err
	}

	// отправляем на удаление
	err = c.noteSrv.DeleteByID(ctx, telectx.Chat().ID, noteID)
	if err != nil {
		if errors.Is(err, api_errors.ErrNotesNotFound) {
			msg := fmt.Sprintf(messages.NoNoteFoundMessage, noteID)
			return telectx.Send(msg)
		}

		err := fmt.Errorf("error while deleting note by ID %d: %w", noteID, err)
		c.HandleError(telectx, err)

		return err
	}

	c.logger.Debugf("Controller: successfully deleted user's note by ID %d\n", noteID)

	msg := fmt.Sprintf(messages.NoteDeletedSuccessMessage, noteID)
	return telectx.Send(msg, &tele.SendOptions{
		ParseMode: htmlParseMode,
	})
}
