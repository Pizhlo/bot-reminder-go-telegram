package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	tele "gopkg.in/telebot.v3"
)

type start struct {
	controller *controller.Controller
	fsm        *FSM
}

func newStartState(FSM *FSM, controller *controller.Controller) *start {
	return &start{controller, FSM}
}

func (n *start) Handle(ctx context.Context, telectx tele.Context) error {
	return telectx.Send(messages.StartMessageLocation)
}
