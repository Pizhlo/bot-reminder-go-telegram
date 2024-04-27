package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние для обработки напоминаний в указанные времена (10:30, 12:30, 15:00)
type timesReminder struct {
	controller *controller.Controller
	fsm        *FSM
	name       stateName
	next       state
}

func newTimesReminderState(controller *controller.Controller, FSM *FSM) *timesReminder {
	return &timesReminder{controller, FSM, timesReminderName, nil}
}

func (n *timesReminder) Handle(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	return n.controller.Times(ctx, telectx)
}

func (n *timesReminder) Name() string {
	return string(n.name)
}

func (n *timesReminder) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.defaultState
}
