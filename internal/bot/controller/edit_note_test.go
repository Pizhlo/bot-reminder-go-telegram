package controller

import (
	"context"
	"fmt"
	"testing"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/mocks"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/note"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/Pizhlo/bot-reminder-go-telegram/pkg/random"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	tele "gopkg.in/telebot.v3"
)

func TestController_AskNoteText(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	telectx := mocks.NewMockteleCtx(ctrl)

	controller := New(nil, nil, nil, nil, 0)

	expectedText := messages.AskNewNoteTextMessage
	expectedKb := view.BackToMenuBtn()

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, kb *tele.ReplyMarkup) {
		assert.Equal(t, expectedText, text)
		assert.Equal(t, expectedKb, kb)
	}).Return(nil)

	err := controller.AskNoteText(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestController_UpdateNote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	telectx := mocks.NewMockteleCtx(ctrl)
	db := mocks.NewMocknoteEditor(ctrl)
	noteSrv := note.New(db)

	controller := New(nil, noteSrv, nil, nil, 0)

	expectedText := fmt.Sprintf(messages.EditNoteSuccessMessage, 1)
	expectedKb := view.BackToMenuAndNotesBtn()

	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat)

	msg := random.String(5)
	telectx.EXPECT().Message().Return(&tele.Message{Text: msg})

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, kb *tele.ReplyMarkup) {
		assert.Equal(t, expectedText, text)
		assert.Equal(t, expectedKb, kb)
	}).Return(nil)

	expectedNote := model.EditNote{
		TgID:   chat.ID,
		ViewID: 1,
		Text:   msg,
	}

	db.EXPECT().UpdateNote(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, note model.EditNote) {
		assert.Equal(t, expectedNote, note)
	}).Return(nil)

	err := controller.UpdateNote(context.Background(), telectx, 1)
	assert.NoError(t, err)
}
