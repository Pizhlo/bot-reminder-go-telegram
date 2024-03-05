package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние для обработки напоминаний раз в неделю
type weekDay struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
	next       state
}

func newWeekDayState(controller *controller.Controller, FSM *FSM) *weekDay {
	return &weekDay{controller, FSM, logger.New(), "every week", nil}
}

func (n *weekDay) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	return n.controller.EveryWeek(ctx, telectx)
}

func (n *weekDay) Name() string {
	return n.name
}

func (n *weekDay) Next() {
	if n.next != nil {
		n.fsm.SetState(n.next)
	} else {
		n.fsm.SetState(n.fsm.DefaultState)
	}
}
