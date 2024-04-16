package controller

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/mocks"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/reminder"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/Pizhlo/bot-reminder-go-telegram/pkg/random"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tele "gopkg.in/telebot.v3"
)

func TestSendReminder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)

	randomReminder := random.Reminder()

	bot, err := tele.NewBot(tele.Settings{Offline: true})
	require.NoError(t, err)

	controller := New(nil, nil, bot, reminderSrv)

	err = controller.SendReminder(context.Background(), &randomReminder)
	assert.True(t, err == tele.ErrNotFound)
}

func TestProcessDeleteReminder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	telectx := mocks.NewMockteleCtx(ctrl)
	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)

	controller := New(nil, nil, nil, reminderSrv)

	randomReminder := random.Reminder()
	chat := &tele.Chat{ID: 1}

	reminderEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Return([]model.Reminder{randomReminder}, nil).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	})

	reminderEditor.EXPECT().SaveJob(gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx interface{}, reminderID uuid.UUID, jobID uuid.UUID) {
		assert.Equal(t, randomReminder.ID, reminderID)
		randomReminder.Job.ID = jobID
	}).Return(nil)

	err := controller.reminderSrv.CreateScheduler(context.Background(), chat.ID, time.Local, controller.SendReminder)
	require.NoError(t, err)

	telectx.EXPECT().Callback().Return(&tele.Callback{Unique: fmt.Sprintf("%d", randomReminder.ViewID)})
	telectx.EXPECT().Chat().Return(chat)

	reminderEditor.EXPECT().GetByViewID(gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64, viewID int) {
		assert.Equal(t, chat.ID, userID)
		assert.Equal(t, int(randomReminder.ViewID), viewID)
	}).Return(&randomReminder, nil)

	reminderEditor.EXPECT().DeleteReminderByID(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, reminderID uuid.UUID) {
		assert.Equal(t, randomReminder.ID, reminderID)
	}).Return(nil)

	expectedText := fmt.Sprintf(messages.ReminderDeletedMessage, randomReminder.Name)
	expectedSendOpts := &tele.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.BackToMenuBtn(),
	}

	telectx.EXPECT().Edit(gomock.Any(), gomock.Any()).Do(func(msg string, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedText, msg)
		assert.Equal(t, expectedSendOpts, sendOpts)
	})

	err = controller.ProcessDeleteReminder(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestProcessDeleteReminder_ReminderDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	telectx := mocks.NewMockteleCtx(ctrl)
	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)

	controller := New(nil, nil, nil, reminderSrv)

	randomReminder := random.Reminder()
	chat := &tele.Chat{ID: 1}

	reminderEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Return([]model.Reminder{randomReminder}, nil).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	})

	reminderEditor.EXPECT().SaveJob(gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx interface{}, reminderID uuid.UUID, jobID uuid.UUID) {
		assert.Equal(t, randomReminder.ID, reminderID)
		randomReminder.Job.ID = jobID
	}).Return(nil)

	err := controller.reminderSrv.CreateScheduler(context.Background(), chat.ID, time.Local, controller.SendReminder)
	require.NoError(t, err)

	telectx.EXPECT().Callback().Return(&tele.Callback{Unique: fmt.Sprintf("%d", randomReminder.ViewID)})
	telectx.EXPECT().Chat().Return(chat)

	reminderEditor.EXPECT().GetByViewID(gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64, viewID int) {
		assert.Equal(t, chat.ID, userID)
		assert.Equal(t, int(randomReminder.ViewID), viewID)
	}).Return(nil, sql.ErrNoRows)

	expectedText := fmt.Sprintf(messages.ReminderDeletedMessage, "")
	expectedSendOpts := &tele.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.BackToMenuBtn(),
	}

	telectx.EXPECT().Edit(gomock.Any(), gomock.Any()).Do(func(msg string, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedText, msg)
		assert.Equal(t, expectedSendOpts, sendOpts)
	})

	err = controller.ProcessDeleteReminder(context.Background(), telectx)
	assert.NoError(t, err)
}
