package timezone

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"gopkg.in/telebot.v3"
)

type timezoneHandler struct {
	userTimezone model.UserTimezone
	srv          server
}

type server interface {
	GetUserID(ctx context.Context, tgID int64) (int, error)
	SaveTimezone(id int, tz model.UserTimezone) error
}

func New(userTimezone model.UserTimezone, srv server) *timezoneHandler {
	return &timezoneHandler{userTimezone: userTimezone, srv: srv}
}

func (h *timezoneHandler) Handle(ctx telebot.Context) error {
	return h.saveTimezone(ctx.Chat().ID)
}

func (h *timezoneHandler) saveTimezone(tgID int64) error {
	c, cancel := context.WithCancel(context.TODO()) // тот ли контекст?
	defer cancel()

	id, err := h.srv.GetUserID(c, tgID) // получаем, какой айди присвоен пользователю в бд
	if err != nil {
		return err
	}

	return h.srv.SaveTimezone(id, h.userTimezone)
}
