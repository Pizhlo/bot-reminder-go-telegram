package controller

import (
	"context"
	"fmt"
	"testing"
	"time"

	mock_controller "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller/mocks"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/mocks"
	model_user "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/user"
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
	tzEditor := mocks.NewMocktimezoneEditor(ctrl)
	userEditor := mocks.NewMockuserEditor(ctrl)
	tzCache := mocks.NewMocktimezoneCache(ctrl)

	loc := time.FixedZone("Europe/Moscow", 1)

	chat := &telebot.Chat{ID: int64(1), FirstName: random.String(5)}

	telectx.EXPECT().Chat().Return(chat)

	// при создании user service
	tzEditor.EXPECT().GetAll(gomock.Any()).Return([]*model_user.User{}, nil)

	// c.userSrv.GetLocation(ctx, telectx.Chat().ID)
	tzCache.EXPECT().Get(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, id int64) {
		assert.Equal(t, chat.ID, id)
	}).Return(loc, nil)

	expectedTxt := fmt.Sprintf(messages.TimezoneMessage, loc.String())

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, sendOpts *telebot.SendOptions) {
		assert.Equal(t, expectedTxt, text)

		expectedOpts := &telebot.SendOptions{
			ReplyMarkup: view.TimezoneMenu(),
			ParseMode:   htmlParseMode,
		}
		assert.Equal(t, sendOpts, expectedOpts)
	}).Return(nil)

	userSrv := user.New(context.Background(), userEditor, tzCache, tzEditor)

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
