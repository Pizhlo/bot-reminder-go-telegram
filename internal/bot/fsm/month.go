package fsm

import (
	"context"
	"errors"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние для обработки напоминаний раз в месяц. Валидирует число месяца
type month struct {
	controller *controller.Controller
	fsm        *FSM
	name       stateName
	next       state
}

func newMonthState(controller *controller.Controller, FSM *FSM) *month {
	return &month{controller, FSM, monthStateName, FSM.ReminderTime}
}

func (n *month) Handle(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	err := n.controller.DaysInMonthDuration(ctx, telectx)
	if err != nil {
		if errors.Is(err, api_errors.ErrInvalidDays) {
			return telectx.EditOrSend(messages.InvalidDaysInMonthMessage, view.BackToReminderMenuBtns())
		}
		return err
	}

	// в случае успеха меняем стейт
	n.fsm.SetNext()

	//return n.fsm.Handle(ctx, telectx)

	return nil
}

func (n *month) Name() string {
	return string(n.name)
}

func (n *month) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.defaultState
}
