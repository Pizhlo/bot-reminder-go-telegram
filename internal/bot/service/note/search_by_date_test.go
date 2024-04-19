package note

import (
	"context"
	"testing"
	"time"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	mock_note "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/mocks"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/Pizhlo/bot-reminder-go-telegram/pkg/random"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchByOneDate_Positive(t *testing.T) {
	userID := int64(1)
	searchNote := model.SearchByOneDate{
		TgID: 1,
		Date: time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, time.Local),
	}

	notes := random.Notes(5)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mock_note.NewMocknoteEditor(ctrl)
	srv := New(noteEditor)
	view := view.NewNote()

	srv.SaveUser(userID)

	noteEditor.EXPECT().SearchByOneDate(gomock.Any(), gomock.All()).Return(notes, nil)

	actualText, _, err := srv.SearchByOneDate(context.Background(), searchNote)
	require.NoError(t, err)

	expectedText := view.Message(notes)

	assert.Equal(t, actualText, expectedText)
}

func TestSearchByOneDate_NotFound(t *testing.T) {
	userID := int64(1)
	searchNote := model.SearchByOneDate{
		TgID: 1,
		Date: time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, time.Local),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mock_note.NewMocknoteEditor(ctrl)
	srv := New(noteEditor)

	srv.SaveUser(userID)

	noteEditor.EXPECT().SearchByOneDate(gomock.Any(), gomock.All()).Return(nil, api_errors.ErrNotesNotFound)

	actualText, _, err := srv.SearchByOneDate(context.Background(), searchNote)
	assert.EqualError(t, err, api_errors.ErrNotesNotFound.Error())

	assert.Equal(t, "", actualText)
}

func TestSaveFirstDate(t *testing.T) {
	userID := int64(1)
	date := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, time.Local)

	srv := New(nil)

	srv.SaveFirstDate(userID, date)

	val, ok := srv.searchMap[userID]
	assert.Equal(t, true, ok)

	assert.Equal(t, userID, val.TgID)
	assert.Equal(t, date, val.FirstDate)
	assert.Equal(t, true, val.SecondDate.IsZero())
}

func TestSaveSecondtDate_NotFound(t *testing.T) {
	userID := int64(1)
	date := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, time.Local)

	srv := New(nil)

	err := srv.SaveSecondDate(userID, date)
	assert.EqualError(t, err, "Note service: no data found for this user")
}

func TestSaveSecondtDate_Positive(t *testing.T) {
	userID := int64(1)
	date := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, time.Local)

	srv := New(nil)

	srv.SaveFirstDate(userID, date)

	err := srv.SaveSecondDate(userID, date)
	assert.NoError(t, err)

	val, ok := srv.searchMap[userID]
	assert.Equal(t, true, ok)

	assert.Equal(t, userID, val.TgID)
	assert.Equal(t, date, val.FirstDate)
	assert.Equal(t, date, val.SecondDate)
}

func TestGetSearchNote_NotFound(t *testing.T) {
	userID := int64(1)

	srv := New(nil)

	_, err := srv.GetSearchNote(userID)
	assert.EqualError(t, err, "Note service: no data found for this user")
}

func TestGetSearchNote_Positive(t *testing.T) {
	userID := int64(1)
	date := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, time.Local)

	srv := New(nil)

	srv.SaveFirstDate(userID, date)

	err := srv.SaveSecondDate(userID, date)
	assert.NoError(t, err)

	val, err := srv.GetSearchNote(userID)
	assert.NoError(t, err)

	assert.Equal(t, userID, val.TgID)
	assert.Equal(t, date, val.FirstDate)
	assert.Equal(t, date, val.SecondDate)
}

func TestValidateSearchDate_Positive(t *testing.T) {
	userID := int64(1)
	date1 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-1, time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, time.Local)
	date2 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, time.Local)

	srv := New(nil)

	srv.SaveFirstDate(userID, date1)

	err := srv.ValidateSearchDate(userID, date2)
	assert.NoError(t, err)
}

func TestValidateSearchDate_Invalid(t *testing.T) {
	userID := int64(1)
	date1 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, time.Local)
	date2 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-1, time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, time.Local)

	srv := New(nil)

	srv.SaveFirstDate(userID, date1)

	err := srv.ValidateSearchDate(userID, date2)
	assert.EqualError(t, err, api_errors.ErrSecondDateBeforeFirst.Error())
}

func TestSearchByTwoDates_Positive(t *testing.T) {
	userID := int64(1)
	date := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, time.Local)
	notes := random.Notes(5)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mock_note.NewMocknoteEditor(ctrl)
	srv := New(noteEditor)

	srv.SaveUser(userID)

	srv.SaveFirstDate(userID, date)

	err := srv.SaveSecondDate(userID, date)
	assert.NoError(t, err)

	noteEditor.EXPECT().SearchByTwoDates(gomock.Any(), gomock.All()).Return(notes, nil)

	actualText, _, err := srv.SearchByTwoDates(context.Background(), userID)
	assert.NoError(t, err)

	view := view.NewNote()

	expectedText := view.Message(notes)

	assert.Equal(t, actualText, expectedText)
}

func TestSearchByTwoDates_NotFoundNotes(t *testing.T) {
	userID := int64(1)
	date := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, time.Local)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mock_note.NewMocknoteEditor(ctrl)
	srv := New(noteEditor)

	srv.SaveUser(userID)

	srv.SaveFirstDate(userID, date)

	err := srv.SaveSecondDate(userID, date)
	assert.NoError(t, err)

	noteEditor.EXPECT().SearchByTwoDates(gomock.Any(), gomock.All()).Return(nil, api_errors.ErrNotesNotFound)

	_, _, err = srv.SearchByTwoDates(context.Background(), userID)
	assert.EqualError(t, err, api_errors.ErrNotesNotFound.Error())
}

func TestSearchByTwoDates_NotFoundInMap(t *testing.T) {
	userID := int64(1)

	srv := New(nil)

	srv.SaveUser(userID)

	_, _, err := srv.SearchByTwoDates(context.Background(), userID)
	assert.EqualError(t, err, "Note service: no data found for this user")
}
