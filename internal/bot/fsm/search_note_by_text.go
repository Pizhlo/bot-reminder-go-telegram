package fsm

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние для поиска заметок по тексту
type searchNoteByTextState struct {
	controller *controller.Controller
	fsm        *FSM
	name       stateName
	next       state
}

func newSearchNoteByTextState(controller *controller.Controller, FSM *FSM) *searchNoteByTextState {
	return &searchNoteByTextState{controller, FSM, searchNoteByTextStateName, nil}
}

func (n *searchNoteByTextState) Handle(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Handling search note by text. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	if strings.HasPrefix(telectx.Message().Text, "/dn") {
		return n.controller.DeleteNoteByID(ctx, telectx)
	}

	if strings.HasPrefix(telectx.Message().Text, "/editn") {
		numberString := strings.TrimPrefix(telectx.Message().Text, "/editn")

		number, err := strconv.Atoi(numberString)
		if err != nil {
			return fmt.Errorf("error converting string note ID '%s' to int: %+v", numberString, err)
		}

		n.fsm.editNote = newEditNoteState(n.controller, n.fsm, number)
		n.fsm.SetState(n.fsm.editNote)
		return n.controller.AskNoteText(ctx, telectx)
	}

	err := n.controller.SearchNoteByText(ctx, telectx)
	if err != nil {
		return err
	}

	// n.fsm.SetNext()

	return nil
}

func (n *searchNoteByTextState) Name() string {
	return string(n.name)
}

func (n *searchNoteByTextState) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.defaultState
}
