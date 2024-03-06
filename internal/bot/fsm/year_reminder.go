package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние для обработки напоминаний раз в год
type year struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
	next       state
}

func newYearState(controller *controller.Controller, FSM *FSM) *year {
	return &year{controller, FSM, logger.New(), "every year", newSelectedDateYear(controller, FSM)}
}

func (n *year) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	return n.controller.Year(ctx, telectx)
}

func (n *year) Name() string {
	return n.name
}

func (n *year) Next() {
	if n.next != nil {
		n.fsm.SetState(n.next)
	} else {
		n.fsm.SetState(n.fsm.DefaultState)
	}
}

type selectedDateYear struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
	next       state
}

func newSelectedDateYear(controller *controller.Controller, FSM *FSM) *selectedDateYear {
	return &selectedDateYear{controller: controller, fsm: FSM, logger: logger.New(), name: "year: selected date", next: nil}
}

func (n *selectedDateYear) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	return n.controller.SaveCalendarDate(ctx, telectx)
}

func (n *selectedDateYear) Name() string {
	return n.name
}

func (n *selectedDateYear) Next() {
	if n.next != nil {
		n.fsm.SetState(n.next)
	} else {
		n.fsm.SetState(n.fsm.DefaultState)
	}
}
