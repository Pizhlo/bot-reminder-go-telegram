package fsm

import (
	"context"
	"sync"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type FSM struct {
	listNote     state
	createNote   state
	defaultState state
	Start        state
	location     state
	current      state
	mu           sync.RWMutex
	logger       *logrus.Logger
}

type state interface {
	Handle(ctx context.Context, telectx tele.Context) error
}

func NewFSM(controller *controller.Controller, known bool) *FSM {
	fsm := &FSM{mu: sync.RWMutex{}, logger: logrus.New()}

	fsm.location = newLocationState(fsm, controller)

	start := newStartState(fsm, controller, fsm.location, fsm.defaultState)
	fsm.Start = start

	defaultState := newDefaultState(controller, fsm, fsm.Start)
	fsm.defaultState = defaultState

	fsm.createNote = newCreateNoteState(controller, fsm)

	fsm.listNote = newListNoteState(fsm, controller)

	// когда пользователь только начал пользоваться, ожидаем команду старт
	if !known {
		fsm.current = fsm.Start
	} else {
		fsm.current = fsm.defaultState
	}

	return fsm
}

func (f *FSM) SetState(state state) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.logger.Debugf("Setting state to: %v\n", state)

	f.current = state
}

func (f *FSM) Handle(ctx context.Context, telectx tele.Context) error {
	f.logger.Debugf("Handling request. Current state: %v. Command: %s\n", f.current, telectx.Message().Text)
	return f.current.Handle(ctx, telectx)
}
