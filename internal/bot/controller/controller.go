package controller

import (
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/note"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/user"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	logger  *logrus.Logger
	userSrv *user.UserService
	noteSrv *note.NoteService
}

func NewMyController(userSrv *user.UserService, noteSrv *note.NoteService) *Controller {
	return &Controller{logger: logrus.New(), userSrv: userSrv, noteSrv: noteSrv}
}
