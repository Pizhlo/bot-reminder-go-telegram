package controller

import (
	"fmt"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/server"
	pkg_err "github.com/pkg/errors"
	tele "gopkg.in/telebot.v3"
)

type Controller struct {
	srv *server.Server
}

func New(srv *server.Server) *Controller {
	return &Controller{srv}
}

const (
	startCmd = `/start`
)

var (
	locationMenu = &tele.ReplyMarkup{ResizeKeyboard: true}
	locationBtn  = locationMenu.Location("Отправить геолокацию")
	rejectBtn    = locationMenu.Text(`Отказаться`)
)

func (c *Controller) SetupBot() error {
	c.srv.Bot.Handle(startCmd, func(ctx tele.Context) error {
		// проверяем, известен ли нам пользователь
		_, err := c.srv.UserCacheEditor.GetUser(ctx.Chat().ID)
		return c.StartMsg(ctx, err)
	})

	c.srv.Bot.Handle(&rejectBtn, func(ctx tele.Context) error {
		return ctx.Send(`text`, tele.RemoveKeyboard)
	})

	c.srv.Bot.Handle(&locationBtn, func(ctx tele.Context) error {
		fmt.Println("location")
		err := c.saveUser(ctx.Chat().ID)
		if err != nil {
			return err
		}

		tz := c.saveTimezone(ctx.Chat().ID, ctx.Query().Location)
		return ctx.Send(fmt.Sprintf(messages.LocationMessage, tz), tele.RemoveKeyboard)
	})

	if err := c.sendStartupMsg(); err != nil {
		return pkg_err.Wrap(err, `unable to send startup message`)
	}

	c.startBot()

	return nil
}

func (c *Controller) sendStartupMsg() error {
	_, err := c.srv.Bot.Send(tele.Recipient(&tele.Chat{ID: -1001890622926}), "#запуск\nБот запущен")
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) startBot() {
	c.srv.Logger.Info().Msg(`successfully loaded app`)
	c.srv.Bot.Start()
}
