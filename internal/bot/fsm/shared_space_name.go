package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние для создания совместного пространства
type sharedSpaceName struct {
	controller *controller.Controller
	fsm        *FSM
	name       stateName
	next       state
}

func newSharedSpaceName(FSM *FSM, controller *controller.Controller) *sharedSpaceName {
	return &sharedSpaceName{controller, FSM, createSharedSpaceName, nil}
}

// Создаем новое совместное пространство
func (n *sharedSpaceName) Handle(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	logrus.Debugf("createSharedSpace state: calling controller.CreateSharedSpace()\n")
	return n.controller.CreateSharedSpace(ctx, telectx)
}

func (n *sharedSpaceName) Name() string {
	return string(n.name)
}

func (n *sharedSpaceName) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.defaultState
}
