package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	tele "gopkg.in/telebot.v3"
)

// дефолтное состояние бота
type defaultState struct {
	fsm        *FSM
	controller *controller.Controller
	start      state
}

const startCommand = "/start"

func newDefaultState(controller *controller.Controller, FSM *FSM, start state) *defaultState {
	return &defaultState{fsm: FSM, controller: controller, start: start}
}

func (n *defaultState) Handle(ctx context.Context, telectx tele.Context) error {
	if telectx.Message().Text != startCommand {
		return n.controller.CreateNote(ctx, telectx)
	}

	//n.fsm.SetState(n.start)
	return n.fsm.Start.Handle(ctx, telectx)
}
