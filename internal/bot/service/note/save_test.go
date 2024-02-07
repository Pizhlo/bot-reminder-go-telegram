package note

import (
	"context"
	"testing"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	mock_note "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/note/mocks"
	"github.com/Pizhlo/bot-reminder-go-telegram/pkg/random"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestSave_Positive(t *testing.T) {
	type test struct {
		name   string
		userID int64
		note   model.Note
	}

	tests := []test{
		{
			name:   "one record",
			userID: 1,
			note:   random.Note(),
		},
		{
			name:   "10 records",
			userID: 1,
			note:   random.Note(),
		},
		{
			name:   "25 records",
			userID: 1,
			note:   random.Note(),
		},
		{
			name:   "1000 records",
			userID: 1,
			note:   random.Note(),
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		noteEditor := mock_note.NewMocknoteEditor(ctrl)
		srv := New(noteEditor)

		noteEditor.EXPECT().Save(gomock.Any(), gomock.All()).Return(nil)

		err := srv.Save(context.Background(), tt.note)
		require.NoError(t, err)
	}
}
