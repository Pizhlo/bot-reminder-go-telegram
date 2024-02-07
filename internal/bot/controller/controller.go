package controller

import (
	"context"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	user_model "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/note"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/user"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type Controller struct {
	logger *logrus.Logger
	bot    *tele.Bot
	// отвечает за информацию о пользователях
	userSrv *user.UserService
	// отвечает за обработку заметок
	noteSrv *note.NoteService
	// последнее отправленное сообщение (для редактирования)
	lastMsg map[int64]*tele.Message
}

const (
	htmlParseMode     = "HTML"
	markdownParseMode = "markdown"
)

func New(userSrv *user.UserService, noteSrv *note.NoteService, bot *tele.Bot) *Controller {
	return &Controller{logger: logger.New(), userSrv: userSrv, noteSrv: noteSrv, bot: bot, lastMsg: make(map[int64]*tele.Message)}
}

// CheckUser проверяет, известен ли пользователь боту
func (c *Controller) CheckUser(ctx context.Context, tgID int64) bool {
	c.noteSrv.SaveUser(tgID)
	return c.userSrv.CheckUser(ctx, tgID)
}

// GetAllUsers возвращает всех зарегистрированных пользователей
func (c *Controller) GetAllUsers(ctx context.Context) []*user_model.User {
	return c.userSrv.GetAll(ctx)
}

// HandleError сообщает об ошибке в канал.
// Также сообщает пользователю об ошибке - редактирует сообщение
func (c *Controller) HandleError(ctx tele.Context, err error) {
	msg := fmt.Sprintf(messages.ErrorMessageChannel, ctx.Message().Text, err)

	// sendErr := ctx.Send(messages.ErrorMessageUser)
	// if sendErr != nil {
	// 	c.logger.Errorf("Error while sending error message to user. Error: %+v\n", sendErr)
	// }

	editErr := ctx.Edit(messages.ErrorMessageUser)
	if editErr != nil {
		c.logger.Errorf("Error while sending error message to user. Error: %+v\n", editErr)
	}

	_, channelErr := c.bot.Send(&tele.Chat{ID: -1001890622926}, msg, &tele.SendOptions{
		ParseMode: markdownParseMode,
	})
	if channelErr != nil {
		c.logger.Errorf("Error while sending error message to channel. Error: %+v\n", editErr)
	}
}

func (c *Controller) SaveLastMsg(user int64, msg *tele.Message) {
	c.lastMsg[user] = msg
}
