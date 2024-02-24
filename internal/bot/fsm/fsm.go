package fsm

import (
	"context"
	"sync"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type FSM struct {
	ListNote         state
	createNote       state
	DefaultState     state
	Start            state
	location         state
	current          state
	SearchNoteByText state
	ReminderName     state
	ReminderTime     state
	SeveralTimesDay  state
	MinutesDuration  state
	HoursDuration    state
	EveryWeek        state
	SeveralDays      state
	DaysDuration     state
	Month            state
	Year             state
	Once             state
	mu               sync.RWMutex
	logger           *logrus.Logger
}

type state interface {
	Handle(ctx context.Context, telectx tele.Context) error
	Name() string
}

func NewFSM(controller *controller.Controller, known bool) *FSM {
	fsm := &FSM{mu: sync.RWMutex{}, logger: logger.New()}

	fsm.location = newLocationState(fsm, controller)

	start := newStartState(fsm, controller, fsm.location, fsm.DefaultState)
	fsm.Start = start

	defaultState := newDefaultState(controller, fsm)
	fsm.DefaultState = defaultState

	// note
	fsm.createNote = newCreateNoteState(controller, fsm)
	fsm.ListNote = newListNoteState(fsm, controller)
	fsm.SearchNoteByText = newSearchNoteByTextState(controller, fsm)

	// reminder
	fsm.ReminderName = newReminderNameState(controller, fsm)
	fsm.ReminderTime = newReminderTimeState(controller, fsm)
	fsm.SeveralTimesDay = newSeveralTimesState(controller, fsm)
	fsm.MinutesDuration = newMinutesDurationState(controller, fsm)
	fsm.HoursDuration = newHoursDurationState(controller, fsm)
	fsm.EveryWeek = newWeekDayState(controller, fsm)
	fsm.SeveralDays = newSeveralDaysState(controller, fsm)
	fsm.DaysDuration = newDaysDurationState(controller, fsm)
	fsm.Month = newMonthState(controller, fsm)
	fsm.Year = newYearState(controller, fsm)
	fsm.Once = newOnceReminderState(controller, fsm)

	// когда пользователь только начал пользоваться, ожидаем команду старт
	if !known {
		//fsm.current = fsm.Start
	} else {
		fsm.current = fsm.DefaultState
	}

	return fsm
}

func (f *FSM) SetState(state state) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.logger.Debugf("Setting state to: %v\n", state.Name())

	f.current = state
}

func (f *FSM) Handle(ctx context.Context, telectx tele.Context) error {
	f.logger.Debugf("Handling request. Current state: %v. Command: %s\n", f.current.Name(), telectx.Message().Text)
	return f.current.Handle(ctx, telectx)
}

func (f *FSM) Name() string {
	return f.current.Name()
}
