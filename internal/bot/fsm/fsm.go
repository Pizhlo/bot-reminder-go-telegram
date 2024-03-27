package fsm

import (
	"context"
	"fmt"
	"sync"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// названия состояний
const (
	defaultStateName           = "default"
	startStateName             = "start"
	listNoteName               = "list_note"
	createNoteName             = "create_note"
	daysDurationName           = "days_duration"
	hoursStateName             = "hours"
	listReminderName           = "list_reminder"
	locationStateName          = "location"
	minutesStateName           = "minutes_duration"
	monthStateName             = "month"
	dateReminderName           = "date_reminder"
	reminderNameState          = "reminder_name"
	reminderTimeState          = "reminder_time"
	searchNoteByTextStateName  = "search_note_by_text"
	searchNoteByDatetStateName = "search_note_by_date"
	searchNoteByTwoDatesState  = "search_note_by_two_dates"
	severalDaysState           = "several_days"
	severalTimesDayState       = "several_times_a_day"
	weekDayState               = "every_week"
	yearReminderState          = "every_year"
)

// Менеджер для управления состояниями бота
type FSM struct {
	// Состояние перечисления заметок
	ListNote state
	// Состояние перечисления напоминаний
	ListReminder state
	// Состояние создания заметки
	createNote state
	// Дефолтное состояние бота, в котором он воспринимает любой текст как заметку
	defaultState state
	// Состояние для обработки команды старт (кнопки меню) от пользователя
	start state
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
	Date state
	// Состояние для поиска заметок по одной дате
	SearchNoteOneDate state
	//Состояние для поиска заметок по двум датам
	SearchNoteTwoDates state
	mu                 sync.RWMutex
	logger             *logrus.Logger
}

// Интерфейс для управления состояниями бота
type state interface {
	Handle(ctx context.Context, telectx tele.Context) error
	Name() string
	Next() state
}

func NewFSM(controller *controller.Controller) *FSM {
	fsm := &FSM{mu: sync.RWMutex{}, logger: logger.New()}

	fsm.location = newLocationState(fsm, controller)

	start := newStartState(fsm, controller, fsm.location, fsm.defaultState)
	fsm.start = start

	defaultState := newDefaultState(controller, fsm)
	fsm.defaultState = defaultState

	// note
	fsm.createNote = newCreateNoteState(controller, fsm)
	fsm.ListNote = newListNoteState(fsm, controller)
	fsm.SearchNoteByText = newSearchNoteByTextState(controller, fsm)
	fsm.SearchNoteOneDate = newSearchNoteOneDateState(controller, fsm)
	fsm.SearchNoteTwoDates = newSearchNoteTwoDateState(controller, fsm)

	// reminder
	fsm.ListReminder = newListReminderState(controller, fsm)
	fsm.ReminderName = newReminderNameState(controller, fsm)
	fsm.ReminderTime = newReminderTimeState(controller, fsm)
	fsm.SeveralTimesDay = newSeveralTimesState(controller, fsm)
	fsm.MinutesDuration = newMinutesDurationState(controller, fsm)
	fsm.HoursDuration = newHoursDurationState(controller, fsm)
	fsm.EveryWeek = newEveryWeekState(controller, fsm)
	fsm.SeveralDays = newSeveralDaysState(controller, fsm)
	fsm.DaysDuration = newDaysDurationState(controller, fsm)
	fsm.Month = newMonthState(controller, fsm)
	fsm.Year = newYearState(controller, fsm)
	fsm.Date = newDateReminderState(controller, fsm)

	// когда пользователь только начал пользоваться, ожидаем команду старт
	fsm.current = fsm.defaultState

	return fsm
}

// SetState устанавливает текущее состояние в переданное
func (f *FSM) SetState(state state) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.logger.Debugf("Setting state to: %v\n", state.Name())

	f.current = state
}

// SetToDefault устанавливает текущее состояние FSM в дефолтное
func (f *FSM) SetToDefault() {
	f.SetState(f.defaultState)
}

func (f *FSM) Handle(ctx context.Context, telectx tele.Context) error {
	f.logger.Debugf("Handling request. Current state: %v. Command: %s\n", f.current.Name(), telectx.Message().Text)
	return f.current.Handle(ctx, telectx)
}

// Name возвращает название текущего состояния
func (f *FSM) Name() string {
	return f.current.Name()
}

// SetNext переключает состояние бота на следующее
func (f *FSM) SetNext() {
	next := f.current.Next()
	f.logger.Debugf("Setting state to next. Next: %v\n", next.Name())
	f.current = next
}

// Current возвращает текущее состояние
func (f *FSM) Current() state {
	return f.current
}

// SetFromString устанавливает текущее состояние в переданное по названию
func (s *FSM) SetFromString(stateStr string) error {
	state, err := s.parseString(stateStr)
	if err != nil {
		return err
	}

	s.SetState(state)
	return nil
}

// parseString парсит переданное название состояния.
// Возвращает ошибку, если такого состояния не найдено
func (s *FSM) parseString(state string) (state, error) {
	switch state {
	case startStateName:
		return s.start, nil
	case defaultStateName:
		return s.defaultState, nil
	case listNoteName:
		return s.ListNote, nil
	case createNoteName:
		return s.createNote, nil
	case daysDurationName:
		return s.DaysDuration, nil
	case hoursStateName:
		return s.HoursDuration, nil
	case listReminderName:
		return s.ListReminder, nil
	case locationStateName:
		return s.location, nil
	case minutesStateName:
		return s.MinutesDuration, nil
	case monthStateName:
		return s.Month, nil
	case dateReminderName:
		return s.Date, nil
	case reminderNameState:
		return s.ReminderName, nil
	case reminderTimeState:
		return s.ReminderTime, nil
	case searchNoteByTextStateName:
		return s.SearchNoteByText, nil
	case searchNoteByDatetStateName:
		return s.SearchNoteOneDate, nil
	case searchNoteByTwoDatesState:
		return s.SearchNoteTwoDates, nil
	case severalDaysState:
		return s.SeveralDays, nil
	case severalTimesDayState:
		return s.SeveralTimesDay, nil
	case weekDayState:
		return s.EveryWeek, nil
	case yearReminderState:
		return s.Year, nil
	default:
		return nil, fmt.Errorf("unknown state: %s", state)
	}
}
