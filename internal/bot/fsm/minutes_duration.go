package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние для обработки напоминаний раз в неск. секунд
type minutesDuration struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       stateName
	next       state
}

func newMinutesDurationState(controller *controller.Controller, FSM *FSM) *minutesDuration {
	return &minutesDuration{controller, FSM, logger.New(), minutesStateName, nil}
}

func (n *minutesDuration) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	return n.controller.MinutesDuration(ctx, telectx)
}

func (n *minutesDuration) Name() string {
	return string(n.name)
}

func (n *minutesDuration) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.defaultState
}
