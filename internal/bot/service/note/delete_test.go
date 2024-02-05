package note

import (
	"context"
	"database/sql"
	"testing"

	mock_note "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/note/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteAll_Positive(t *testing.T) {
	type test struct {
		userID int64
	}

	tests := []test{
		{
			userID: 1,
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		noteEditor := mock_note.NewMocknoteEditor(ctrl)
		srv := New(noteEditor)

		srv.SaveUser(tt.userID)

		noteEditor.EXPECT().DeleteAllByUserID(gomock.Any(), gomock.All()).Return(nil)

		err := srv.DeleteAll(context.Background(), 1)
		require.NoError(t, err)
	}
}

func TestDeleteAll_DbErr(t *testing.T) {
	type test struct {
		userID int64
		err    error
	}

	tests := []test{
		{
			userID: 1,
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		noteEditor := mock_note.NewMocknoteEditor(ctrl)
		srv := New(noteEditor)

		srv.SaveUser(tt.userID)

		tt.err = sql.ErrNoRows

		noteEditor.EXPECT().DeleteAllByUserID(gomock.Any(), gomock.All()).Return(tt.err)

		err := srv.DeleteAll(context.Background(), 1)
		assert.EqualError(t, err, tt.err.Error())
	}
}
