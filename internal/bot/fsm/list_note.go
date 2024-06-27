package fsm

import (
	"context"
	"strconv"
	"strings"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние перечисления заметок
type listNote struct {
	controller *controller.Controller
	fsm        *FSM
	name       stateName
	next       state
}

func newListNoteState(FSM *FSM, controller *controller.Controller) *listNote {
	return &listNote{controller, FSM, listNoteName, nil}
}

const (
	// префикс для удаления заметки
	deleteNotePrefix = "/dn"
	// префикс для редактирования заметки
	editNotePrefix = "/editn"
)

func (n *listNote) Handle(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)
	msg := telectx.Message().Text

	if strings.HasPrefix(msg, deleteNotePrefix) {
		return n.controller.DeleteNoteByID(ctx, telectx)
	} else if strings.HasPrefix(msg, editNotePrefix) {
		noteIDString, _ := strings.CutPrefix(msg, editNotePrefix)
		noteID, err := strconv.Atoi(noteIDString)
		if err != nil {
			return err
		}

		// создаем состояние, которое будет хранить номер редактируемой заметки
		n.fsm.editNote = newEditNoteState(n.controller, n.fsm, noteID)

		// проставляем состояние в текущее
		n.fsm.current = n.fsm.editNote

		// запрашиваем текст заметки
		return n.controller.AskNoteText(ctx, telectx)
	} else {
		n.fsm.SetToDefault()
		return n.fsm.Handle(ctx, telectx)
	}

	//return n.controller.ListNotes(ctx, telectx)
}

func (n *listNote) Name() string {
	return string(n.name)
}

func (n *listNote) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.defaultState
}
