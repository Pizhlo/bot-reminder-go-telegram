package note

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	mock_note "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/note/mocks"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/Pizhlo/bot-reminder-go-telegram/pkg/random"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
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
		view := view.New()

		srv.SaveUser(tt.userID)

		tt.notes = generateNotes(tt.notesNum)

		tt.expectedText = view.Message(tt.notes)

		noteEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.All()).Return(tt.notes, nil)

		actualText, _, err := srv.GetAll(context.Background(), 1)
		require.NoError(t, err)

		assert.Equal(t, actualText, tt.expectedText, fmt.Sprintf("texts are not equal. Expected:======\n\n %s. Actual:======\n\n %s.", tt.expectedText, actualText))
	}
}

// generateNotes генерирует необходимое количество заметок
func generateNotes(n int) []model.Note {
	var notes []model.Note

	for i := 0; i < n; i++ {
		note := model.Note{
			ID:      i,
			TgID:    1,
			Text:    random.String(10),
			Created: time.Now(),
		}
		notes = append(notes, note)
	}

	return notes
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
		view := view.New()

		srv.SaveUser(tt.userID)

		tt.notes = generateNotes(tt.notesNum)
		tt.err = sql.ErrNoRows
		tt.expectedText = view.Message(tt.notes)

		noteEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.All()).Return(nil, tt.err)

		actualText, _, err := srv.GetAll(context.Background(), 1)

		assert.Equal(t, actualText, "")
		assert.Equal(t, err, sql.ErrNoRows)
	}
}
