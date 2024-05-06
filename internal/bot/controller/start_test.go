package controller

import (
	"context"
	"fmt"
	"testing"
	"time"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	mock_controller "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/mocks"
	user_model "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
	user_srv "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/user"
	cache "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/storage/cache/timezone"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/Pizhlo/bot-reminder-go-telegram/pkg/random"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/telebot.v3"
)

func TestStartCmd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	telectx := mock_controller.NewMockteleCtx(ctrl)

	chat := &telebot.Chat{ID: int64(1), FirstName: random.String(5)}

	telectx.EXPECT().Chat().Return(chat)

	expectedTxt := fmt.Sprintf(messages.StartMessage, chat.FirstName)
	expectedKb := view.MainMenu()

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, kb *telebot.ReplyMarkup) {
		assert.Equal(t, expectedTxt, text)
		assert.Equal(t, expectedKb, kb)
	})

	controller := New(nil, nil, nil, nil, 0)

	err := controller.StartCmd(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestMenuCmd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	telectx := mock_controller.NewMockteleCtx(ctrl)

	expectedTxt := messages.MenuMessage
	expectedKb := view.MainMenu()

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, kb *telebot.ReplyMarkup) {
		assert.Equal(t, expectedTxt, text)
		assert.Equal(t, expectedKb, kb)
	})

	controller := New(nil, nil, nil, nil, 0)

	err := controller.MenuCmd(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestHelpCmd_KnownUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	telectx := mock_controller.NewMockteleCtx(ctrl)

	tzCache := cache.New()
	db := mock_controller.NewMocktimezoneEditor(ctrl)
	userDB := mock_controller.NewMockuserEditor(ctrl)

	user := &telebot.User{FirstName: random.String(10)}

	db.EXPECT().GetAll(gomock.Any()).Return([]*user_model.User{{TGID: user.ID,
		Timezone: user_model.Timezone{Name: time.Local.String()}}}, nil)

	userSrv := user_srv.New(context.Background(), userDB, tzCache, db)

	userDB.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Do(func(ctx interface{}, userID int64, u *user_model.User) {
		assert.Equal(t, user.ID, userID)
	})

	err := userSrv.SaveUser(context.Background(), user.ID, &user_model.User{TGID: user.ID})
	require.NoError(t, err)

	expectedTxt := fmt.Sprintf(messages.HelpMessage, user.FirstName)
	expectedOpts := &telebot.SendOptions{
		ReplyMarkup: view.MainMenu(),
		ParseMode:   htmlParseMode,
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, kb *telebot.SendOptions) {
		assert.Equal(t, expectedTxt, text)
		assert.Equal(t, expectedOpts, kb)
	})

	chat := &telebot.Chat{ID: user.ID}

	telectx.EXPECT().Chat().Return(chat)
	telectx.EXPECT().Sender().Return(user)

	controller := New(userSrv, nil, nil, nil, 0)

	err = controller.HelpCmd(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestHelpCmd_UnknownUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	telectx := mock_controller.NewMockteleCtx(ctrl)

	tzCache := cache.New()
	db := mock_controller.NewMocktimezoneEditor(ctrl)
	userDB := mock_controller.NewMockuserEditor(ctrl)

	user := &telebot.User{FirstName: random.String(10)}

	db.EXPECT().GetAll(gomock.Any()).Return([]*user_model.User{}, nil)
	userDB.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, nil).Do(func(ctx interface{}, tgID int64) {
		assert.Equal(t, user.ID, tgID)
	})

	userSrv := user_srv.New(context.Background(), userDB, tzCache, db)

	expectedTxt := fmt.Sprintf(messages.HelpMessage, user.FirstName)
	expectedOpts := &telebot.SendOptions{
		ParseMode: htmlParseMode,
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, kb *telebot.SendOptions) {
		assert.Equal(t, expectedTxt, text)
		assert.Equal(t, expectedOpts, kb)
	})

	chat := &telebot.Chat{ID: user.ID}

	telectx.EXPECT().Chat().Return(chat)
	telectx.EXPECT().Sender().Return(user)

	controller := New(userSrv, nil, nil, nil, 0)

	err := controller.HelpCmd(context.Background(), telectx)
	assert.NoError(t, err)
}
