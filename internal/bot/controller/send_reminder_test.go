package controller

import (
	"context"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = controller.SendReminder(ctx, randomReminder)
	assert.NoError(t, err)
}

func TestProcessDeleteReminder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chat := &tele.Chat{ID: 1}

	telectx := mocks.NewMockteleCtx(ctrl)
	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)

	controller := New(nil, nil, nil, reminderSrv)

	randomReminder := random.Reminder()

	reminderEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Return([]model.Reminder{*randomReminder}, nil).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	}).Times(2)

	reminderEditor.EXPECT().SaveJob(gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx interface{}, reminderID uuid.UUID, jobID uuid.UUID) {
		assert.Equal(t, randomReminder.ID, reminderID)
		randomReminder.Job.ID = jobID
	}).Return(nil).Times(2)

	reminderSrv.CreateScheduler(context.Background(), chat.ID, time.Local, controller.SendReminder)

	err := controller.reminderSrv.CreateScheduler(ctx, chat.ID, time.Local, controller.SendReminder)
	require.NoError(t, err)

	reminderEditor.EXPECT().DeleteReminderByID(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, reminderID uuid.UUID) {
		assert.Equal(t, randomReminder.ID, reminderID)
	}).Return(randomReminder.Job.ID, nil)

	expectedText := fmt.Sprintf(messages.ReminderDeletedMessage, randomReminder.Name)
	expectedSendOpts := &tele.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.BackToRemindersAndMenu(),
	}

	telectx.EXPECT().Edit(gomock.Any(), gomock.Any()).Do(func(msg string, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedText, msg)
		assert.Equal(t, expectedSendOpts, sendOpts)
	})

	err = controller.ProcessDeleteReminder(ctx, telectx, randomReminder)
	assert.NoError(t, err)
}
