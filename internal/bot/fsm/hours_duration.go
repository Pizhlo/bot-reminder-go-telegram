package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние для обработки напоминаний раз в неск. часов
type hoursDuration struct {
	controller *controller.Controller
	fsm        *FSM
	name       stateName
	next       state
}

func newHoursDurationState(controller *controller.Controller, FSM *FSM) *hoursDuration {
	return &hoursDuration{controller, FSM, hoursStateName, nil}
}

func (n *hoursDuration) Handle(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	return n.controller.HoursDuration(ctx, telectx)
}

func (n *hoursDuration) Name() string {
	return string(n.name)
}

func (n *hoursDuration) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.defaultState
}
