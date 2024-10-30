package reminder

import (
	"context"
	"errors"
	"testing"
	"time"

	mock_reminder "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/mocks"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/Pizhlo/bot-reminder-go-telegram/pkg/random"
	"github.com/go-co-op/gocron/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userID := int64(1)

	reminderEditor := mock_reminder.NewMockreminderEditor(ctrl)
	n := New(reminderEditor)
	n.SaveUser(userID)

	// ожидаемые напоминания, которые возвращает база
	expectedResult := random.Reminders(5)

	// один раз - при создании планировщика, еще раз - при вызове n.GetAll()
	reminderEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Return(expectedResult, nil).Times(2)

	reminderEditor.EXPECT().SaveJob(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(len(expectedResult)).Do(func(ctx interface{}, reminderID uuid.UUID, jobID uuid.UUID) {
		for i := 0; i < len(expectedResult); i++ {
			if expectedResult[i].ID == reminderID {
				expectedResult[i].Job.ID = jobID
				return
			}
		}
	})

	err := n.CreateScheduler(ctx, userID, time.Local, func(ctx context.Context, reminder *model.Reminder) error { return nil })
	assert.NoError(t, err)

	sch, err := n.GetScheduler(userID)
	require.NoError(t, err)

	jobs := sch.Jobs()

	jobsMap := map[uuid.UUID]gocron.Job{}

	for _, j := range jobs {
		jobsMap[j.ID()] = j
	}

	// заполняем поле nextRun у всех сгенерированных напоминаний
	for i := 0; i < len(expectedResult); i++ {
		j, ok := jobsMap[expectedResult[i].Job.ID]
		if !ok {
			t.Fatalf("job not found in jobs map")
			continue
		}

		nextRun, err := j.NextRun()
		if err != nil {
			t.Fatalf("error getting next run for job: %+v", err)
			continue
		}

		expectedResult[i].Job.NextRun = nextRun
	}

	// создаем view для генерации сообщения на основе данных из БД
	view := view.NewReminder()
	// подготавливаем сообщение - сохраняем во view
	view.Message(expectedResult)

	reminders, err := n.GetAll(context.Background(), userID)
	assert.NoError(t, err)

	actual, err := n.Message(userID, reminders)
	assert.NoError(t, err)

	expected := view.First()
	assert.Equal(t, expected, actual)

	kb := n.Keyboard(userID)
	assert.Equal(t, view.Keyboard(), kb)
}

func TestGetAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userID := int64(1)

	reminderEditor := mock_reminder.NewMockreminderEditor(ctrl)
	n := New(reminderEditor)
	n.SaveUser(userID)

	testErr := errors.New("test err")

	reminderEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Return(nil, testErr)

	_, err := n.GetAll(context.Background(), userID)
	assert.EqualError(t, err, testErr.Error())

}
