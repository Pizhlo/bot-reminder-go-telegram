package controller

import (
	"context"
	"fmt"
	"testing"
	"time"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	mock_controller "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/mocks"
	model_user "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/note"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/reminder"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/user"
	tz_cache "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/storage/cache/timezone"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/Pizhlo/bot-reminder-go-telegram/pkg/random"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/telebot.v3"
)

func TestTimezone(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	telectx := mock_controller.NewMockteleCtx(ctrl)
	tzEditor := mock_controller.NewMocktimezoneEditor(ctrl)
	userEditor := mock_controller.NewMockuserEditor(ctrl)
	tz := tz_cache.New()

	loc := time.FixedZone("Europe/Moscow", 1)

	chat := &telebot.Chat{ID: int64(1), FirstName: random.String(5)}

	tz.Save(context.Background(), chat.ID, loc)

	telectx.EXPECT().Chat().Return(chat)

	// при создании user service
	tzEditor.EXPECT().GetAll(gomock.Any()).Return([]*model_user.User{}, nil)

	expectedTxt := fmt.Sprintf(messages.TimezoneMessage, loc.String())

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, sendOpts *telebot.SendOptions) {
		assert.Equal(t, expectedTxt, text)

		expectedOpts := &telebot.SendOptions{
			ReplyMarkup: view.TimezoneMenu(),
			ParseMode:   htmlParseMode,
		}
		assert.Equal(t, expectedOpts, sendOpts)
	}).Return(nil)

	userSrv := user.New(context.Background(), userEditor, tz, tzEditor)

	controller := New(userSrv, nil, nil, nil)

	err := controller.Timezone(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestRequestLocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	telectx := mock_controller.NewMockteleCtx(ctrl)

	chat := &telebot.Chat{ID: int64(1), FirstName: random.String(5)}
	telectx.EXPECT().Chat().Return(chat)

	expectedTxt := fmt.Sprintf(messages.StartMessageLocation, chat.FirstName)
	expectedKb := view.LocationMenu()

	telectx.EXPECT().Send(gomock.Any(), gomock.Any()).Do(func(text string, kb *telebot.ReplyMarkup) {
		assert.Equal(t, expectedTxt, text)
		assert.Equal(t, expectedKb, kb)
	})

	controller := New(nil, nil, nil, nil)

	err := controller.RequestLocation(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestAcceptTimezone(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// storage
	telectx := mock_controller.NewMockteleCtx(ctrl)
	tzEditor := mock_controller.NewMocktimezoneEditor(ctrl)
	userEditor := mock_controller.NewMockuserEditor(ctrl)
	reminderEditor := mock_controller.NewMockreminderEditor(ctrl)
	tz := tz_cache.New()

	// при создании user service
	tzEditor.EXPECT().GetAll(gomock.Any()).Return([]*model_user.User{}, nil)

	// srv
	userSrv := user.New(context.Background(), userEditor, tz, tzEditor)
	reminderSrv := reminder.New(reminderEditor)
	noteSrv := note.New(nil)

	chat := &telebot.Chat{
		ID: int64(1),
	}
	message := &telebot.Message{
		Location: &telebot.Location{
			Lat: 55.681717,
			Lng: 37.686791,
		},
	}

	user := &model_user.User{
		TGID: chat.ID,
		Timezone: model_user.Timezone{
			Name: "Europe/Moscow",
		},
	}

	telectx.EXPECT().Chat().Return(chat).Times(3)
	telectx.EXPECT().Message().Return(message).Times(2)

	// s.userEditor.Save(ctx, userID, u)
	userEditor.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64, u *model_user.User) {
		assert.Equal(t, chat.ID, userID)
		assert.Equal(t, user, u)
	}).Return(nil)

	// s.timezoneEditor.Save(ctx, userID, tz)
	tzEditor.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64, tz *model_user.Timezone) {
		assert.Equal(t, chat.ID, userID)
		assert.Equal(t, user.Timezone, *tz)
	}).Return(nil)

	// c.reminderEditor.GetAllByUserID(ctx, userID)
	reminderEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	}).Return(nil, api_errors.ErrRemindersNotFound)

	expectedTxt := fmt.Sprintf(messages.LocationMessage, user.Timezone.Name)

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, sendOpts *telebot.SendOptions) {
		assert.Equal(t, expectedTxt, text)

		expectedOpts := &telebot.SendOptions{
			ParseMode: htmlParseMode,
			ReplyMarkup: &telebot.ReplyMarkup{
				RemoveKeyboard: true,
			},
		}
		assert.Equal(t, expectedOpts, sendOpts)
	})

	controller := New(userSrv, noteSrv, nil, reminderSrv)

	err := controller.AcceptTimezone(context.Background(), telectx)
	assert.NoError(t, err)
}
