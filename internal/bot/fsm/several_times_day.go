package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние для обработки напоминаний несколько раз в день (раз в неск. секунд / минут)
type severalTimes struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
	next       state
}

func newSeveralTimesState(controller *controller.Controller, FSM *FSM) *severalTimes {
	return &severalTimes{controller, FSM, logger.New(), "several times a day", nil}
}

func (n *severalTimes) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	return n.controller.SeveralTimesADayReminder(ctx, telectx)
}

func (n *severalTimes) Name() string {
	return n.name
}

func (n *severalTimes) Next() {
	if n.next != nil {
		n.fsm.SetState(n.next)
	} else {
		n.fsm.SetState(n.fsm.DefaultState)
	}
}
