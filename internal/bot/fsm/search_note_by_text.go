package fsm

import (
	"context"
	"strings"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние для поиска заметок по тексту
type searchNoteByTextState struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
	next       state
}

func newSearchNoteByTextState(controller *controller.Controller, FSM *FSM) *searchNoteByTextState {
	return &searchNoteByTextState{controller, FSM, logger.New(), searchNoteByTextStateName, nil}
}

func (n *searchNoteByTextState) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling search note by text. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	if strings.HasPrefix(telectx.Message().Text, "/del") {
		return n.controller.DeleteNoteByID(ctx, telectx)
	}

	err := n.controller.SearchNoteByText(ctx, telectx)
	if err != nil {
		return err
	}

	n.fsm.SetNext()

	return nil
}

func (n *searchNoteByTextState) Name() string {
	return n.name
}

func (n *searchNoteByTextState) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.defaultState
}
