package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	tele "gopkg.in/telebot.v3"
)

// дефолтное состояние бота
type defaultState struct {
	FSM        *FSM
	controller *controller.Controller
}

func newDefaultState(controller *controller.Controller, FSM *FSM) *defaultState {
	return &defaultState{FSM: FSM, controller: controller}
}

func (n *defaultState) Handle(ctx context.Context, telectx tele.Context) error {
	return n.controller.CreateNote(ctx, telectx)
}
