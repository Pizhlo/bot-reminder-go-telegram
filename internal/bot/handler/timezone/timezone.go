package timezone

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/server"
	"gopkg.in/telebot.v3"
)

type timezoneHandler struct {
	userTimezone model.UserTimezone
	srv          *server.Server
}

func New(userTimezone model.UserTimezone, srv *server.Server) *timezoneHandler {
	return &timezoneHandler{userTimezone, srv}
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
