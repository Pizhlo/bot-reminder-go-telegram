package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние для обработки названия напоминания
type reminderName struct {
	controller *controller.Controller
	fsm        *FSM
	name       stateName
	next       state
}

func newReminderNameState(controller *controller.Controller, FSM *FSM) *reminderName {
	return &reminderName{controller, FSM, reminderNameState, nil}
}

func (n *reminderName) Handle(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)
	return n.controller.ReminderName(ctx, telectx)
}

func (n *reminderName) Name() string {
	return string(n.name)
}

func (n *reminderName) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.defaultState
}
