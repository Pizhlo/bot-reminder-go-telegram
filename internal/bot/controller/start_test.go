package controller

import (
	"context"
	"fmt"
	"testing"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	mock_controller "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/mocks"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/Pizhlo/bot-reminder-go-telegram/pkg/random"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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

	controller := New(nil, nil, nil, nil)

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

	controller := New(nil, nil, nil, nil)

	err := controller.MenuCmd(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestHelpCmd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	telectx := mock_controller.NewMockteleCtx(ctrl)

	expectedTxt := messages.HelpMessage
	expectedKb := view.MainMenu()

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, kb *telebot.ReplyMarkup) {
		assert.Equal(t, expectedTxt, text)
		assert.Equal(t, expectedKb, kb)
	})

	controller := New(nil, nil, nil, nil)

	err := controller.HelpCmd(context.Background(), telectx)
	assert.NoError(t, err)
}
