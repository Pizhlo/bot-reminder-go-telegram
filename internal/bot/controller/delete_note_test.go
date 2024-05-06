package controller

import (
	"context"
	"fmt"
	"testing"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/mocks"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/note"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/Pizhlo/bot-reminder-go-telegram/pkg/random"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	tele "gopkg.in/telebot.v3"
)

func TestDeleteNoteByID_Positive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mocks.NewMocknoteEditor(ctrl)
	noteSrv := note.New(noteEditor)
	controller := New(nil, noteSrv, nil, nil, 0)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat)

	telectx.EXPECT().Message().Return(&tele.Message{Text: "/dn1"})

	note := random.Note()

	// n.noteEditor.GetByViewID(ctx, userID, noteID)
	noteEditor.EXPECT().GetByViewID(gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64, noteID int) {
		assert.Equal(t, chat.ID, userID)
		assert.Equal(t, 1, noteID)
	}).Return(&note, nil)

	// n.noteEditor.DeleteNoteByViewID(ctx, userID, noteID)
	noteEditor.EXPECT().DeleteNoteByViewID(gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64, noteID int) {
		assert.Equal(t, chat.ID, userID)
		assert.Equal(t, 1, noteID)
	}).Return(nil)

	expectedText := fmt.Sprintf(messages.NoteDeletedSuccessMessage, 1)
	expectedSendOpts := &tele.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.NotesAndMenuBtns(),
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedText, text)
		assert.Equal(t, expectedSendOpts, sendOpts)
	}).Return(nil)

	err := controller.DeleteNoteByID(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestDeleteNoteByID_NotFoundSuffix(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mocks.NewMocknoteEditor(ctrl)
	noteSrv := note.New(noteEditor)
	controller := New(nil, noteSrv, nil, nil, 0)

	telectx := mocks.NewMockteleCtx(ctrl)

	msg := &tele.Message{Text: random.String(5)}

	telectx.EXPECT().Message().Return(msg)

	expectedErr := fmt.Errorf("error in controller.DeleteNoteByID(): not found suffix %s in message text: %s", deleteNotePrefix, msg.Text)

	err := controller.DeleteNoteByID(context.Background(), telectx)
	assert.EqualError(t, err, expectedErr.Error())
}

func TestDeleteNoteByID_IdNotInt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mocks.NewMocknoteEditor(ctrl)
	noteSrv := note.New(noteEditor)
	controller := New(nil, noteSrv, nil, nil, 0)

	telectx := mocks.NewMockteleCtx(ctrl)

	rnd := random.String(5)
	msg := &tele.Message{Text: "/dn" + rnd}

	telectx.EXPECT().Message().Return(msg)

	errTxt := fmt.Sprintf("strconv.Atoi: parsing \"%s\": invalid syntax", rnd)

	expectedErr := fmt.Errorf("error while convertion string %s to int while handling command %s: %s", rnd, msg.Text, errTxt)

	err := controller.DeleteNoteByID(context.Background(), telectx)
	assert.EqualError(t, err, expectedErr.Error())
}

func TestDeleteNoteByID_ErrNotesNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mocks.NewMocknoteEditor(ctrl)
	noteSrv := note.New(noteEditor)
	controller := New(nil, noteSrv, nil, nil, 0)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat)

	telectx.EXPECT().Message().Return(&tele.Message{Text: "/dn1"})

	// n.noteEditor.GetByViewID(ctx, userID, noteID)
	noteEditor.EXPECT().GetByViewID(gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64, noteID int) {
		assert.Equal(t, chat.ID, userID)
		assert.Equal(t, 1, noteID)
	}).Return(nil, api_errors.ErrNotesNotFound)

	expectedText := fmt.Sprintf(messages.NoNoteFoundByNumberMessage, 1)
	expectedKb := view.BackToMenuBtn()

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, kb *tele.ReplyMarkup) {
		assert.Equal(t, expectedText, text)
		assert.Equal(t, expectedKb, kb)
	}).Return(nil)

	err := controller.DeleteNoteByID(context.Background(), telectx)
	assert.NoError(t, err)
}
