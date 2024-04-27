package controller

import (
	"context"
	"fmt"
	"testing"
	"time"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
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
	"gopkg.in/telebot.v3"
)

func TestSaveReminder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	telectx := mocks.NewMockteleCtx(ctrl)
	tzEditor := mocks.NewMocktimezoneEditor(ctrl)
	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)

	// при создании user service
	tzEditor.EXPECT().GetAll(gomock.Any()).Return([]*model_user.User{}, nil)

	tzCache := tz_cache.New()
	userSrv := user.New(context.Background(), nil, tzCache, tzEditor)

	randomReminder := random.Reminder()

	controller := New(userSrv, nil, nil, reminderSrv)

	loc := time.FixedZone("Europe/Moscow", 1)

	chat := &telebot.Chat{ID: int64(1), FirstName: random.String(5)}

	telectx.EXPECT().Chat().Return(chat).Times(5)

	tzCache.Save(context.Background(), chat.ID, loc)

	reminderEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Return([]model.Reminder{}, nil)
	reminderEditor.EXPECT().SaveJob(gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx interface{}, reminderID uuid.UUID, jobID uuid.UUID) {
		assert.Equal(t, randomReminder.ID, reminderID, "reminder IDs are not equal")
	}).Times(2).Return(nil)

	reminderSrv.SaveName(chat.ID, randomReminder.Name)
	err := reminderSrv.SaveDate(chat.ID, randomReminder.Date)
	assert.NoError(t, err)
	err = reminderSrv.SaveTime(chat.ID, randomReminder.Time)
	assert.NoError(t, err)
	err = reminderSrv.SaveType(chat.ID, randomReminder.Type)
	assert.NoError(t, err)
	err = reminderSrv.SaveCreatedField(chat.ID, loc)
	assert.NoError(t, err)

	var verb string
	if randomReminder.Type == model.DateType {
		verb = "сработает"
	} else {
		verb = "будет срабатывать"
	}

	reminderEditor.EXPECT().Save(gomock.Any(), gomock.Any()).Return(randomReminder.ID, nil).Do(func(ctx interface{}, reminder *model.Reminder) {
		assert.Equal(t, randomReminder.ID, reminder.ID)
		assert.Equal(t, randomReminder.TgID, reminder.TgID)
		assert.Equal(t, randomReminder.Name, reminder.Name)
		assert.Equal(t, randomReminder.Type, reminder.Type)
		assert.Equal(t, randomReminder.Date, reminder.Date)
		assert.Equal(t, randomReminder.Time, reminder.Time)
		assert.Equal(t, randomReminder.Created, reminder.Created)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	nextRun, err := reminderSrv.SaveAndStartReminder(ctx, chat.ID, loc, controller.SendReminder, randomReminder)
	assert.NoError(t, err)

	nextRunMsg, err := view.ProcessTypeAndDate(randomReminder.Type, randomReminder.Date, randomReminder.Time)
	assert.NoError(t, err)

	layout := "02.01.2006 15:04:05"

	expectedTxt := fmt.Sprintf(messages.SuccessCreationMessage, randomReminder.Name, verb, nextRunMsg, nextRun.NextRun.Format(layout))
	expectedOpts := &telebot.SendOptions{
		ReplyMarkup: view.BackToMenuAndCreateOneElse(),
		ParseMode:   htmlParseMode,
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, sendOpts *telebot.SendOptions) {
		assert.Equal(t, expectedTxt, text)
		assert.Equal(t, expectedOpts, sendOpts)
	}).Return(nil)

	reminderEditor.EXPECT().Save(gomock.Any(), gomock.Any()).Return(randomReminder.ID, nil).Do(func(ctx interface{}, reminder *model.Reminder) {
		assert.Equal(t, randomReminder.TgID, reminder.TgID)
		assert.Equal(t, randomReminder.Date, reminder.Date)
		assert.Equal(t, randomReminder.Time, reminder.Time)
		assert.Equal(t, randomReminder.Type, reminder.Type)
		assert.Equal(t, randomReminder.Name, reminder.Name)
	})
	reminderEditor.EXPECT().DeleteMemory(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	}).Return(nil)

	err = controller.saveReminder(ctx, telectx)
	assert.NoError(t, err)
}
