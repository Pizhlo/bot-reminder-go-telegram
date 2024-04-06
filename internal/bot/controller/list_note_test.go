package controller

import (
	"context"
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

func TestListNotes_ErrNotesNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mocks.NewMocknoteEditor(ctrl)
	noteSrv := note.New(noteEditor)
	controller := New(nil, noteSrv, nil, nil)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat)

	telectx.EXPECT().Message().Return(&tele.Message{}).Times(2)

	// n.noteEditor.GetAllByUserID(ctx, userID)
	noteEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	}).Return(nil, api_errors.ErrNotesNotFound)

	expectedText := messages.UserDoesntHaveNotesMessage
	expectedKb := view.BackToMenuBtn()

	telectx.EXPECT().Edit(gomock.Any(), gomock.Any()).Do(func(text string, kb *tele.ReplyMarkup) {
		assert.Equal(t, expectedText, text)
		assert.Equal(t, expectedKb, kb)
	}).Return(nil)

	err := controller.ListNotes(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestListNotes_Positive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mocks.NewMocknoteEditor(ctrl)
	noteSrv := note.New(noteEditor)
	controller := New(nil, noteSrv, nil, nil)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat)

	noteSrv.SaveUser(chat.ID)

	telectx.EXPECT().Message().Return(&tele.Message{}).Times(2)

	notes := random.Notes(5)

	// n.noteEditor.GetAllByUserID(ctx, userID)
	noteEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	}).Return(notes, nil)

	noteView := view.NewNote()

	expectedText := noteView.Message(notes)
	expectedSendOptions := &tele.SendOptions{
		ReplyMarkup: noteView.Keyboard(),
		ParseMode:   htmlParseMode,
	}

	telectx.EXPECT().Edit(gomock.Any(), gomock.Any()).Do(func(text string, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedText, text)
		assert.Equal(t, expectedSendOptions, sendOpts)
	}).Return(nil)

	err := controller.ListNotes(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestNextPageNotes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mocks.NewMocknoteEditor(ctrl)
	noteSrv := note.New(noteEditor)
	controller := New(nil, noteSrv, nil, nil)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat)

	noteSrv.SaveUser(chat.ID)

	notes := random.Notes(5)

	// n.noteEditor.GetAllByUserID(ctx, userID)
	noteEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	}).Return(notes, nil)

	noteView := view.NewNote()

	noteSrv.GetAll(context.Background(), 1)
	noteView.Message(notes)

	expectedText := noteView.Next()
	expectedSendOptions := &tele.SendOptions{
		ReplyMarkup: noteView.Keyboard(),
		ParseMode:   htmlParseMode,
	}

	telectx.EXPECT().Edit(gomock.Any(), gomock.Any()).Do(func(text string, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedText, text)
		assert.Equal(t, expectedSendOptions, sendOpts)
	}).Return(nil)

	err := controller.NextPageNotes(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestPrevPageNotes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mocks.NewMocknoteEditor(ctrl)
	noteSrv := note.New(noteEditor)
	controller := New(nil, noteSrv, nil, nil)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat)

	noteSrv.SaveUser(chat.ID)

	notes := random.Notes(5)

	// n.noteEditor.GetAllByUserID(ctx, userID)
	noteEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	}).Return(notes, nil)

	noteView := view.NewNote()

	noteSrv.GetAll(context.Background(), 1)
	noteView.Message(notes)

	expectedText := noteView.Previous()
	expectedSendOptions := &tele.SendOptions{
		ReplyMarkup: noteView.Keyboard(),
		ParseMode:   htmlParseMode,
	}

	telectx.EXPECT().Edit(gomock.Any(), gomock.Any()).Do(func(text string, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedText, text)
		assert.Equal(t, expectedSendOptions, sendOpts)
	}).Return(nil)

	err := controller.PrevPageNotes(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestLastPageNotes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mocks.NewMocknoteEditor(ctrl)
	noteSrv := note.New(noteEditor)
	controller := New(nil, noteSrv, nil, nil)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat)

	noteSrv.SaveUser(chat.ID)

	notes := random.Notes(5)

	// n.noteEditor.GetAllByUserID(ctx, userID)
	noteEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	}).Return(notes, nil)

	noteView := view.NewNote()

	noteSrv.GetAll(context.Background(), 1)
	noteView.Message(notes)

	expectedText := noteView.Last()
	expectedSendOptions := &tele.SendOptions{
		ReplyMarkup: noteView.Keyboard(),
		ParseMode:   htmlParseMode,
	}

	telectx.EXPECT().Edit(gomock.Any(), gomock.Any()).Do(func(text string, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedText, text)
		assert.Equal(t, expectedSendOptions, sendOpts)
	}).Return(nil)

	err := controller.LastPageNotes(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestFirstPageNotes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mocks.NewMocknoteEditor(ctrl)
	noteSrv := note.New(noteEditor)
	controller := New(nil, noteSrv, nil, nil)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat)

	noteSrv.SaveUser(chat.ID)

	notes := random.Notes(5)

	// n.noteEditor.GetAllByUserID(ctx, userID)
	noteEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	}).Return(notes, nil)

	noteView := view.NewNote()

	noteSrv.GetAll(context.Background(), 1)
	noteView.Message(notes)

	expectedText := noteView.First()
	expectedSendOptions := &tele.SendOptions{
		ReplyMarkup: noteView.Keyboard(),
		ParseMode:   htmlParseMode,
	}

	telectx.EXPECT().Edit(gomock.Any(), gomock.Any()).Do(func(text string, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedText, text)
		assert.Equal(t, expectedSendOptions, sendOpts)
	}).Return(nil)

	err := controller.FirstPageNotes(context.Background(), telectx)
	assert.NoError(t, err)
}
