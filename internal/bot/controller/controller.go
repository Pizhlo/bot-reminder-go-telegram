package controller

import (
	"context"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	user_model "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
	gocron "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/scheduler"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/note"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/reminder"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/user"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
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
	// отвечает за напоминания
	reminderSrv *reminder.ReminderService
	// scheduler вызывает функции в указанное время
	scheduler *gocron.Scheduler
	// отражает, используется ли календарь заметок
	noteCalendar map[int64]bool
	// отражает, используется ли календарь напоминаний
	reminderCalendar map[int64]bool
}

const (
	htmlParseMode     = "HTML"
	markdownParseMode = "markdown"
)

func New(userSrv *user.UserService, noteSrv *note.NoteService, bot *tele.Bot, reminderSrv *reminder.ReminderService) (*Controller, error) {
	sch, err := gocron.New()
	if err != nil {
		return nil, err
	}

	return &Controller{logger: logger.New(),
		userSrv:          userSrv,
		noteSrv:          noteSrv,
		bot:              bot,
		reminderSrv:      reminderSrv,
		scheduler:        sch,
		noteCalendar:     make(map[int64]bool),
		reminderCalendar: make(map[int64]bool)}, nil
}

// CheckUser проверяет, известен ли пользователь боту
func (c *Controller) CheckUser(ctx context.Context, tgID int64) bool {
	c.noteSrv.SaveUser(tgID)
	c.reminderSrv.SaveUser(tgID)
	c.noteCalendar[tgID] = false
	c.reminderCalendar[tgID] = false
	return c.userSrv.CheckUser(ctx, tgID)
}

// GetAllUsers возвращает всех зарегистрированных пользователей
func (c *Controller) GetAllUsers(ctx context.Context) []*user_model.User {
	return c.userSrv.GetAll(ctx)
}

// HandleError сообщает об ошибке в канал.
// Также сообщает пользователю об ошибке - редактирует сообщение
func (c *Controller) HandleError(ctx tele.Context, err error, state string) {
	var msg string
	if ctx.Message().Sender != c.bot.Me {
		msg = fmt.Sprintf(messages.ErrorMessageChannel, ctx.Message().Text, state, err)
	} else {
		msg = fmt.Sprintf(messages.ErrorMessageChannel, ctx.Callback().Unique, state, err)
	}

	editErr := ctx.EditOrSend(messages.ErrorMessageUser, view.BackToMenuBtn())
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

// SaveUsers сохраняет пользователей в сервисах
func (c *Controller) SaveUsers(ctx context.Context, users []*user_model.User) {
	for _, u := range users {
		c.noteSrv.SaveUser(u.TGID)
		c.userSrv.SaveUser(ctx, u.TGID, u)
		c.reminderSrv.SaveUser(u.TGID)
	}
}
