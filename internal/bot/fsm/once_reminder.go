package fsm

import (
	"context"
	"errors"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние для обработки одноразового напоминания (дата выбирается в календаре)
type dateReminder struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       stateName
	next       state
}

func newDateReminderState(controller *controller.Controller, FSM *FSM) *dateReminder {
	return &dateReminder{controller, FSM, logger.New(), dateReminderName, FSM.ReminderTime}
}

func (n *dateReminder) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	err := n.controller.ProcessDate(ctx, telectx)
	if err != nil {
		if errors.Is(err, api_errors.ErrInvalidDate) {
			return n.controller.InvalidDate(ctx, telectx)
		}
	}

	n.fsm.SetNext()

	return nil
}

func (n *dateReminder) Name() string {
	return string(n.name)
}

func (n *dateReminder) Next() state {
	if n.next != nil {
		return n.next
	}

	return n.fsm.defaultState
}

// Next state

// Состояние, когда пользователь выбрал дату
// type selectedDateOnce struct {
// 	controller *controller.Controller
// 	fsm        *FSM
// 	logger     *logrus.Logger
// 	name       string
// 	next       state
// }

// func newSelectedDateOnceState(controller *controller.Controller, FSM *FSM) *selectedDateOnce {
// 	s := &selectedDateOnce{controller, FSM, logger.New(), "date reminder: selected date", FSM.ReminderTime}

// 	// для случаев, когда пользователь несколько раз ввел не ту дату, проверка не прошла (дата уже прошла)
// 	s.next = s
// 	return s
// }

// func (n *selectedDateOnce) Handle(ctx context.Context, telectx tele.Context) error {
// 	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

// 	err := n.controller.ProcessDate(ctx, telectx)
// 	if err != nil {
// 		if errors.Is(err, api_errors.ErrInvalidDate) {
// 			return n.controller.InvalidDate(ctx, telectx)
// 		}
// 	}

// 	n.fsm.SetNext()

// 	return nil
// }

// func (n *selectedDateOnce) Name() string {
// 	return string(n.name)
// }

// func (n *selectedDateOnce) Next() state {
// 	if n.next != nil {
// 		return n.next
// 	}

// 	return n.fsm.DefaultState
// }
