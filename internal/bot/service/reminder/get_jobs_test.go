package reminder

import (
	"context"
	"database/sql"
	"testing"

	mock_reminder "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/reminder/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAllJobs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mock_reminder.NewMockreminderEditor(ctrl)
	n := New(reminderEditor)

	expectResult := []uuid.UUID{uuid.New(), uuid.New(), uuid.New(), uuid.New()}

	reminderEditor.EXPECT().GetAllJobs(gomock.Any(), gomock.Any()).Return(expectResult, nil)

	jobs, err := n.GetAllJobs(context.Background(), int64(1))
	assert.NoError(t, err)

	assert.Equal(t, expectResult, jobs)
}

func TestGetAllJobs_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mock_reminder.NewMockreminderEditor(ctrl)
	n := New(reminderEditor)

	sqlErr := sql.ErrNoRows

	reminderEditor.EXPECT().GetAllJobs(gomock.Any(), gomock.Any()).Return(nil, sqlErr)

	_, err := n.GetAllJobs(context.Background(), int64(1))
	assert.EqualError(t, sqlErr, err.Error())
}

func TestGetJobID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mock_reminder.NewMockreminderEditor(ctrl)
	n := New(reminderEditor)

	expectResult := uuid.New()

	reminderEditor.EXPECT().GetJobID(gomock.Any(), gomock.Any(), gomock.Any()).Return(expectResult, nil)

	id, err := n.GetJobID(context.Background(), int64(1), 1)
	assert.NoError(t, err)

	assert.Equal(t, expectResult, id)
}

func TestGetJobID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mock_reminder.NewMockreminderEditor(ctrl)
	n := New(reminderEditor)

	sqlErr := sql.ErrNoRows

	reminderEditor.EXPECT().GetJobID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.UUID{}, sqlErr)

	_, err := n.GetJobID(context.Background(), int64(1), 1)
	assert.EqualError(t, sqlErr, err.Error())
}
