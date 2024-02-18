package reminder

import (
	"context"
	"database/sql"
	"testing"

	mock_reminder "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/reminder/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDeleteAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mock_reminder.NewMockreminderEditor(ctrl)
	n := New(reminderEditor)

	reminderEditor.EXPECT().DeleteAllByUserID(gomock.Any(), gomock.Any()).Return(nil)

	err := n.DeleteAll(context.Background(), int64(1))
	assert.NoError(t, err)
}

func TestDeleteAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mock_reminder.NewMockreminderEditor(ctrl)
	n := New(reminderEditor)

	sqlErr := sql.ErrNoRows

	reminderEditor.EXPECT().DeleteAllByUserID(gomock.Any(), gomock.Any()).Return(sqlErr)

	err := n.DeleteAll(context.Background(), int64(1))
	assert.EqualError(t, sqlErr, err.Error())
}

func TestDeleteReminderByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mock_reminder.NewMockreminderEditor(ctrl)
	n := New(reminderEditor)

	reminderEditor.EXPECT().DeleteReminderByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	err := n.DeleteReminderByID(context.Background(), int64(1), 1)
	assert.NoError(t, err)
}

func TestDeleteReminderByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mock_reminder.NewMockreminderEditor(ctrl)
	n := New(reminderEditor)

	sqlErr := sql.ErrNoRows

	reminderEditor.EXPECT().DeleteReminderByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(sqlErr)

	err := n.DeleteReminderByID(context.Background(), int64(1), 1)
	assert.EqualError(t, sqlErr, err.Error())
}
