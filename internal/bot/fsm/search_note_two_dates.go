package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние для поиска заметок по одной дате
type searchNoteTwoDate struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       stateName
	next       state
}

func newSearchNoteTwoDateState(controller *controller.Controller, FSM *FSM) *searchNoteTwoDate {
	return &searchNoteTwoDate{controller, FSM, logger.New(), searchNoteByTwoDatesState, newSelectedDayFirst(controller, FSM)}
}

func (n *searchNoteTwoDate) Handle(ctx context.Context, telectx tele.Context) error {
	// если пользователь прислал текст - сохраняем
	if !telectx.Message().Sender.IsBot && telectx.Message().Text != "" {
		return n.controller.CreateNote(ctx, telectx)
	}

	return n.controller.SearchNoteByTwoDates(ctx, telectx)
}

func (n *searchNoteTwoDate) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.defaultState
}

func (n *searchNoteTwoDate) Name() string {
	return string(n.name)
}

// Next state

// Состояние, в котором бот обрабатывает первый выбранный день
type selectedDayFirst struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
	next       state
}

func newSelectedDayFirst(controller *controller.Controller, FSM *FSM) *selectedDayFirst {
	return &selectedDayFirst{controller: controller,
		fsm: FSM, logger: logger.New(),
		name: "search notes by two dates: first selected day",
		next: newSelectedDaySecond(controller, FSM)}
}

func (n *selectedDayFirst) Name() string {
	return string(n.name)
}

func (n *selectedDayFirst) Handle(ctx context.Context, telectx tele.Context) error {
	return n.controller.SearchNoteByTwoDatesFirstDate(ctx, telectx)
}

func (n *selectedDayFirst) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.defaultState
}

// Next state

// Состояние, в котором бот обрабатывает второй выбранный день
type selectedDaySecond struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
	next       state
}

func newSelectedDaySecond(controller *controller.Controller, FSM *FSM) *selectedDaySecond {
	s := &selectedDaySecond{controller: controller,
		fsm:    FSM,
		logger: logger.New(),
		name:   "search notes by two dates: second selected day",
		next:   nil}

	s.next = s

	return s
}

func (n *selectedDaySecond) Name() string {
	return string(n.name)
}

func (n *selectedDaySecond) Handle(ctx context.Context, telectx tele.Context) error {
	return n.controller.SearchNoteByTwoDatesSecondDate(ctx, telectx)
}

func (n *selectedDaySecond) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.defaultState
}
