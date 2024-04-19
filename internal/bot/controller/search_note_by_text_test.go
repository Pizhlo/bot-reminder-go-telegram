package controller

import (
	"context"
	"testing"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
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

func TestSearchNoteByText(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mocks.NewMocknoteEditor(ctrl)
	noteSrv := note.New(noteEditor)
	controller := New(nil, noteSrv, nil, nil)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat)

	noteSrv.SaveUser(chat.ID)

	noteView := view.NewNote()

	text := random.String(5)

	note := random.Note()
	note.Text = text

	telectx.EXPECT().Message().Return(&tele.Message{Text: text})

	noteEditor.EXPECT().SearchByText(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, searchNote model.SearchByText) {
		assert.Equal(t, chat.ID, searchNote.TgID)
		assert.Equal(t, text, searchNote.Text)
	}).Return([]model.Note{note}, nil)

	expectedTxt := noteView.Message([]model.Note{note})
	expectedSendOpts := &tele.SendOptions{
		ReplyMarkup: noteView.Keyboard(),
		ParseMode:   htmlParseMode,
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(message string, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedTxt, message)
		assert.Equal(t, expectedSendOpts, sendOpts)
	})

	err := controller.SearchNoteByText(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestSearchNoteByText_NoNotesFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mocks.NewMocknoteEditor(ctrl)
	noteSrv := note.New(noteEditor)
	controller := New(nil, noteSrv, nil, nil)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat)

	noteSrv.SaveUser(chat.ID)

	text := random.String(5)

	telectx.EXPECT().Message().Return(&tele.Message{Text: text})

	noteEditor.EXPECT().SearchByText(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, searchNote model.SearchByText) {
		assert.Equal(t, chat.ID, searchNote.TgID)
		assert.Equal(t, text, searchNote.Text)
	}).Return([]model.Note{}, api_errors.ErrNotesNotFound)

	expectedTxt := messages.NoNotesFoundByTextMessage
	expectedKb := view.BackToMenuBtn()

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(message string, kb *tele.ReplyMarkup) {
		assert.Equal(t, expectedTxt, message)
		assert.Equal(t, expectedKb, kb)
	})

	err := controller.SearchNoteByText(context.Background(), telectx)
	assert.NoError(t, err)
}
