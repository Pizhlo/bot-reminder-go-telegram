package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type severalDays struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
}

func newSeveralDaysState(controller *controller.Controller, FSM *FSM) *severalDays {
	return &severalDays{controller, FSM, logger.New(), "once in several days"}
}

func (n *severalDays) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	return n.controller.SeveralDays(ctx, telectx)
}

func (n *severalDays) Name() string {
	return n.name
}
