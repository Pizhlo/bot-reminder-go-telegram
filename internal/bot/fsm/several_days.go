package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние для обработки напоминаний раз в неск. дней
type severalDays struct {
	controller *controller.Controller
	fsm        *FSM
	name       stateName
	next       state
}

func newSeveralDaysState(controller *controller.Controller, FSM *FSM) *severalDays {
	return &severalDays{controller, FSM, severalDaysState, nil}
}

func (n *severalDays) Handle(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	return n.controller.SeveralDays(ctx, telectx)
}

func (n *severalDays) Name() string {
	return string(n.name)
}

func (n *severalDays) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.defaultState
}
