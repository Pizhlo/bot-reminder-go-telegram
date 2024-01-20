package fsm

import (
	"context"
	"sync"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	tele "gopkg.in/telebot.v3"
)

type FSM struct {
	listNote     state
	createNote   state
	defaultState state
	start        state
	location     state
	current      state
	mu           sync.RWMutex
}

type state interface {
	Handle(ctx context.Context, telectx tele.Context) error
}

func NewFSM(controller *controller.Controller) *FSM {
	fsm := &FSM{mu: sync.RWMutex{}}

	fsm.createNote = newCreateNoteState(controller, fsm)
	fsm.defaultState = newDefaultState(controller, fsm)
	fsm.listNote = newListNoteState(fsm, controller)
	fsm.location = newLocationState(fsm, controller)
	fsm.start = newStartState(fsm, controller, fsm.location)

	// когда пользователь только начал пользоваться, ожидаем команду старт
	fsm.current = fsm.start

	return fsm
}

func (f *FSM) setState(state state) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.current = state
}

func (f *FSM) Handle(ctx context.Context, telectx tele.Context) error {
	return f.current.Handle(ctx, telectx)
}
