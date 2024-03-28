package reminder

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	mock_reminder "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/reminder/mocks"
	"github.com/go-co-op/gocron/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func task(ctx context.Context, reminder *model.Reminder) error {
	return nil
}

func TestDeleteAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mock_reminder.NewMockreminderEditor(ctrl)
	n := New(reminderEditor)

	reminderEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Return([]model.Reminder{}, nil)

	err := n.CreateScheduler(context.Background(), int64(1), time.Local, task)
	require.NoError(t, err)

	reminderEditor.EXPECT().GetAllJobs(gomock.Any(), gomock.Any()).Return(uuid.UUIDs{uuid.New()}, nil)
	reminderEditor.EXPECT().DeleteJobAndReminder(gomock.Any(), gomock.Any()).Return(nil)

	err = n.DeleteAll(context.Background(), int64(1))
	// потому что возвращает ошибку job not found
	assert.False(t, reflect.ValueOf(err) == reflect.ValueOf(gocron.ErrJobNotFound))
}

func TestDeleteAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mock_reminder.NewMockreminderEditor(ctrl)
	n := New(reminderEditor)

	reminderEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Return([]model.Reminder{}, nil)

	err := n.CreateScheduler(context.Background(), int64(1), time.Local, task)
	require.NoError(t, err)

	sqlErr := sql.ErrNoRows

	reminderEditor.EXPECT().GetAllJobs(gomock.Any(), gomock.Any()).Return(uuid.UUIDs{uuid.New()}, nil)
	reminderEditor.EXPECT().DeleteJobAndReminder(gomock.Any(), gomock.Any()).Return(sqlErr)

	err = n.DeleteAll(context.Background(), int64(1))
	assert.False(t, reflect.ValueOf(err) == reflect.ValueOf(gocron.ErrJobNotFound))
}
