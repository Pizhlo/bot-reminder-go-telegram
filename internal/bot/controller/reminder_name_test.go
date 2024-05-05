package controller

import (
	"context"
	"fmt"
	"testing"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/mocks"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/reminder"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/Pizhlo/bot-reminder-go-telegram/pkg/random"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	tele "gopkg.in/telebot.v3"
)

func TestReminderName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	controller := New(nil, nil, nil, reminderSrv)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat).Times(2)

	randomName := random.String(5)

	telectx.EXPECT().Message().Return(&tele.Message{Sender: &tele.User{IsBot: false}, Text: randomName}).Times(2)

	expectedText := fmt.Sprintf(messages.TypeOfReminderMessage, randomName)
	expectedSendOpts := &tele.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.ReminderTypes(),
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedText, text)
		assert.Equal(t, expectedSendOpts, sendOpts)
	}).Return(nil)

	err := controller.ReminderName(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestReminderName_BotMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	controller := New(nil, nil, nil, reminderSrv)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat)

	randomText := random.String(5)
	randomName := random.String(5)

	reminderSrv.SaveName(chat.ID, randomName)

	telectx.EXPECT().Message().Return(&tele.Message{Sender: &tele.User{IsBot: true}, Text: randomText})

	expectedText := fmt.Sprintf(messages.TypeOfReminderMessage, randomName)
	expectedSendOpts := &tele.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.ReminderTypes(),
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedText, text)
		assert.Equal(t, expectedSendOpts, sendOpts)
	}).Return(nil)

	err := controller.ReminderName(context.Background(), telectx)
	assert.NoError(t, err)

	r, err := reminderSrv.GetFromMemory(chat.ID)
	assert.Equal(t, randomName, r.Name)
	assert.NoError(t, err)
}
