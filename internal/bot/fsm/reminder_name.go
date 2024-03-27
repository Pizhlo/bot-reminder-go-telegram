package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние для обработки названия напоминания
type reminderName struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       stateName
	next       state
}

func newReminderNameState(controller *controller.Controller, FSM *FSM) *reminderName {
	return &reminderName{controller, FSM, logger.New(), reminderNameState, nil}
}

func (n *reminderName) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)
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
