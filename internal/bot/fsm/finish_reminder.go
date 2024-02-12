package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type finishReminder struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
}

func newFinishReminderState(controller *controller.Controller, FSM *FSM) *finishReminder {
	return &finishReminder{controller, FSM, logger.New(), "finish reminder"}
}

func (n *finishReminder) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	n.fsm.SetState(n.fsm.DefaultState)

	return n.controller.FinishReminder(ctx, telectx)
}

func (n *finishReminder) Name() string {
	return n.name
}
