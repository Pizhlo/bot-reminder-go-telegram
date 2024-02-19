package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type daysDuration struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
}

func newDaysDurationState(controller *controller.Controller, FSM *FSM) *daysDuration {
	return &daysDuration{controller, FSM, logger.New(), "days duration"}
}

func (n *daysDuration) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	err := n.controller.DaysDuration(ctx, telectx)
	if err != nil {
		return err
	}

	// в случае успеха меняем стейт
	n.fsm.SetState(n.fsm.ReminderTime)

	return nil
}

func (n *daysDuration) Name() string {
	return n.name
}
