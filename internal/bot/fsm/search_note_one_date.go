package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние для поиска заметок по одной дате
type searchNoteOneDate struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
	next       state
}

func newSearchNoteOneDateState(controller *controller.Controller, FSM *FSM) *searchNoteOneDate {
	return &searchNoteOneDate{controller, FSM, logger.New(), "search note by one date", newSelectedDay(controller, FSM)}
}

func (n *searchNoteOneDate) Handle(ctx context.Context, telectx tele.Context) error {
	// если пользователь прислал текст - сохраняем
	if !telectx.Message().Sender.IsBot && telectx.Message().Text != "" {
		return n.controller.CreateNote(ctx, telectx)
	}

	return n.controller.SearchNoteByOnedate(ctx, telectx)
}

func (n *searchNoteOneDate) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.defaultState
}

func (n *searchNoteOneDate) Name() string {
	return n.name
}

// Next state

// Состояние, в котором бот обрабатывает выбранный день
type selectedDay struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
	next       state
}

func newSelectedDay(controller *controller.Controller, FSM *FSM) *selectedDay {
	return &selectedDay{controller: controller, fsm: FSM, logger: logger.New(), name: "search notes by one date: selected day", next: nil}
}

func (n *selectedDay) Name() string {
	return n.name
}

func (n *selectedDay) Handle(ctx context.Context, telectx tele.Context) error {
	return n.controller.SearchNoteBySelectedDate(ctx, telectx)
}

func (n *selectedDay) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.defaultState
}
