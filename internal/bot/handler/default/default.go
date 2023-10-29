package default_handler

import (
	"context"
	"errors"

	api_err "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/server"
	tele "gopkg.in/telebot.v3"
)

type DefaultHandler struct {
	srv *server.Server
}

func New(srv *server.Server) *DefaultHandler {
	return &DefaultHandler{srv}
}

var (
	locationMenu = &tele.ReplyMarkup{ResizeKeyboard: true}
	LocationBtn  = locationMenu.Location("Отправить геолокацию")
	RejectBtn    = locationMenu.Text(`Отказаться`)
)

// handles /start command
func (h *DefaultHandler) Handle(ctx tele.Context) error {
	c, cancel := context.WithCancel(context.TODO()) // тот ли контекст?
	defer cancel()

	_, err := h.srv.GetUserID(c, ctx.Chat().ID) // проверяем, известен ли нам пользователь
	if err != nil {
		if errors.Is(err, api_err.ErrUserNotFound) { // если пользователь новый - отправляем клавиатуру с запросом гео
			locationMenu.Reply(
				locationMenu.Row(LocationBtn),
				locationMenu.Row(RejectBtn),
			)
			return ctx.Send(messages.StartMessageLocation, locationMenu)
		}
		return err
	}
	return ctx.Send(messages.StartMessage) // если пользователь известен
}
