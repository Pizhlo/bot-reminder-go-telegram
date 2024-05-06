package controller

import (
	"context"
	"fmt"
	"strconv"
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

func TestSearchNoteByTwoDates(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mocks.NewMocknoteEditor(ctrl)
	noteSrv := note.New(noteEditor)
	controller := New(nil, noteSrv, nil, nil, 0)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat)

	noteSrv.SaveUser(chat.ID)

	noteView := view.NewNote()

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(message string, kb *tele.ReplyMarkup) {
		assert.Equal(t, messages.SearchByTwoDatesFirstDateMessage, message)
		assert.Equal(t, noteView.Calendar(), kb)
	})

	err := controller.SearchNoteByTwoDates(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestSearchNoteByTwoDatesFirstDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mocks.NewMocknoteEditor(ctrl)
	tzEditor := mocks.NewMocktimezoneEditor(ctrl)
	tzCache := tz_cache.New()
	// при создании user service
	tzEditor.EXPECT().GetAll(gomock.Any()).Return([]*model_user.User{}, nil)

	userSrv := user.New(context.Background(), nil, tzCache, tzEditor)

	noteSrv := note.New(noteEditor)
	controller := New(userSrv, noteSrv, nil, nil, 0)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat).Times(5)

	loc := time.FixedZone("Europe/Moscow", 1)
	tzCache.Save(context.Background(), chat.ID, loc)

	telectx.EXPECT().Callback().Return(&tele.Callback{Unique: "1"})

	noteSrv.SaveUser(chat.ID)
	noteSrv.SetupCalendar(chat.ID)

	noteView := view.NewNote()
	noteView.SetCurMonth()
	noteView.SetCurYear()

	expectedText := messages.SearchByTwoDatesSecondDateMessage
	expectedKb := noteView.Calendar()

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(message string, kb *tele.ReplyMarkup) {
		assert.Equal(t, expectedText, message)
		assert.Equal(t, expectedKb, kb)
	})

	err := controller.SearchNoteByTwoDatesFirstDate(context.Background(), telectx)
	assert.NoError(t, err)

	expectedNote := &model.SearchByTwoDates{
		TgID:      chat.ID,
		FirstDate: time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, loc),
	}
	n, err := noteSrv.GetSearchNote(chat.ID)
	assert.NoError(t, err)
	assert.Equal(t, expectedNote, n)
}

func TestSearchNoteByTwoDatesFirstDate_InvalidCallback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	controller := New(nil, nil, nil, nil, 0)

	telectx := mocks.NewMockteleCtx(ctrl)

	callback := &tele.Callback{Unique: "abc"}
	telectx.EXPECT().Callback().Return(callback).Times(2)

	expectedErr := fmt.Errorf("error while converting string %s to type int: strconv.Atoi: parsing \"abc\": invalid syntax", callback.Unique)
	err := controller.SearchNoteByTwoDatesFirstDate(context.Background(), telectx)
	assert.EqualError(t, err, expectedErr.Error())
}

func TestSearchNoteByTwoDatesFirstDate_FirstDayAfterToday(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mocks.NewMocknoteEditor(ctrl)
	tzEditor := mocks.NewMocktimezoneEditor(ctrl)
	tzCache := tz_cache.New()
	// при создании user service
	tzEditor.EXPECT().GetAll(gomock.Any()).Return([]*model_user.User{}, nil)

	userSrv := user.New(context.Background(), nil, tzCache, tzEditor)

	noteSrv := note.New(noteEditor)
	controller := New(userSrv, noteSrv, nil, nil, 0)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat).Times(3)

	noteSrv.SaveUser(chat.ID)
	noteSrv.SetupCalendar(chat.ID)

	loc := time.FixedZone("Europe/Moscow", 1)
	tzCache.Save(context.Background(), chat.ID, loc)

	callback := &tele.Callback{Unique: strconv.Itoa(time.Now().Day() + 1)}
	telectx.EXPECT().Callback().Return(callback)

	err := controller.SearchNoteByTwoDatesFirstDate(context.Background(), telectx)
	assert.EqualError(t, err, api_errors.ErrFirstDayFuture.Error())
}

func TestSearchNoteByTwoDatesSecondDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mocks.NewMocknoteEditor(ctrl)
	tzEditor := mocks.NewMocktimezoneEditor(ctrl)
	tzCache := tz_cache.New()
	// при создании user service
	tzEditor.EXPECT().GetAll(gomock.Any()).Return([]*model_user.User{}, nil)

	userSrv := user.New(context.Background(), nil, tzCache, tzEditor)

	noteSrv := note.New(noteEditor)
	controller := New(userSrv, noteSrv, nil, nil, 0)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat).Times(7)

	loc := time.FixedZone("Europe/Moscow", 1)
	tzCache.Save(context.Background(), chat.ID, loc)

	telectx.EXPECT().Callback().Return(&tele.Callback{Unique: "2"})

	noteSrv.SaveUser(chat.ID)
	noteSrv.SetupCalendar(chat.ID)

	noteView := view.NewNote()
	noteView.SetCurMonth()
	noteView.SetCurYear()

	firstDate := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, loc)
	noteSrv.SaveFirstDate(chat.ID, firstDate)

	expectedNote := &model.SearchByTwoDates{
		TgID:       chat.ID,
		FirstDate:  firstDate,
		SecondDate: time.Date(time.Now().Year(), time.Now().Month(), 2, 0, 0, 0, 0, loc),
	}

	notes := random.Notes(2)
	noteEditor.EXPECT().SearchByTwoDates(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, searchNote *model.SearchByTwoDates) {
		assert.Equal(t, expectedNote, searchNote)
	}).Return(notes, nil)

	expectedText := noteView.Message(notes)
	expectedOpts := &tele.SendOptions{
		ReplyMarkup: noteView.Keyboard(),
		ParseMode:   htmlParseMode,
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(message string, opts *tele.SendOptions) {
		assert.Equal(t, expectedText, message)
		assert.Equal(t, expectedOpts, opts)
	})

	err := controller.SearchNoteByTwoDatesSecondDate(context.Background(), telectx)
	assert.NoError(t, err)

	n, err := noteSrv.GetSearchNote(chat.ID)
	assert.NoError(t, err)
	assert.Equal(t, expectedNote, n)
}

func TestSearchNoteByTwoDatesSecondDate_InvalidCallback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	controller := New(nil, nil, nil, nil, 0)

	telectx := mocks.NewMockteleCtx(ctrl)

	callback := &tele.Callback{Unique: "abc"}
	telectx.EXPECT().Callback().Return(callback).Times(2)

	expectedErr := fmt.Errorf("error while converting string %s to type int: strconv.Atoi: parsing \"abc\": invalid syntax", callback.Unique)
	err := controller.SearchNoteByTwoDatesSecondDate(context.Background(), telectx)
	assert.EqualError(t, err, expectedErr.Error())
}

func TestSearchNoteByTwoDatesSecondDate_NoNotesFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mocks.NewMocknoteEditor(ctrl)
	tzEditor := mocks.NewMocktimezoneEditor(ctrl)
	tzCache := tz_cache.New()
	// при создании user service
	tzEditor.EXPECT().GetAll(gomock.Any()).Return([]*model_user.User{}, nil)

	userSrv := user.New(context.Background(), nil, tzCache, tzEditor)

	noteSrv := note.New(noteEditor)
	controller := New(userSrv, noteSrv, nil, nil, 0)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat).Times(7)

	loc := time.FixedZone("Europe/Moscow", 1)
	tzCache.Save(context.Background(), chat.ID, loc)

	telectx.EXPECT().Callback().Return(&tele.Callback{Unique: "2"})

	noteSrv.SaveUser(chat.ID)
	noteSrv.SetupCalendar(chat.ID)

	noteView := view.NewNote()
	noteView.SetCurMonth()
	noteView.SetCurYear()

	firstDate := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, loc)
	noteSrv.SaveFirstDate(chat.ID, firstDate)

	expectedNote := &model.SearchByTwoDates{
		TgID:       chat.ID,
		FirstDate:  firstDate,
		SecondDate: time.Date(time.Now().Year(), time.Now().Month(), 2, 0, 0, 0, 0, loc),
	}

	noteEditor.EXPECT().SearchByTwoDates(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, searchNote *model.SearchByTwoDates) {
		assert.Equal(t, expectedNote, searchNote)
	}).Return(nil, api_errors.ErrNotesNotFound)

	expectedText := fmt.Sprintf(messages.NoNotesFoundByTwoDatesMessage, expectedNote.FirstDate.Format("02.01.2006"), expectedNote.SecondDate.Format("02.01.2006"))
	expectedKb := view.BackToMenuAndNotesBtn()

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(message string, kb *tele.ReplyMarkup) {
		assert.Equal(t, expectedText, message)
		assert.Equal(t, expectedKb, kb)
	})

	err := controller.SearchNoteByTwoDatesSecondDate(context.Background(), telectx)
	assert.NoError(t, err)

	n, err := noteSrv.GetSearchNote(chat.ID)
	assert.NoError(t, err)
	assert.Equal(t, expectedNote, n)
}

