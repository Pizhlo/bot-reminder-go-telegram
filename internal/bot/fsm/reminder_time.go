package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние для обработки времени напоминания
type reminderTime struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
	next       state
}

func newReminderTimeState(controller *controller.Controller, FSM *FSM) *reminderTime {
	return &reminderTime{controller, FSM, logger.New(), "reminder time", FSM.DefaultState}
}

func (n *reminderTime) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	return n.controller.ReminderTime(ctx, telectx)
}

func (n *reminderTime) Name() string {
	return n.name
}

func (n *reminderTime) Next() {
	if n.next != nil {
		n.fsm.SetState(n.next)
	} else {
		n.fsm.SetState(n.fsm.DefaultState)
	}
}
