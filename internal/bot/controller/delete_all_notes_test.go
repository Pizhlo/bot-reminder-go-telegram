package controller

import (
	"context"
	"testing"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/mocks"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/note"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	tele "gopkg.in/telebot.v3"
)

func TestDeleteAllNotes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	telectx := mocks.NewMockteleCtx(ctrl)
	noteEditor := mocks.NewMocknoteEditor(ctrl)

	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat)

	// n.noteEditor.DeleteAllByUserID(ctx, userID)
	noteEditor.EXPECT().DeleteAllByUserID(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	}).Return(nil)

	noteSrv := note.New(noteEditor)
	controller := New(nil, noteSrv, nil, nil)

	noteSrv.SaveUser(chat.ID)

	expectedText := messages.AllNotesDeletedMessage
	expectedKb := view.BackToMenuBtn()

	telectx.EXPECT().Edit(gomock.Any(), gomock.Any()).Do(func(text string, kb *tele.ReplyMarkup) {
		assert.Equal(t, expectedText, text)
		assert.Equal(t, expectedKb, kb)
	}).Return(nil)

	err := controller.DeleteAllNotes(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestConfirmDeleteAllNotes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	telectx := mocks.NewMockteleCtx(ctrl)

	controller := New(nil, nil, nil, nil)

	expectedText := messages.ConfirmDeleteNotesMessage
	expectedKb := &tele.ReplyMarkup{OneTimeKeyboard: true}
	expectedKb.Inline(
		expectedKb.Row(BtnDeleteAllNotes, BtnNotDeleteAllNotes),
	)

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, kb *tele.ReplyMarkup) {
		assert.Equal(t, expectedText, text)
		assert.Equal(t, expectedKb, kb)
	}).Return(nil)

	err := controller.ConfirmDeleteAllNotes(context.Background(), telectx)
	assert.NoError(t, err)
}
