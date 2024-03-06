package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние для обработки одноразового напоминания (дата выбирается в календаре)
type onceReminder struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
	next       state
}

func newOnceReminderState(controller *controller.Controller, FSM *FSM) *onceReminder {
	return &onceReminder{controller, FSM, logger.New(), "date reminder", newSelectedDateOnceState(controller, FSM)}
}

func (n *onceReminder) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	return n.controller.Date(ctx, telectx)
}

func (n *onceReminder) Name() string {
	return n.name
}

func (n *onceReminder) Next() {
	if n.next != nil {
		n.fsm.SetState(n.next)
	} else {
		n.fsm.SetState(n.fsm.DefaultState)
	}
}

// Next state

// Состояние, когда пользователь выбрал дату
type selectedDateOnce struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
	next       state
}

func newSelectedDateOnceState(controller *controller.Controller, FSM *FSM) *selectedDateOnce {
	s := &selectedDateOnce{controller, FSM, logger.New(), "date reminder: selected date", nil}

	// для случаев, когда пользователь несколько раз ввел не ту дату, проверка не прошла (дата уже прошла)
	s.next = s
	return s
}

func (n *selectedDateOnce) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	err := n.controller.ProcessDate(ctx, telectx)
	if err != nil {
		return err
	}

	return nil
}

func (n *selectedDateOnce) Name() string {
	return n.name
}

func (n *selectedDateOnce) Next() {
	if n.next != nil {
		n.fsm.SetState(n.next)
	} else {
		n.fsm.SetState(n.fsm.DefaultState)
	}
}
