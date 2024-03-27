package fsm

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние показа напоминаний
type listReminder struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
	next       state
}

// префикс для удаления напоминания. Например: /dr1
const deleteReminderPrefix = "/dr"

func newListReminderState(c *controller.Controller, fsm *FSM) *listReminder {
	return &listReminder{controller: c, fsm: fsm, logger: logger.New(), name: listReminderName, next: nil}
}

func (n *listReminder) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)
	msg := telectx.Message().Text

	if !strings.HasPrefix(msg, deleteReminderPrefix) {
		return n.controller.CreateNote(ctx, telectx)
	} else {
		before, found := strings.CutPrefix(msg, deleteReminderPrefix)
		if !found {
			return fmt.Errorf("list reminder state: not found prefix %s in message: %s", deleteReminderPrefix, msg)
		}

		reminderID, err := strconv.Atoi(before)
		if err != nil {
			err := fmt.Errorf("error while convertion string %s to int while handling command %s: %w", before, msg, err)
			return err
		}

		if err := n.controller.DeleteReminderByViewID(ctx, telectx, reminderID); err != nil {
			return err
		}

		return telectx.EditOrSend(fmt.Sprintf(messages.ReminderDeletedByIDMessage, reminderID), &tele.SendOptions{
			ParseMode:   "html",
			ReplyMarkup: view.BackToReminderMenuBtns(),
		})
	}

	//return n.controller.ListNotes(ctx, telectx)
}

func (n *listReminder) Name() string {
	return n.name
}

func (n *listReminder) Next() state {
	if n.next != nil {
		return n.next
	}

	return n.fsm.defaultState
}
