package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type year struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
}

func newYearState(controller *controller.Controller, FSM *FSM) *year {
	return &year{controller, FSM, logger.New(), "every year"}
}

func (n *year) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	return n.controller.Year(ctx, telectx)
}

func (n *year) Name() string {
	return n.name
}
