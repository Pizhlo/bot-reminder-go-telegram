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

// Состояние для обработки времени напоминания
type reminderTime struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
	next       state
}

func newReminderTimeState(controller *controller.Controller, FSM *FSM) *reminderTime {
	return &reminderTime{controller, FSM, logger.New(), "reminder time", FSM.DefaultState}
}

func (n *reminderTime) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	err := n.controller.ReminderTime(ctx, telectx)
	if err != nil {
		if errors.Is(err, api_errors.ErrInvalidTime) {
			return telectx.EditOrSend(messages.InvalidTimeMessage, view.BackToMenuBtn())
		}
		return err
	}

	n.fsm.SetNext()

	return nil
}

func (n *reminderTime) Name() string {
	return n.name
}

func (n *reminderTime) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.DefaultState
}
