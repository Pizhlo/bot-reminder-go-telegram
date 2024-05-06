package controller

import (
	"context"
	"testing"
	"time"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/mocks"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	model_user "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/reminder"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/user"
	tz_cache "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/storage/cache/timezone"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/Pizhlo/bot-reminder-go-telegram/pkg/random"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	tele "gopkg.in/telebot.v3"
)

func TestReminderTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	tzEditor := mocks.NewMocktimezoneEditor(ctrl)
	userEditor := mocks.NewMockuserEditor(ctrl)
	tz := tz_cache.New()

	chat := &tele.Chat{ID: int64(1)}

	randomReminder := random.Reminder()

	// при создании user service
	tzEditor.EXPECT().GetAll(gomock.Any()).Return([]*model_user.User{
		{
			TGID: chat.ID,
			Timezone: model_user.Timezone{
				Name: "Europe/Moscow",
			},
		},
	}, nil)

	reminderEditor.EXPECT().Save(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, r *model.Reminder) {
		assert.Equal(t, randomReminder.TgID, r.TgID)
		assert.Equal(t, randomReminder.Name, r.Name)
		assert.Equal(t, randomReminder.Date, r.Date)
		assert.Equal(t, randomReminder.Type, r.Type)
		assert.Equal(t, randomReminder.Time, r.Time)
	}).Return(randomReminder.ID, nil)

	reminderEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	}).Return([]model.Reminder{}, nil)

	reminderEditor.EXPECT().SaveJob(gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx interface{}, reminderID uuid.UUID, jobID interface{}) {
		assert.Equal(t, randomReminder.ID, reminderID)
	}).Return(nil)

	userSrv := user.New(context.Background(), userEditor, tz, tzEditor)

	controller := New(userSrv, nil, nil, reminderSrv, 0)

	telectx := mocks.NewMockteleCtx(ctrl)

	telectx.EXPECT().Chat().Return(chat).Times(9)

	reminderSrv.SaveUser(chat.ID)
	reminderSrv.SaveName(chat.ID, randomReminder.Name)
	err := reminderSrv.SaveType(chat.ID, randomReminder.Type)
	assert.NoError(t, err)
	err = reminderSrv.SaveDate(chat.ID, randomReminder.Date)
	assert.NoError(t, err)

	telectx.EXPECT().Message().Return(&tele.Message{Text: randomReminder.Time}).Times(3)

	expectedOpts := &tele.SendOptions{
		ReplyMarkup: view.BackToMenuAndCreateOneElse(),
		ParseMode:   htmlParseMode,
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(msg interface{}, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedOpts, sendOpts)
	}).Return(nil)

	reminderEditor.EXPECT().DeleteMemory(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	}).Return(nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = controller.ReminderTime(ctx, telectx)
	assert.NoError(t, err)
}

func TestReminderTime_InvalidTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	tzEditor := mocks.NewMocktimezoneEditor(ctrl)
	userEditor := mocks.NewMockuserEditor(ctrl)
	tz := tz_cache.New()

	chat := &tele.Chat{ID: int64(1)}

	randomReminder := random.Reminder()

	// делаем время невалидным
	randomReminder.Time = random.String(10)

	// при создании user service
	tzEditor.EXPECT().GetAll(gomock.Any()).Return([]*model_user.User{
		{
			TGID: chat.ID,
			Timezone: model_user.Timezone{
				Name: "Europe/Moscow",
			},
		},
	}, nil)

	userSrv := user.New(context.Background(), userEditor, tz, tzEditor)

	controller := New(userSrv, nil, nil, reminderSrv, 0)

	telectx := mocks.NewMockteleCtx(ctrl)

	telectx.EXPECT().Chat().Return(chat)

	reminderSrv.SaveUser(chat.ID)
	reminderSrv.SaveName(chat.ID, randomReminder.Name)

	telectx.EXPECT().Message().Return(&tele.Message{Text: randomReminder.Time})

	err := controller.ReminderTime(context.Background(), telectx)
	assert.EqualError(t, err, api_errors.ErrInvalidTime.Error())
}

func TestReminderTime_TimeInPast(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	tzEditor := mocks.NewMocktimezoneEditor(ctrl)
	userEditor := mocks.NewMockuserEditor(ctrl)
	tz := tz_cache.New()

	chat := &tele.Chat{ID: int64(1)}

	randomReminder := random.Reminder()
	randomReminder.Type = model.DateType
	randomReminder.Date = time.Now().Format("02.01.2006")

	reminderSrv.SaveUser(chat.ID)
	reminderSrv.SaveName(chat.ID, randomReminder.Name)
	err := reminderSrv.SaveType(chat.ID, randomReminder.Type)
	assert.NoError(t, err)
	err = reminderSrv.SaveDate(chat.ID, randomReminder.Date)
	assert.NoError(t, err)

	// при создании user service
	tzEditor.EXPECT().GetAll(gomock.Any()).Return([]*model_user.User{
		{
			TGID: chat.ID,
			Timezone: model_user.Timezone{
				Name: "Europe/Moscow",
			},
		},
	}, nil)

	userSrv := user.New(context.Background(), userEditor, tz, tzEditor)

	controller := New(userSrv, nil, nil, reminderSrv, 0)

	telectx := mocks.NewMockteleCtx(ctrl)

	telectx.EXPECT().Chat().Return(chat).Times(3)

	telectx.EXPECT().Message().Return(&tele.Message{Text: "00:00"}).Times(2)

	err = controller.ReminderTime(context.Background(), telectx)
	assert.EqualError(t, err, api_errors.ErrTimeInPast.Error())
}