func TestSearchNoteByTwoDatesSecondDate_ErrSecondDateBeforeFirst(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mocks.NewMocknoteEditor(ctrl)
	tzEditor := mocks.NewMocktimezoneEditor(ctrl)
	tzCache := tz_cache.New()
	// при создании user service
	tzEditor.EXPECT().GetAll(gomock.Any()).Return([]*model_user.User{}, nil)

	userSrv := user.New(context.Background(), nil, tzCache, tzEditor)

	noteSrv := note.New(noteEditor)
	controller := New(userSrv, noteSrv, nil, nil, 0)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat).Times(4)

	loc := time.FixedZone("Europe/Moscow", 1)
	tzCache.Save(context.Background(), chat.ID, loc)

	telectx.EXPECT().Callback().Return(&tele.Callback{Unique: "1"})

	noteSrv.SaveUser(chat.ID)
	noteSrv.SetupCalendar(chat.ID)

	firstDate := time.Date(time.Now().Year(), time.Now().Month(), 2, 0, 0, 0, 0, loc)
	noteSrv.SaveFirstDate(chat.ID, firstDate)

	err := controller.SearchNoteByTwoDatesSecondDate(context.Background(), telectx)
	assert.EqualError(t, err, api_errors.ErrSecondDateBeforeFirst.Error())
}

func TestSearchNoteByTwoDatesSecondDate_ErrSecondDateFuture(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mocks.NewMocknoteEditor(ctrl)
	tzEditor := mocks.NewMocktimezoneEditor(ctrl)
	tzCache := tz_cache.New()
	// при создании user service
	tzEditor.EXPECT().GetAll(gomock.Any()).Return([]*model_user.User{}, nil)

	userSrv := user.New(context.Background(), nil, tzCache, tzEditor)

	noteSrv := note.New(noteEditor)
	controller := New(userSrv, noteSrv, nil, nil, 0)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat).Times(4)

	loc := time.FixedZone("Europe/Moscow", 1)
	tzCache.Save(context.Background(), chat.ID, loc)

	telectx.EXPECT().Callback().Return(&tele.Callback{Unique: strconv.Itoa(time.Now().Day() + 1)})

	noteSrv.SaveUser(chat.ID)
	noteSrv.SetupCalendar(chat.ID)

	firstDate := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, loc)
	noteSrv.SaveFirstDate(chat.ID, firstDate)

	err := controller.SearchNoteByTwoDatesSecondDate(context.Background(), telectx)
	assert.EqualError(t, err, api_errors.ErrSecondDateFuture.Error())
}

func TestSecondDateBeforeFirst(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mocks.NewMocknoteEditor(ctrl)

	noteSrv := note.New(noteEditor)
	controller := New(nil, noteSrv, nil, nil, 0)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat)

	noteSrv.SaveUser(chat.ID)
	noteSrv.SetupCalendar(chat.ID)

	noteView := view.NewNote()
	noteView.SetCurMonth()
	noteView.SetCurYear()

	expectedText := messages.FirstDateBeforeSecondMessage
	expectedOpts := &tele.SendOptions{
		ReplyMarkup: noteView.Calendar(),
		ParseMode:   htmlParseMode,
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(message string, opts *tele.SendOptions) {
		assert.Equal(t, expectedText, message)
		assert.Equal(t, expectedOpts, opts)
	})

	err := controller.SecondDateBeforeFirst(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestSecondDateInFuture(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mocks.NewMocknoteEditor(ctrl)

	noteSrv := note.New(noteEditor)
	controller := New(nil, noteSrv, nil, nil, 0)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat)

	noteSrv.SaveUser(chat.ID)
	noteSrv.SetupCalendar(chat.ID)

	noteView := view.NewNote()
	noteView.SetCurMonth()
	noteView.SetCurYear()

	expectedText := messages.SecondDateInFutureMessage
	expectedOpts := &tele.SendOptions{
		ReplyMarkup: noteView.Calendar(),
		ParseMode:   htmlParseMode,
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(message string, opts *tele.SendOptions) {
		assert.Equal(t, expectedText, message)
		assert.Equal(t, expectedOpts, opts)
	})

	err := controller.SecondDateInFuture(context.Background(), telectx)
	assert.NoError(t, err)
}

func TestFirstDateInFuture(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mocks.NewMocknoteEditor(ctrl)

	noteSrv := note.New(noteEditor)
	controller := New(nil, noteSrv, nil, nil, 0)

	telectx := mocks.NewMockteleCtx(ctrl)
	chat := &tele.Chat{ID: int64(1)}
	telectx.EXPECT().Chat().Return(chat)

	noteSrv.SaveUser(chat.ID)
	noteSrv.SetupCalendar(chat.ID)

	noteView := view.NewNote()
	noteView.SetCurMonth()
	noteView.SetCurYear()

	expectedText := messages.FirstDateInFutureMessage
	expectedOpts := &tele.SendOptions{
		ReplyMarkup: noteView.Calendar(),
		ParseMode:   htmlParseMode,
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(message string, opts *tele.SendOptions) {
		assert.Equal(t, expectedText, message)
		assert.Equal(t, expectedOpts, opts)
	})

	err := controller.FirstDateInFuture(context.Background(), telectx)
	assert.NoError(t, err)
}
