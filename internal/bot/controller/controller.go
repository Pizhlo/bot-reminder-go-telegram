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
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/user"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type Controller struct {
	mu  sync.Mutex
	bot *tele.Bot
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
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
//go:generate mockgen -source ./controller.go -destination=./mocks/context.go
type teleCtx interface {
	// Bot returns the bot instance.
	Bot() *tele.Bot

	// Update returns the original update.
	Update() tele.Update

	// Message returns stored message if such presented.
	Message() *tele.Message

	// Callback returns stored callback if such presented.
	Callback() *tele.Callback

	// Query returns stored query if such presented.
	Query() *tele.Query

	// InlineResult returns stored inline result if such presented.
	InlineResult() *tele.InlineResult

	// ShippingQuery returns stored shipping query if such presented.
	ShippingQuery() *tele.ShippingQuery

	// PreCheckoutQuery returns stored pre checkout query if such presented.
	PreCheckoutQuery() *tele.PreCheckoutQuery

	// Poll returns stored poll if such presented.
	Poll() *tele.Poll

	// PollAnswer returns stored poll answer if such presented.
	PollAnswer() *tele.PollAnswer

	// ChatMember returns chat member changes.
	ChatMember() *tele.ChatMemberUpdate

	// ChatJoinRequest returns cha
	ChatJoinRequest() *tele.ChatJoinRequest

	// Migration returns both migration from and to chat IDs.
	Migration() (int64, int64)

	// Sender returns the current recipient, depending on the context type.
	// Returns nil if user is not presented.
	Sender() *tele.User

	// Chat returns the current chat, depending on the context type.
	// Returns nil if chat is not presented.
	Chat() *tele.Chat

	// Recipient combines both Sender and Chat functions. If there is no user
	// the chat will be returned. The native context cannot be without sender,
	// but it is useful in the case when the context created intentionally
	// by the NewContext constructor and have only Chat field inside.
	Recipient() tele.Recipient

	// Text returns the message text, depending on the context type.
	// In the case when no related data presented, returns an empty string.
	Text() string

	// Entities returns the message entities, whether it's media caption's or the text's.
	// In the case when no entities presented, returns a nil.
	Entities() tele.Entities

	// Data returns the current data, depending on the context type.
	// If the context contains command, returns its arguments string.
	// If the context contains payment, returns its payload.
	// In the case when no related data presented, returns an empty string.
	Data() string

	// Args returns a raw slice of command or callback arguments as strings.
	// The message arguments split by space, while the callback's ones by a "|" symbol.
	Args() []string

	// Send sends a message to the current recipient.
	// See Send from bot.go.
	Send(what interface{}, opts ...interface{}) error

	// SendAlbum sends an album to the current recipient.
	// See SendAlbum from bot.go.
	SendAlbum(a tele.Album, opts ...interface{}) error

	// Reply replies to the current message.
	// See Reply from bot.go.
	Reply(what interface{}, opts ...interface{}) error

	// Forward forwards the given message to the current recipient.
	// See Forward from bot.go.
	Forward(msg tele.Editable, opts ...interface{}) error

	// ForwardTo forwards the current message to the given recipient.
	// See Forward from bot.go
	ForwardTo(to tele.Recipient, opts ...interface{}) error

	// Edit edits the current message.
	// See Edit from bot.go.
	Edit(what interface{}, opts ...interface{}) error

	// EditCaption edits the caption of the current message.
	// See EditCaption from bot.go.
	EditCaption(caption string, opts ...interface{}) error

	// EditOrSend edits the current message if the update is callback,
	// otherwise the content is sent to the chat as a separate message.
	EditOrSend(what interface{}, opts ...interface{}) error

	// EditOrReply edits the current message if the update is callback,
	// otherwise the content is replied as a separate message.
	EditOrReply(what interface{}, opts ...interface{}) error

	// Delete removes the current message.
	// See Delete from bot.go.
	Delete() error

	// DeleteAfter waits for the duration to elapse and then removes the
	// message. It handles an error automatically using b.OnError callback.
	// It returns a Timer that can be used to cancel the call using its Stop method.
	DeleteAfter(d time.Duration) *time.Timer

	// Notify updates the chat action for the current recipient.
	// See Notify from bot.go.
	Notify(action tele.ChatAction) error

	// Ship replies to the current shipping query.
	// See Ship from bot.go.
	Ship(what ...interface{}) error

	// Accept finalizes the current deal.
	// See Accept from bot.go.
	Accept(errorMessage ...string) error

	// Answer sends a response to the current inline query.
	// See Answer from bot.go.
	Answer(resp *tele.QueryResponse) error

	// Respond sends a response for the current callback query.
	// See Respond from bot.go.
	Respond(resp ...*tele.CallbackResponse) error

	// Get retrieves data from the context.
	Get(key string) interface{}

	// Set saves data in the context.
	Set(key string, val interface{})
}

const (
	htmlParseMode     = "HTML"
	markdownParseMode = "markdown"
)

func New(userSrv *user.UserService, noteSrv *note.NoteService, bot *tele.Bot, reminderSrv *reminder.ReminderService) *Controller {
	return &Controller{
		userSrv:          userSrv,
		noteSrv:          noteSrv,
		bot:              bot,
		reminderSrv:      reminderSrv,
		schedulers:       make(map[int64]*gocron.Scheduler),
		noteCalendar:     make(map[int64]bool),
		reminderCalendar: make(map[int64]bool),
		mu:               sync.Mutex{}}
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
		msg = fmt.Sprintf(messages.ErrorMessageChannel, ctx.Message().Text, state, err)
	} else {
		msg = fmt.Sprintf(messages.ErrorMessageChannel, ctx.Callback().Unique, state, err)
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

// SaveUsers сохраняет пользователей в сервисах
func (c *Controller) SaveUsers(ctx context.Context, users []*user_model.User) {
	errors := []error{}
	for _, u := range users {
		c.noteSrv.SaveUser(u.TGID)

		c.reminderSrv.SaveUser(u.TGID)

		loc, err := c.userSrv.GetLocation(ctx, u.TGID)
		if err != nil {
			errors = append(errors, err)
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

func (c *Controller) saveUser(ctx context.Context, tgID int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.reminderCalendar[tgID] = false
	c.noteCalendar[tgID] = false
	c.noteSrv.SaveUser(tgID)
	c.reminderSrv.SaveUser(tgID)

	loc, err := c.userSrv.GetLocation(ctx, tgID)
	if err != nil {
		return err
	}

	return c.reminderSrv.CreateScheduler(ctx, tgID, loc, c.SendReminder)
}

// SaveState сохраняет в БД состояние бота с переданным пользователем
func (c *Controller) SaveState(ctx context.Context, tgID int64, state string) error {
	return c.userSrv.SaveState(ctx, tgID, state)
}

func (c *Controller) GetState(ctx context.Context, tgID int64) (string, error) {
	return c.userSrv.GetState(ctx, tgID)
}
