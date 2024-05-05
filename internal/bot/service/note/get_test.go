package note

import (
	"context"
	"errors"
	"fmt"
	"testing"

	mock_note "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/mocks"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/Pizhlo/bot-reminder-go-telegram/pkg/random"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAll_Positive(t *testing.T) {
	type test struct {
		name         string
		userID       int64
		notesNum     int
		notes        []model.Note
		expectedText string
	}

	tests := []test{
		{
			name:     "one record",
			userID:   1,
			notesNum: 1,
		},
		{
			name:     "10 records",
			userID:   1,
			notesNum: 10,
		},
		{
			name:     "25 records",
			userID:   1,
			notesNum: 25,
		},
		{
			name:     "1000 records",
			userID:   1,
			notesNum: 1000,
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		noteEditor := mock_note.NewMocknoteEditor(ctrl)
		srv := New(noteEditor)
		view := view.NewNote()

		srv.SaveUser(tt.userID)

		tt.notes = random.Notes(tt.notesNum)

		tt.expectedText = view.Message(tt.notes)

		noteEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.All()).Return(tt.notes, nil)

		actualText, _, err := srv.GetAll(context.Background(), 1)
		require.NoError(t, err)

		assert.Equal(t, actualText, tt.expectedText, fmt.Sprintf("texts are not equal. Expected:======\n\n %s. Actual:======\n\n %s.", tt.expectedText, actualText))
	}
}

func TestGetAll_DBError(t *testing.T) {
	type test struct {
		name         string
		userID       int64
		notesNum     int
		notes        []model.Note
		expectedText string
		err          error
	}

	tests := []test{
		{
			name:     "one record",
			userID:   1,
			notesNum: 1,
		},
		{
			name:     "10 records",
			userID:   1,
			notesNum: 10,
		},
		{
			name:     "25 records",
			userID:   1,
			notesNum: 25,
		},
		{
			name:     "1000 records",
			userID:   1,
			notesNum: 1000,
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		noteEditor := mock_note.NewMocknoteEditor(ctrl)
		srv := New(noteEditor)
		view := view.NewNote()

		srv.SaveUser(tt.userID)

		tt.notes = random.Notes(tt.notesNum)
		tt.err = errors.New("test error")
		tt.expectedText = view.Message(tt.notes)

		noteEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.All()).Return(nil, tt.err)

		actualText, _, err := srv.GetAll(context.Background(), 1)

		assert.Equal(t, actualText, "")
		assert.Equal(t, err, tt.err)
	}
}

func TestNextPage_Positive(t *testing.T) {
	type test struct {
		name         string
		userID       int64
		notesNum     int
		notes        []model.Note
		expectedText string
	}

	tests := []test{
		{
			name:     "10 record",
			userID:   1,
			notesNum: 10,
		},
		{
			name:     "6 records",
			userID:   1,
			notesNum: 6,
		},
		{
			name:     "11 records",
			userID:   1,
			notesNum: 11,
		},
		{
			name:     "21 records",
			userID:   1,
			notesNum: 21,
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		noteEditor := mock_note.NewMocknoteEditor(ctrl)
		srv := New(noteEditor)
		view := view.NewNote()

		srv.SaveUser(tt.userID)

		tt.notes = random.Notes(tt.notesNum)
		view.Message(tt.notes)

		noteEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.All()).Return(tt.notes, nil)

		_, _, err := srv.GetAll(context.Background(), tt.userID)
		require.NoError(t, err)

		tt.expectedText = view.Next()

		actualText, _ := srv.NextPage(tt.userID)

		assert.Equal(t, actualText, tt.expectedText)
	}
}

func TestPrevPage_Positive(t *testing.T) {
	type test struct {
		name         string
		userID       int64
		notesNum     int
		notes        []model.Note
		expectedText string
	}

	tests := []test{
		{
			name:     "10 record",
			userID:   1,
			notesNum: 10,
		},
		{
			name:     "6 records",
			userID:   1,
			notesNum: 6,
		},
		{
			name:     "11 records",
			userID:   1,
			notesNum: 11,
		},
		{
			name:     "21 records",
			userID:   1,
			notesNum: 21,
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		noteEditor := mock_note.NewMocknoteEditor(ctrl)
		srv := New(noteEditor)
		view := view.NewNote()

		srv.SaveUser(tt.userID)

		tt.notes = random.Notes(tt.notesNum)
		view.Message(tt.notes)

		noteEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.All()).Return(tt.notes, nil)

		_, _, err := srv.GetAll(context.Background(), tt.userID)
		require.NoError(t, err)

		tt.expectedText = view.Previous()

		actualText, _ := srv.PrevPage(tt.userID)

		assert.Equal(t, actualText, tt.expectedText)
	}
}

func TestLastPage_Positive(t *testing.T) {
	type test struct {
		name         string
		userID       int64
		notesNum     int
		notes        []model.Note
		expectedText string
	}

	tests := []test{
		{
			name:     "10 record",
			userID:   1,
			notesNum: 10,
		},
		{
			name:     "6 records",
			userID:   1,
			notesNum: 6,
		},
		{
			name:     "11 records",
			userID:   1,
			notesNum: 11,
		},
		{
			name:     "21 records",
			userID:   1,
			notesNum: 21,
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		noteEditor := mock_note.NewMocknoteEditor(ctrl)
		srv := New(noteEditor)
		view := view.NewNote()

		srv.SaveUser(tt.userID)

		tt.notes = random.Notes(tt.notesNum)
		view.Message(tt.notes)

		noteEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.All()).Return(tt.notes, nil)

		_, _, err := srv.GetAll(context.Background(), tt.userID)
		require.NoError(t, err)

		tt.expectedText = view.Last()

		actualText, _ := srv.LastPage(tt.userID)

		assert.Equal(t, actualText, tt.expectedText)
	}
}

func TestFirstPage_Positive(t *testing.T) {
	type test struct {
		name         string
		userID       int64
		notesNum     int
		notes        []model.Note
		expectedText string
	}

	tests := []test{
		{
			name:     "10 record",
			userID:   1,
			notesNum: 10,
		},
		{
			name:     "6 records",
			userID:   1,
			notesNum: 6,
		},
		{
			name:     "11 records",
			userID:   1,
			notesNum: 11,
		},
		{
			name:     "21 records",
			userID:   1,
			notesNum: 21,
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		noteEditor := mock_note.NewMocknoteEditor(ctrl)
		srv := New(noteEditor)
		view := view.NewNote()

		srv.SaveUser(tt.userID)

		tt.notes = random.Notes(tt.notesNum)
		view.Message(tt.notes)

		noteEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.All()).Return(tt.notes, nil)

		_, _, err := srv.GetAll(context.Background(), tt.userID)
		require.NoError(t, err)

		tt.expectedText = view.First()

		actualText, _ := srv.FirstPage(tt.userID)

		assert.Equal(t, actualText, tt.expectedText)
	}
}
