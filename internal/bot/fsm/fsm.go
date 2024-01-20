package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	tele "gopkg.in/telebot.v3"
)

type FSM struct {
	listNote     state
	createNote   state
	defaultState state
	start        state
	current      state
}

type state interface {
	Handle(ctx context.Context, telectx tele.Context) error
}

func NewFSM(controller *controller.Controller) *FSM {
	fsm := &FSM{}

	fsm.createNote = newCreateNoteState(controller, fsm)
	fsm.defaultState = newDefaultState(controller, fsm)
	fsm.listNote = newListNoteState(fsm, controller)
	fsm.start = newStartState(fsm, controller)

	// когда пользователь только начал пользоваться, ожидаем команду старт
	fsm.current = fsm.start

	return fsm
}

func (f *FSM) Handle(ctx context.Context, telectx tele.Context) error {
	return f.current.Handle(ctx, telectx)
}
