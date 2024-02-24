package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type onceReminder struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
}

func newOnceReminderState(controller *controller.Controller, FSM *FSM) *onceReminder {
	return &onceReminder{controller, FSM, logger.New(), "date reminder"}
}

func (n *onceReminder) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	return n.controller.Date(ctx, telectx)
}

func (n *onceReminder) Name() string {
	return n.name
}
