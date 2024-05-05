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

// Состояние для обработки напоминаний раз в несколько дней. Валидирует число.
// Тип: напоминание раз в несколько дней
type daysDuration struct {
	controller *controller.Controller
	fsm        *FSM
	name       stateName
	next       state
}

func newDaysDurationState(controller *controller.Controller, FSM *FSM) *daysDuration {
	return &daysDuration{controller, FSM, daysDurationName, FSM.ReminderTime}
}

func (n *daysDuration) Handle(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	err := n.controller.DaysDuration(ctx, telectx)
	if err != nil {
		if errors.Is(err, api_errors.ErrInvalidDays) {
			return telectx.EditOrSend(messages.InvalidDaysMessage, view.BackToReminderMenuBtns())
		}
		return err
	}

	// в случае успеха меняем стейт
	n.fsm.SetNext()

	return nil
}

func (n *daysDuration) Name() string {
	return string(n.name)
}

func (n *daysDuration) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.defaultState
}
