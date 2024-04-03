package controller

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

const deleteNotePrefix = "/dn"

// DeleteNoteByID удаляет одну заметку пользователя по ID заметки
func (c *Controller) DeleteNoteByID(ctx context.Context, telectx tele.Context) error {
	text := telectx.Message().Text

	logrus.Debugf("Controller: handling %s command\n", text)

	// достаем айди заметки из сообщения (напр. /del6 -> 6)
	noteIDString, found := strings.CutPrefix(text, deleteNotePrefix)
	if !found {
		err := fmt.Errorf("error in controller.DeleteNoteByID(): not found suffix %s in message text: %s", deleteNotePrefix, text)

		return err
	}

	// приводим к типу int
	noteID, err := strconv.Atoi(noteIDString)
	if err != nil {
		err := fmt.Errorf("error while convertion string %s to int while handling command %s: %w", noteIDString, text, err)

		return err
	}

	// отправляем на удаление
	err = c.noteSrv.DeleteByID(ctx, telectx.Chat().ID, noteID)
	if err != nil {
		if errors.Is(err, api_errors.ErrNotesNotFound) {
			msg := fmt.Sprintf(messages.NoNoteFoundByNumberMessage, noteID)
			return telectx.EditOrSend(msg, view.BackToMenuBtn())
		}

		err := fmt.Errorf("error while deleting note by ID %d: %w", noteID, err)
		return err
	}

	logrus.Debugf("Controller: successfully deleted user's note by ID %d\n", noteID)

	msg := fmt.Sprintf(messages.NoteDeletedSuccessMessage, noteID)
	return telectx.EditOrSend(msg, &tele.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.NotesAndMenuBtns(),
	})
}
