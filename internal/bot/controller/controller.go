package controller

import (
	"context"
	"fmt"
	"sync"
	"time"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	user_model "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
	gocron "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/scheduler"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/note"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/reminder"
	sharedaccess "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/shared_space"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/user"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type Controller struct {
	mu  sync.Mutex
	bot *tele.Bot
	// совместные пространства пользователей
	sharedSpace *sharedaccess.SharedSpace
	// отвечает за информацию о пользователях
	userSrv *user.UserService
	// отвечает за обработку заметок
	noteSrv *note.NoteService
	// отвечает за напоминания
	reminderSrv *reminder.ReminderService
	// schedulers - мапа с планировщиками для каждого пользователя
	schedulers map[int64]*gocron.Scheduler
	// отражает, используется ли календарь заметок
	noteCalendar map[int64]bool
	// отражает, используется ли календарь напоминаний
	reminderCalendar map[int64]bool
	channelID        int64
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
//go:generate mockgen -source ./controller.go -destination=../mocks/controller.go -package=mocks
type teleCtx interface {
	tele.Context
}

const (
	htmlParseMode     = "HTML"
	markdownParseMode = "markdown"
)

func New(userSrv *user.UserService,
	noteSrv *note.NoteService,
	bot *tele.Bot,
	reminderSrv *reminder.ReminderService,
	channelID int64,
	sharedSpace *sharedaccess.SharedSpace) *Controller {

	return &Controller{
		userSrv:          userSrv,
		noteSrv:          noteSrv,
		bot:              bot,
		reminderSrv:      reminderSrv,
		schedulers:       make(map[int64]*gocron.Scheduler),
		noteCalendar:     make(map[int64]bool),
		reminderCalendar: make(map[int64]bool),
		mu:               sync.Mutex{},
		channelID:        channelID,
		sharedSpace:      sharedSpace,
	}
}

// CheckUser проверяет, известен ли пользователь боту
func (c *Controller) CheckUser(ctx context.Context, tgID int64) bool {
	// c.noteSrv.SaveUser(tgID)
	// c.reminderSrv.SaveUser(tgID)
	// c.noteCalendar[tgID] = false
	// c.reminderCalendar[tgID] = false

	return c.userSrv.CheckUser(ctx, tgID)
}

// GetAllUsers возвращает всех зарегистрированных пользователей
func (c *Controller) GetAllUsers(ctx context.Context) ([]*user_model.User, error) {
	return c.userSrv.GetAll(ctx)
}

// HandleError сообщает об ошибке в канал.
// Также сообщает пользователю об ошибке - редактирует сообщение
func (c *Controller) HandleError(ctx tele.Context, err error, state string) {
	var msg string
	if ctx.Message().Sender != c.bot.Me {
		msg = fmt.Sprintf(messages.ErrorMessageChannelMessage, ctx.Message().Text, state, err)
	} else {
		msg = fmt.Sprintf(messages.ErrorMessageChannelMessage, ctx.Callback().Unique, state, err)
	}

	editErr := ctx.EditOrSend(messages.ErrorMessageUser, view.BackToMenuBtn())
	if editErr != nil {
		logrus.Errorf("Error while sending error message to user. Error: %+v\n", editErr)
	}

	_, channelErr := c.bot.Send(&tele.Chat{ID: -1001890622926}, msg, &tele.SendOptions{
		ParseMode: htmlParseMode,
	})
	if channelErr != nil {
		logrus.Errorf("Error while sending error message to channel. Error: %+v\n", channelErr)
	}
}

// LoadMemoryReminders выгружает недоделанные напоминания из БД
func (c *Controller) LoadMemoryReminders(ctx context.Context) error {
	return c.reminderSrv.LoadMemory(ctx)
}

// SaveUsers сохраняет пользователей в сервисах
func (c *Controller) SaveUsers(ctx context.Context, users []*user_model.User) {
	errors := []error{}
	for _, u := range users {
		c.noteSrv.SaveUser(u.TGID)

		c.reminderSrv.SaveUser(u.TGID)

		c.sharedSpace.SaveUser(u.TGID)

		loc, err := c.userSrv.GetLocation(ctx, u.TGID)
		if err != nil {
			errors = append(errors, err)
		}

		if loc == nil {
			logrus.Error("Controller: cannot create scheduler for user: timezone is empty")
			continue
		}

		err = c.reminderSrv.CreateScheduler(ctx, u.TGID, loc, c.SendReminder)
		if err != nil {
			errors = append(errors, err)
		}

		// err = c.reminderSrv.StartAllJobs(ctx, u.TGID, loc, c.SendReminder)
		// if err != nil {
		// 	errors = append(errors, err)
		// }
	}

	if len(errors) > 0 {
		logrus.Fatalf("error while saving users on start: errors: %v", errors)
	}
}

// createScheduler создает планировщика для конкретного пользователя
// func (c *Controller) createScheduler(ctx context.Context, tgID int64) error {
// 	if _, ok := c.schedulers[tgID]; !ok {
// 		loc, err := c.userSrv.GetLocation(ctx, tgID)
// 		if err != nil {
// 			return err
// 		}

// 		sch, err := gocron.New(loc)
// 		if err != nil {
// 			return err
// 		}

// 		c.schedulers[tgID] = sch
// 	}

// 	return nil
// }

// func (c *Controller) getScheduler(tgID int64) (*gocron.Scheduler, error) {
// 	if val, ok := c.schedulers[tgID]; ok {
// 		return val, nil
// 	}

// 	return nil, errors.New("no scheduler found for this user")
// }

func (c *Controller) saveUser(ctx context.Context, tgID int64, loc *time.Location) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.reminderCalendar[tgID] = false
	c.noteCalendar[tgID] = false
	c.noteSrv.SaveUser(tgID)
	c.reminderSrv.SaveUser(tgID)
	c.sharedSpace.SaveUser(tgID)

	return c.reminderSrv.CreateScheduler(ctx, tgID, loc, c.SendReminder)
}

// SaveState сохраняет в БД состояние бота с переданным пользователем
func (c *Controller) SaveState(ctx context.Context, tgID int64, state string) error {
	return c.userSrv.SaveState(ctx, tgID, state)
}

func (c *Controller) GetState(ctx context.Context, tgID int64) (string, error) {
	return c.userSrv.GetState(ctx, tgID)
}

// SaveMemoryReminder сохраняет напоминания, хранящиеся в памяти, в БД
func (c *Controller) SaveMemoryReminder(ctx context.Context) error {
	return c.reminderSrv.SaveMemory(ctx)
}
