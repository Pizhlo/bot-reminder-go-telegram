package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние для обработки напоминаний раз в неск. часов
type hoursDuration struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
	next       state
}

func newHoursDurationState(controller *controller.Controller, FSM *FSM) *hoursDuration {
	return &hoursDuration{controller, FSM, logger.New(), "hours duration", nil}
}

func (n *hoursDuration) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	return n.controller.HoursDuration(ctx, telectx)
}

func (n *hoursDuration) Name() string {
	return n.name
}

func (n *hoursDuration) Next() {
	if n.next != nil {
		n.fsm.SetState(n.next)
	} else {
		n.fsm.SetState(n.fsm.DefaultState)
	}
}
