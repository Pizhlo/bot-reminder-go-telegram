package note

import (
	"context"
	"testing"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	mock_note "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/note/mocks"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/Pizhlo/bot-reminder-go-telegram/pkg/random"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchByText_Positive(t *testing.T) {
	type test struct {
		name         string
		userID       int64
		notesNum     int
		notes        []model.Note
		searchNote   model.SearchNoteByText
		expectedText string
	}

	tests := []test{
		{
			name:   "10 record",
			userID: 1,
			searchNote: model.SearchNoteByText{
				TgID: 1,
				Text: "test",
			},
			notesNum: 10,
		},
		{
			name:   "6 records",
			userID: 1,
			searchNote: model.SearchNoteByText{
				TgID: 1,
				Text: "test",
			},
			notesNum: 6,
		},
		{
			name:   "11 records",
			userID: 1,
			searchNote: model.SearchNoteByText{
				TgID: 1,
				Text: "test",
			},
			notesNum: 11,
		},
		{
			name:   "21 records",
			userID: 1,
			searchNote: model.SearchNoteByText{
				TgID: 1,
				Text: "test",
			},
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

		noteEditor.EXPECT().SearchByText(gomock.Any(), gomock.All()).Return(tt.notes, nil)

		actualText, _, err := srv.SearchByText(context.Background(), tt.searchNote)
		require.NoError(t, err)

		tt.expectedText = view.Message(tt.notes)

		assert.Equal(t, actualText, tt.expectedText)
	}
}

func TestSearchByText_NotesNotFound(t *testing.T) {
	type test struct {
		name       string
		userID     int64
		notesNum   int
		notes      []model.Note
		searchNote model.SearchNoteByText
		err        error
	}

	tests := []test{
		{
			name:   "10 record",
			userID: 1,
			searchNote: model.SearchNoteByText{
				TgID: 1,
				Text: "test",
			},
			notesNum: 10,
			err:      api_errors.ErrNotesNotFound,
		},
		{
			name:   "6 records",
			userID: 1,
			searchNote: model.SearchNoteByText{
				TgID: 1,
				Text: "test",
			},
			notesNum: 6,
			err:      api_errors.ErrNotesNotFound,
		},
		{
			name:   "11 records",
			userID: 1,
			searchNote: model.SearchNoteByText{
				TgID: 1,
				Text: "test",
			},
			notesNum: 11,
			err:      api_errors.ErrNotesNotFound,
		},
		{
			name:   "21 records",
			userID: 1,
			searchNote: model.SearchNoteByText{
				TgID: 1,
				Text: "test",
			},
			notesNum: 21,
			err:      api_errors.ErrNotesNotFound,
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		noteEditor := mock_note.NewMocknoteEditor(ctrl)
		srv := New(noteEditor)

		srv.SaveUser(tt.userID)

		tt.notes = random.Notes(tt.notesNum)

		noteEditor.EXPECT().SearchByText(gomock.Any(), gomock.All()).Return(nil, tt.err)

		_, _, err := srv.SearchByText(context.Background(), tt.searchNote)
		assert.EqualError(t, err, tt.err.Error())
	}
}
