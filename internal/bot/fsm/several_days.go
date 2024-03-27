package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние для обработки напоминаний раз в неск. дней
type severalDays struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
	next       state
}

func newSeveralDaysState(controller *controller.Controller, FSM *FSM) *severalDays {
	return &severalDays{controller, FSM, logger.New(), severalDaysState, nil}
}

func (n *severalDays) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	return n.controller.SeveralDays(ctx, telectx)
}

func (n *severalDays) Name() string {
	return n.name
}

func (n *severalDays) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.defaultState
}
