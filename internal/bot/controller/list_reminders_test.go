package controller

import (
	"context"
	"testing"
	"time"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/mocks"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/reminder"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/Pizhlo/bot-reminder-go-telegram/pkg/random"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	tele "gopkg.in/telebot.v3"
)

func TestListReminders_Positive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := mocks.NewMockreminderEditor(ctrl)
	srv := reminder.New(db)
	controller := New(nil, nil, nil, srv, 0)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat).Times(3)

	srv.SaveUser(chat.ID)

	reminders := random.Reminders(5)

	db.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	}).Return(reminders, nil).Times(2)

	db.EXPECT().SaveJob(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(len(reminders)).Do(func(ctx interface{}, reminderID uuid.UUID, jobID uuid.UUID) {
		for i := 0; i < len(reminders); i++ {
			if reminders[i].ID == reminderID {
				reminders[i].Job.ID = jobID
			}
		}
	})

	err := srv.CreateScheduler(ctx, chat.ID, time.Local, func(ctx context.Context, reminder *model.Reminder) error { return nil })
	assert.NoError(t, err)

	view := view.NewReminder()

	expectedText, err := view.Message(reminders)
	assert.NoError(t, err)

	expectedSendOptions := &tele.SendOptions{
		ReplyMarkup: view.Keyboard(),
		ParseMode:   htmlParseMode,
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedText, text)
		assert.Equal(t, expectedSendOptions, sendOpts)
	}).Return(nil)

	err = controller.ListReminders(ctx, telectx)
	assert.NoError(t, err)
}

