package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние для обработки напоминаний раз в год
type year struct {
	controller *controller.Controller
	fsm        *FSM
	name       stateName
	next       state
}

func newYearState(controller *controller.Controller, FSM *FSM) *year {
	return &year{controller, FSM, yearReminderState, FSM.ReminderTime}
}

func (n *year) Handle(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	err := n.controller.SaveCalendarDate(ctx, telectx)
	if err != nil {
		return err
	}

	n.fsm.SetNext()

	return nil
}

func (n *year) Name() string {
	return string(n.name)
}

func (n *year) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.defaultState
}

// type selectedDateYear struct {
// 	controller *controller.Controller
// 	fsm        *FSM
//
// 	name       string
// 	next       state
// }

// func newSelectedDateYear(controller *controller.Controller, FSM *FSM) *selectedDateYear {
// 	return &selectedDateYear{controller: controller, fsm: FSM,  name: "year: selected date", next: nil}
// }

// func (n *selectedDateYear) Handle(ctx context.Context, telectx tele.Context) error {
// 	logrus.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

// 	return n.controller.SaveCalendarDate(ctx, telectx)
// }

// func (n *selectedDateYear) Name() string {
// 	return string(n.name)
// }

// func (n *selectedDateYear) Next() state {
// 	if n.next != nil {
// 		return n.next
// 	}
// 	return n.fsm.DefaultState
// }
