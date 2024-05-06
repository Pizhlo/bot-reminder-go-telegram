package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние обработки баг репорта
type bugReport struct {
	controller *controller.Controller
	fsm        *FSM
	name       stateName
	next       state
}

func newBugReportState(controller *controller.Controller, FSM *FSM) *bugReport {
	return &bugReport{controller, FSM, bugReportState, nil}
}

func (n *bugReport) Handle(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)
	err := n.controller.BugReport(ctx, telectx)
	if err != nil {
		return err
	}

	n.fsm.SetNext()

	return nil
}

func (n *bugReport) Name() string {
	return string(n.name)
}

func (n *bugReport) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.defaultState
}
