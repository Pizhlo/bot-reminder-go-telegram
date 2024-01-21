package controller

import (
	"context"

	user_model "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/note"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/user"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	logger  *logrus.Logger
	userSrv *user.UserService
	noteSrv *note.NoteService
}

func New(userSrv *user.UserService, noteSrv *note.NoteService) *Controller {
	return &Controller{logger: logrus.New(), userSrv: userSrv, noteSrv: noteSrv}
}

// CheckUser проверяет, известен ли пользователь боту
func (c *Controller) CheckUser(ctx context.Context, tgID int64) bool {
	return c.userSrv.CheckUser(ctx, tgID)
}

func (c *Controller) GetAllUsers(ctx context.Context) []*user_model.User {
	return c.userSrv.GetAll(ctx)
}
