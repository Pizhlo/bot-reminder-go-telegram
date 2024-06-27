package fsm

import (
	"context"
	"fmt"
	"sync"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type stateName string

// названия состояний
const (
	defaultStateName                   stateName = "default"
	startStateName                     stateName = "start"
	listNoteName                       stateName = "list_note"
	createNoteName                     stateName = "create_note"
	editNoteName                       stateName = "edit_note"
	daysDurationName                   stateName = "days_duration"
	hoursStateName                     stateName = "hours"
	listReminderName                   stateName = "list_reminder"
	locationStateName                  stateName = "location"
	minutesStateName                   stateName = "minutes_duration"
	timesReminderName                  stateName = "times_reminder"
	monthStateName                     stateName = "month"
	dateReminderName                   stateName = "date_reminder"
	reminderNameState                  stateName = "reminder_name"
	reminderTimeState                  stateName = "reminder_time"
	searchNoteByTextStateName          stateName = "search_note_by_text"
	searchNoteByDatetStateName         stateName = "search_note_by_date"
	searchNoteByTwoDatesState          stateName = "search_note_by_two_dates"
	searchNoteByTwoDatesStateFirstDay  stateName = "search_note_by_two_dates_first_day"
	searchNoteByTwoDatesStateSecondDay stateName = "search_note_by_two_dates_second_day"
	severalDaysState                   stateName = "several_days"
	severalTimesDayState               stateName = "several_times_a_day"
	weekDayState                       stateName = "every_week"
	yearReminderState                  stateName = "every_year"

	bugReportState stateName = "bug_report"
)

// Менеджер для управления состояниями бота
type FSM struct {
	// Состояние перечисления заметок
	ListNote state
	// Состояние перечисления напоминаний
	ListReminder state
	// Состояние создания заметки
	createNote state
	// Состояние для редактирования заметки
	editNote state
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
	// Состояние для обработки напоминаний в указанные времена (10:30, 12:30, 15:00)
	Times state
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
	// Состояние для поиска заметок по двум датам
	SearchNoteTwoDates state
	// Состояние для баг репорта
	BugReportState state
	mu             sync.RWMutex
}

// Интерфейс для управления состояниями бота
type state interface {
	Handle(ctx context.Context, telectx tele.Context) error
	Name() string
	Next() state
}

func NewFSM(controller *controller.Controller) *FSM {
	fsm := &FSM{mu: sync.RWMutex{}}

	fsm.location = newLocationState(fsm, controller)

	start := newStartState(fsm, controller, fsm.location, fsm.defaultState)
	fsm.start = start

	defaultState := newDefaultState(controller, fsm)
	fsm.defaultState = defaultState

	// note
	fsm.createNote = newCreateNoteState(controller, fsm)
	fsm.ListNote = newListNoteState(fsm, controller)
	// fsm.editNote = newEditNoteState(controller, fsm)
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
	fsm.Times = newTimesReminderState(controller, fsm)
	fsm.EveryWeek = newEveryWeekState(controller, fsm)
	fsm.SeveralDays = newSeveralDaysState(controller, fsm)
	fsm.DaysDuration = newDaysDurationState(controller, fsm)
	fsm.Month = newMonthState(controller, fsm)
	fsm.Year = newYearState(controller, fsm)
	fsm.Date = newDateReminderState(controller, fsm)

	fsm.BugReportState = newBugReportState(controller, fsm)

	// когда пользователь только начал пользоваться, ожидаем команду старт
	fsm.current = fsm.defaultState

	return fsm
}

// SetState устанавливает текущее состояние в переданное
func (f *FSM) SetState(state state) {
	f.mu.Lock()
	defer f.mu.Unlock()

	logrus.Debugf("Setting state to: %v\n", state.Name())

	f.current = state
}

// SetToDefault устанавливает текущее состояние FSM в дефолтное
func (f *FSM) SetToDefault() {
	f.SetState(f.defaultState)
}

func (f *FSM) Handle(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Handling request. Current state: %v. Command: %s\n", f.current.Name(), telectx.Message().Text)
	return f.current.Handle(ctx, telectx)
}

// Name возвращает название текущего состояния
func (f *FSM) Name() string {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.current.Name()
}

// SetNext переключает состояние бота на следующее
func (f *FSM) SetNext() {
	f.SetState(f.current.Next())
}

// Current возвращает текущее состояние
func (f *FSM) Current() state {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.current
}

// SetFromString устанавливает текущее состояние в переданное по названию
func (s *FSM) SetFromString(stateStr string) error {
	stateName := stateName(stateStr)
	state, err := s.parseString(stateName)
	if err != nil {
		logrus.Errorf("error while setting state on start: %v", err)
		s.SetToDefault()
		return nil
	}

	s.SetState(state)
	return nil
}

// parseString парсит переданное название состояния.
// Возвращает ошибку, если такого состояния не найдено
func (s *FSM) parseString(state stateName) (state, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch state {
	case startStateName:
		return s.start, nil
	case defaultStateName:
		return s.defaultState, nil
	case listNoteName:
		return s.ListNote, nil
	case createNoteName:
		return s.createNote, nil
	case editNoteName:
		return s.editNote, nil
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
	case timesReminderName:
		return s.Times, nil
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
