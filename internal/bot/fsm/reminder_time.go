package fsm

import (
	"context"
	"errors"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние для обработки времени напоминания
type reminderTime struct {
	controller *controller.Controller
	fsm        *FSM
	name       stateName
	next       state
}

func newReminderTimeState(controller *controller.Controller, FSM *FSM) *reminderTime {
	return &reminderTime{controller, FSM, reminderTimeState, FSM.defaultState}
}

func (n *reminderTime) Handle(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	err := n.controller.ReminderTime(ctx, telectx)
	if err != nil {
		if errors.Is(err, api_errors.ErrInvalidTime) {
			return telectx.EditOrSend(messages.InvalidTimeMessage, view.BackToMenuBtn())
		}

		if errors.Is(err, api_errors.ErrTimeInPast) {
			return telectx.EditOrSend(messages.TimeInPastMessage, view.BackToMenuBtn())
		}

		return err
	}

	n.fsm.SetNext()

	return nil
}

func (n *reminderTime) Name() string {
	return string(n.name)
}

func (n *reminderTime) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.defaultState
}
