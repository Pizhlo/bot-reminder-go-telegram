package fsm

import (
	"context"
	"errors"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние для обработки числа месяца, в которое присылать уведомление.
// Тип: напоминание раз в месяц
type daysDuration struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
	next       state
}

func newDaysDurationState(controller *controller.Controller, FSM *FSM) *daysDuration {
	return &daysDuration{controller, FSM, logger.New(), "days duration", nil}
}

func (n *daysDuration) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	err := n.controller.DaysDuration(ctx, telectx)
	if err != nil {
		if errors.Is(err, api_errors.ErrInvalidDays) {
			return telectx.EditOrSend(messages.InvalidDaysMessage, view.BackToReminderMenuBtns())
		}
		return err
	}

	// в случае успеха меняем стейт
	n.fsm.SetState(n.fsm.ReminderTime)

	return nil
}

func (n *daysDuration) Name() string {
	return n.name
}

func (n *daysDuration) Next() {
	if n.next != nil {
		n.fsm.SetState(n.next)
	} else {
		n.fsm.SetState(n.fsm.DefaultState)
	}
}