func TestListReminders_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	controller := New(nil, nil, nil, reminderSrv, 0)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat)

	reminderSrv.SaveUser(chat.ID)

	reminderEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	}).Return(nil, api_errors.ErrRemindersNotFound)

	expectedText := messages.NoRemindersMessage
	expectedKb := view.CreateReminderAndBackToMenu()

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, kb *tele.ReplyMarkup) {
		assert.Equal(t, expectedText, text)
		assert.Equal(t, expectedKb, kb)
	}).Return(nil)

	err := controller.ListReminders(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestNextPageReminders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	controller := New(nil, nil, nil, reminderSrv, 0)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat).Times(2)

	reminderSrv.SaveUser(chat.ID)

	reminders := random.Reminders(5)

	reminderEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	}).Return(reminders, nil)

	// сформировать сообщения во вью
	reminderSrv.GetAll(context.Background(), chat.ID)
	reminderSrv.Message(chat.ID, reminders)

	reminderView := view.NewReminder()
	reminderView.Message(reminders)

	expectedText := reminderView.Next()
	expectedOpts := &tele.SendOptions{
		ReplyMarkup: reminderView.Keyboard(),
		ParseMode:   htmlParseMode,
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedText, text)
		assert.Equal(t, expectedOpts, sendOpts)
	}).Return(nil)

	err := controller.NextPageReminders(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestNextPageReminders_MessageNotEdited(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	controller := New(nil, nil, nil, reminderSrv, 0)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat).Times(2)

	reminderSrv.SaveUser(chat.ID)

	reminders := random.Reminders(5)

	reminderEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	}).Return(reminders, nil)

	// сформировать сообщения во вью
	reminderSrv.GetAll(context.Background(), chat.ID)
	reminderSrv.Message(chat.ID, reminders)

	reminderView := view.NewReminder()
	reminderView.Message(reminders)

	expectedText := reminderView.Next()
	expectedOpts := &tele.SendOptions{
		ReplyMarkup: reminderView.Keyboard(),
		ParseMode:   htmlParseMode,
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedText, text)
		assert.Equal(t, expectedOpts, sendOpts)
	}).Return(tele.ErrMessageNotModified)

	err := controller.NextPageReminders(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestPrevPageReminders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	controller := New(nil, nil, nil, reminderSrv, 0)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat).Times(2)

	reminderSrv.SaveUser(chat.ID)

	reminders := random.Reminders(5)

	reminderEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	}).Return(reminders, nil)

	// сформировать сообщения во вью
	reminderSrv.GetAll(context.Background(), chat.ID)
	reminderSrv.Message(chat.ID, reminders)

	reminderView := view.NewReminder()
	reminderView.Message(reminders)

	expectedText := reminderView.Previous()
	expectedOpts := &tele.SendOptions{
		ReplyMarkup: reminderView.Keyboard(),
		ParseMode:   htmlParseMode,
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedText, text)
		assert.Equal(t, expectedOpts, sendOpts)
	}).Return(nil)

	err := controller.PrevPageReminders(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestPrevPageReminders_MessageNotModified(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	controller := New(nil, nil, nil, reminderSrv, 0)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat).Times(2)

	reminderSrv.SaveUser(chat.ID)

	reminders := random.Reminders(5)

	reminderEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	}).Return(reminders, nil)

	// сформировать сообщения во вью
	reminderSrv.GetAll(context.Background(), chat.ID)
	reminderSrv.Message(chat.ID, reminders)

	reminderView := view.NewReminder()
	reminderView.Message(reminders)

	expectedText := reminderView.Previous()
	expectedOpts := &tele.SendOptions{
		ReplyMarkup: reminderView.Keyboard(),
		ParseMode:   htmlParseMode,
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedText, text)
		assert.Equal(t, expectedOpts, sendOpts)
	}).Return(tele.ErrMessageNotModified)

	err := controller.PrevPageReminders(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestLastPageReminders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	controller := New(nil, nil, nil, reminderSrv, 0)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat).Times(2)

	reminderSrv.SaveUser(chat.ID)

	reminders := random.Reminders(5)

	reminderEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	}).Return(reminders, nil)

	// сформировать сообщения во вью
	reminderSrv.GetAll(context.Background(), chat.ID)
	reminderSrv.Message(chat.ID, reminders)

	reminderView := view.NewReminder()
	reminderView.Message(reminders)

	expectedText := reminderView.Last()
	expectedOpts := &tele.SendOptions{
		ReplyMarkup: reminderView.Keyboard(),
		ParseMode:   htmlParseMode,
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedText, text)
		assert.Equal(t, expectedOpts, sendOpts)
	}).Return(nil)

	err := controller.LastPageReminders(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestLastPageReminders_MessageNotModified(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	controller := New(nil, nil, nil, reminderSrv, 0)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat).Times(2)

	reminderSrv.SaveUser(chat.ID)

	reminders := random.Reminders(5)

	reminderEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	}).Return(reminders, nil)

	// сформировать сообщения во вью
	reminderSrv.GetAll(context.Background(), chat.ID)
	reminderSrv.Message(chat.ID, reminders)

	reminderView := view.NewReminder()
	reminderView.Message(reminders)

	expectedText := reminderView.Last()
	expectedOpts := &tele.SendOptions{
		ReplyMarkup: reminderView.Keyboard(),
		ParseMode:   htmlParseMode,
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedText, text)
		assert.Equal(t, expectedOpts, sendOpts)
	}).Return(tele.ErrMessageNotModified)

	err := controller.LastPageReminders(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestFirstPageReminders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	controller := New(nil, nil, nil, reminderSrv, 0)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat).Times(2)

	reminderSrv.SaveUser(chat.ID)

	reminders := random.Reminders(5)

	reminderEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	}).Return(reminders, nil)

	// сформировать сообщения во вью
	reminderSrv.GetAll(context.Background(), chat.ID)
	reminderSrv.Message(chat.ID, reminders)

	reminderView := view.NewReminder()
	reminderView.Message(reminders)

	expectedText := reminderView.First()
	expectedOpts := &tele.SendOptions{
		ReplyMarkup: reminderView.Keyboard(),
		ParseMode:   htmlParseMode,
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedText, text)
		assert.Equal(t, expectedOpts, sendOpts)
	}).Return(nil)

	err := controller.FirstPageReminders(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestFirstPageReminders_MessageNotModified(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	controller := New(nil, nil, nil, reminderSrv, 0)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat).Times(2)

	reminderSrv.SaveUser(chat.ID)

	reminders := random.Reminders(5)

	reminderEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	}).Return(reminders, nil)

	// сформировать сообщения во вью
	reminderSrv.GetAll(context.Background(), chat.ID)
	reminderSrv.Message(chat.ID, reminders)

	reminderView := view.NewReminder()
	reminderView.Message(reminders)

	expectedText := reminderView.First()
	expectedOpts := &tele.SendOptions{
		ReplyMarkup: reminderView.Keyboard(),
		ParseMode:   htmlParseMode,
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedText, text)
		assert.Equal(t, expectedOpts, sendOpts)
	}).Return(tele.ErrMessageNotModified)

	err := controller.FirstPageReminders(context.Background(), telectx)
	assert.NoError(t, err)
}
