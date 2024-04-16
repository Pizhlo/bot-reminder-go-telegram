package controller

import (
	"context"
	"fmt"
	"testing"
	"time"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/mocks"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	model_user "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/note"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/user"
	tz_cache "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/storage/cache/timezone"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/Pizhlo/bot-reminder-go-telegram/pkg/random"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	tele "gopkg.in/telebot.v3"
)

func TestSearchNoteByOnedate(t *testing.T) {
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

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(message string, kb *tele.ReplyMarkup) {
		assert.Equal(t, messages.SearchOneDateMessage, message)
		assert.Equal(t, noteView.Calendar(), kb)
	})

	err := controller.SearchNoteByOnedate(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestSearchNoteBySelectedDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mocks.NewMocknoteEditor(ctrl)
	tzEditor := mocks.NewMocktimezoneEditor(ctrl)
	tzCache := tz_cache.New()
	// при создании user service
	tzEditor.EXPECT().GetAll(gomock.Any()).Return([]*model_user.User{}, nil)

	userSrv := user.New(context.Background(), nil, tzCache, tzEditor)

	noteSrv := note.New(noteEditor)
	controller := New(userSrv, noteSrv, nil, nil)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat).Times(4)

	loc := time.FixedZone("Europe/Moscow", 1)
	tzCache.Save(context.Background(), chat.ID, loc)

	telectx.EXPECT().Callback().Return(&tele.Callback{Unique: "1"})

	noteSrv.SaveUser(chat.ID)
	noteSrv.SetupCalendar(chat.ID)

	noteView := view.NewNote()

	expectedNote := model.SearchByOneDate{
		TgID: chat.ID,
		Date: time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, loc),
	}

	note := random.Note()

	noteEditor.EXPECT().SearchByOneDate(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, searchNote model.SearchByOneDate) {
		assert.Equal(t, expectedNote, searchNote)
	}).Return([]model.Note{note}, nil)

	expectedText := noteView.Message([]model.Note{note})
	expectedSendOpts := &tele.SendOptions{
		ReplyMarkup: noteView.Keyboard(),
		ParseMode:   htmlParseMode,
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(message string, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedText, message)
		assert.Equal(t, expectedSendOpts, sendOpts)
	})

	err := controller.SearchNoteBySelectedDate(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestSearchNoteBySelectedDate_InvalidCallback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	controller := New(nil, nil, nil, nil)

	telectx := mocks.NewMockteleCtx(ctrl)

	callback := &tele.Callback{Unique: "abc"}
	telectx.EXPECT().Callback().Return(callback).Times(2)

	expectedErr := fmt.Errorf("error while converting string %s to type int: strconv.Atoi: parsing \"abc\": invalid syntax", callback.Unique)
	err := controller.SearchNoteBySelectedDate(context.Background(), telectx)
	assert.EqualError(t, err, expectedErr.Error())
}

func TestSearchNoteBySelectedDate_NoNotesFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mocks.NewMocknoteEditor(ctrl)
	tzEditor := mocks.NewMocktimezoneEditor(ctrl)
	tzCache := tz_cache.New()
	// при создании user service
	tzEditor.EXPECT().GetAll(gomock.Any()).Return([]*model_user.User{}, nil)

	userSrv := user.New(context.Background(), nil, tzCache, tzEditor)

	noteSrv := note.New(noteEditor)
	controller := New(userSrv, noteSrv, nil, nil)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat).Times(4)

	loc := time.FixedZone("Europe/Moscow", 1)
	tzCache.Save(context.Background(), chat.ID, loc)

	telectx.EXPECT().Callback().Return(&tele.Callback{Unique: "1"})

	noteSrv.SaveUser(chat.ID)
	noteSrv.SetupCalendar(chat.ID)

	expectedNote := model.SearchByOneDate{
		TgID: chat.ID,
		Date: time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, loc),
	}

	noteEditor.EXPECT().SearchByOneDate(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, searchNote model.SearchByOneDate) {
		assert.Equal(t, expectedNote, searchNote)
	}).Return([]model.Note{}, api_errors.ErrNotesNotFound)

	expectedText := fmt.Sprintf(messages.NoNotesFoundByDateMessage, expectedNote.Date.Format("02.01.2006"))
	expectedKb := view.BackToMenuAndNotesBtn()

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(message string, kb *tele.ReplyMarkup) {
		assert.Equal(t, expectedText, message)
		assert.Equal(t, expectedKb, kb)
	})

	err := controller.SearchNoteBySelectedDate(context.Background(), telectx)
	assert.NoError(t, err)
}
