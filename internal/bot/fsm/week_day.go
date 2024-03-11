package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние для обработки напоминаний раз в неделю
type everyWeek struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
	next       state
}

func newEveryWeekState(controller *controller.Controller, FSM *FSM) *everyWeek {
	return &everyWeek{controller, FSM, logger.New(), "every week", newWeekDayState(controller, FSM)}
}

func (n *everyWeek) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	n.fsm.SetNext()

	return n.controller.EveryWeek(ctx, telectx)
}

func (n *everyWeek) Name() string {
	return n.name
}

func (n *everyWeek) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.DefaultState
}

type weekDay struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
	next       state
}

func newWeekDayState(controller *controller.Controller, FSM *FSM) *weekDay {
	return &weekDay{controller, FSM, logger.New(), "every week", FSM.ReminderTime}
}

func (n *weekDay) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	// если пользователь нажал кнопку
	if telectx.Callback() != nil {
		n.fsm.SetNext()
	}

	return n.controller.WeekDay(ctx, telectx)
}

func (n *weekDay) Name() string {
	return n.name
}

func (n *weekDay) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.DefaultState
}
