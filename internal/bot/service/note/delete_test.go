package note

import (
	"context"
	"errors"
	"testing"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"

	mock_note "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/mocks"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDeleteAll_Positive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mock_note.NewMocknoteEditor(ctrl)
	srv := New(noteEditor)

	srv.SaveUser(1)

	noteEditor.EXPECT().DeleteAllByUserID(gomock.Any(), gomock.All()).Return(nil)

	err := srv.DeleteAll(context.Background(), 1)

	assert.NoError(t, err)
}

func TestDeleteAll_DBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mock_note.NewMocknoteEditor(ctrl)
	srv := New(noteEditor)

	srv.SaveUser(1)

	testErr := errors.New("test error")

	noteEditor.EXPECT().DeleteAllByUserID(gomock.Any(), gomock.All()).Return(testErr)

	err := srv.DeleteAll(context.Background(), 1)

	assert.Equal(t, err, testErr)
}

func TestDeleteByID_Positive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mock_note.NewMocknoteEditor(ctrl)
	srv := New(noteEditor)

	srv.SaveUser(1)

	noteEditor.EXPECT().GetByViewID(gomock.Any(), gomock.All(), gomock.Any()).Return(&model.Note{}, nil)
	noteEditor.EXPECT().DeleteByID(gomock.Any(), gomock.All()).Return(nil)

	err := srv.DeleteByID(context.Background(), 1, 1)

	assert.NoError(t, err)
}

func TestDeleteByID_NotesNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noteEditor := mock_note.NewMocknoteEditor(ctrl)
	srv := New(noteEditor)

	srv.SaveUser(1)

	noteEditor.EXPECT().GetByViewID(gomock.Any(), gomock.All(), gomock.Any()).Return(nil, api_errors.ErrNotesNotFound)

	err := srv.DeleteByID(context.Background(), 1, 1)

	assert.EqualError(t, err, api_errors.ErrNotesNotFound.Error())
}
