package fsm

import (
	"context"
	"sync"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Менеджер для управления состояниями бота
type FSM struct {
	// Состояние перечисления напомианий
	ListNote state
	// Состояние создания заметки
	createNote state
	// Дефолтное состояние бота, в котором он воспринимает любой текст как заметку
	DefaultState state
	// Состояние для обработки команды старт (кнопки меню) от пользователя
	Start state
	// Состояние для обработки геолокации от пользователя
	location state
	// Текущее состояние, в котором находится бот
	current state
	// Состояние для поиска заметок по тексту
	SearchNoteByText state
	// Состояние для обработки названия напоминания
	ReminderName state
	// Состояние для обработки времени напоминания
	ReminderTime state
	// Состояние для обработки напоминаний несколько раз в день (раз в неск. минут / часов)
	SeveralTimesDay state
	// Состояние для обработки напоминаний раз в неск. минут
	MinutesDuration state
	// Состояние для обработки напоминаний раз в неск. часов
	HoursDuration state
	// Состояние для обработки напоминаний раз в неделю
	EveryWeek state
	// Состояние для обработки напоминаний раз в неск. дней
	SeveralDays state
	// Состояние для обработки числа месяца, в которое присылать уведомление.
	// Тип: напоминание раз в месяц
	DaysDuration state
	// Состояние для обработки напоминаний раз в месяц
	Month state
	// Состояние для обработки напоминаний раз в год
	Year state
	// Состояние для обработки одноразового напоминания (дата выбирается в календаре)
	Once state
	// Состояние для поиска заметок по одной дате
	SearchNoteOneDate state
	mu                sync.RWMutex
	logger            *logrus.Logger
}

// Интерфейс для управления состояниями бота
type state interface {
	Handle(ctx context.Context, telectx tele.Context) error
	Name() string
	Next()
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
	fsm.SearchNoteOneDate = newSearchNoteOneDateState(controller, fsm)

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

func (f *FSM) SetNext() {
	f.current.Next()
}

func (f *FSM) Current() state {
	return f.current
}
