package controller

import (
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/calendar"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/server"
)

type Controller struct {
	srv      *server.Server
	calendar *calendar.Calendar
	logger   *logger.Logger
}

func New(srv *server.Server, calendar *calendar.Calendar, logger *logger.Logger) *Controller {
	return &Controller{srv, calendar, logger}
}
