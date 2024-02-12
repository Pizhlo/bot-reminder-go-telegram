package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type reminderName struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
}

func newReminderNameState(controller *controller.Controller, FSM *FSM) *reminderName {
	return &reminderName{controller, FSM, logger.New(), "reminder name"}
}

func (n *reminderName) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)
	return n.controller.ReminderName(ctx, telectx)
}

func (n *reminderName) Name() string {
	return n.name
}
