package controller

import (
	"context"
	"testing"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/mocks"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/reminder"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	tele "gopkg.in/telebot.v3"
)

func TestConfirmDeleteAllReminders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	telectx := mocks.NewMockteleCtx(ctrl)
	reminderEditor := mocks.NewMockreminderEditor(ctrl)

	reminderSrv := reminder.New(reminderEditor)
	controller := New(nil, nil, nil, reminderSrv, 0)

	expectedText := messages.ConfirmDeleteRemindersMessage
	expectedKb := &tele.ReplyMarkup{OneTimeKeyboard: true}
	expectedKb.Inline(
		expectedKb.Row(BtnDeleteAllReminders, BtnNotDeleteAllReminders),
	)

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, kb *tele.ReplyMarkup) {
		assert.Equal(t, expectedText, text)
		assert.Equal(t, expectedKb, kb)
	}).Return(nil)

	err := controller.ConfirmDeleteAllReminders(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestDeleteAllReminders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	telectx := mocks.NewMockteleCtx(ctrl)
	reminderEditor := mocks.NewMockreminderEditor(ctrl)

	reminderSrv := reminder.New(reminderEditor)
	controller := New(nil, nil, nil, reminderSrv, 0)

	chat := &tele.Chat{
		ID: int64(1),
	}

	telectx.EXPECT().Chat().Return(chat)

	reminderEditor.EXPECT().GetAllJobs(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	}).Return([]uuid.UUID{}, nil)

	expectedText := messages.AllRemindersDeletedMessage
	expectedKb := view.BackToMenuBtn()

	telectx.EXPECT().Edit(gomock.Any(), gomock.Any()).Do(func(text string, kb *tele.ReplyMarkup) {
		assert.Equal(t, expectedText, text)
		assert.Equal(t, expectedKb, kb)
	}).Return(nil)

	err := controller.DeleteAllReminders(context.Background(), telectx)
	assert.NoError(t, err)
}
