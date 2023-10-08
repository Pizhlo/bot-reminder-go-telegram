package controller

import (
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/server"
	tele "gopkg.in/telebot.v3"
)

type Controller struct {
	srv *server.Server
}

func New(srv *server.Server) *Controller {
	return &Controller{srv}
}

func (c *Controller) SetupBot() {
	c.srv.Bot.Handle(`/start`, func(c tele.Context) error {
		return c.Send("Hello!")
	})

	c.startBot()
}

func (c *Controller) startBot() {
	c.srv.Logger.Info().Msg(`successfully loaded app`)
	c.srv.Bot.Start()
}
