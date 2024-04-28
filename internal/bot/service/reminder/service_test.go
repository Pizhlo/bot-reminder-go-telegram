package reminder

import (
	"context"
	"testing"

	mock_reminder "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/mocks"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/pkg/random"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLoadMemory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminders := random.Reminders(5)

	for i := 0; i < len(reminders); i++ {
		reminders[i].TgID = int64(random.Int(2, 1000))
	}

	reminderEditor := mock_reminder.NewMockreminderEditor(ctrl)
	reminderEditor.EXPECT().GetMemory(gomock.Any()).Return(reminders, nil)

	n := New(reminderEditor)
	err := n.LoadMemory(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, len(reminders), len(n.reminderMap))

	for _, r := range reminders {
		val, ok := n.reminderMap[r.TgID]
		assert.True(t, ok)

		assert.Equal(t, r.Name, val.Name)
		assert.Equal(t, r.Date, val.Date)
		assert.Equal(t, r.Time, val.Time)
		assert.Equal(t, r.Type, val.Type)

	}
}

func TestLoadMemory_WithoutType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminders := random.Reminders(5)

	for i := 0; i < len(reminders); i++ {
		reminders[i].TgID = int64(random.Int(2, 1000))
		reminders[i].Type = ""
	}

	reminderEditor := mock_reminder.NewMockreminderEditor(ctrl)
	reminderEditor.EXPECT().GetMemory(gomock.Any()).Return(reminders, nil)

	n := New(reminderEditor)
	err := n.LoadMemory(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, len(reminders), len(n.reminderMap))

	for _, r := range reminders {
		val, ok := n.reminderMap[r.TgID]
		assert.True(t, ok)

		assert.Equal(t, r.Name, val.Name)
		assert.Equal(t, r.Date, val.Date)
		assert.Equal(t, r.Time, val.Time)
		assert.Equal(t, r.Type, val.Type)

	}
}

func TestLoadMemory_WithoutDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminders := random.Reminders(5)

	for i := 0; i < len(reminders); i++ {
		reminders[i].TgID = int64(random.Int(2, 1000))
		reminders[i].Date = ""
	}

	reminderEditor := mock_reminder.NewMockreminderEditor(ctrl)
	reminderEditor.EXPECT().GetMemory(gomock.Any()).Return(reminders, nil)

	n := New(reminderEditor)
	err := n.LoadMemory(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, len(reminders), len(n.reminderMap))

	for _, r := range reminders {
		val, ok := n.reminderMap[r.TgID]
		assert.True(t, ok)

		assert.Equal(t, r.Name, val.Name)
		assert.Equal(t, r.Date, val.Date)
		assert.Equal(t, r.Time, val.Time)
		assert.Equal(t, r.Type, val.Type)

	}
}

func TestLoadMemory_WithoutTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminders := random.Reminders(5)

	for i := 0; i < len(reminders); i++ {
		reminders[i].TgID = int64(random.Int(2, 1000))
		reminders[i].Time = ""
	}

	reminderEditor := mock_reminder.NewMockreminderEditor(ctrl)
	reminderEditor.EXPECT().GetMemory(gomock.Any()).Return(reminders, nil)

	n := New(reminderEditor)
	err := n.LoadMemory(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, len(reminders), len(n.reminderMap))

	for _, r := range reminders {
		val, ok := n.reminderMap[r.TgID]
		assert.True(t, ok)

		assert.Equal(t, r.Name, val.Name)
		assert.Equal(t, r.Date, val.Date)
		assert.Equal(t, r.Time, val.Time)
		assert.Equal(t, r.Type, val.Type)

	}
}

func TestLoadMemory_OnlyTwoFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminders := random.Reminders(5)

	for i := 0; i < len(reminders); i++ {
		reminders[i].TgID = int64(random.Int(2, 1000))
		reminders[i].Time = ""
		reminders[i].Date = ""
		reminders[i].Type = ""
	}

	reminderEditor := mock_reminder.NewMockreminderEditor(ctrl)
	reminderEditor.EXPECT().GetMemory(gomock.Any()).Return(reminders, nil)

	n := New(reminderEditor)
	err := n.LoadMemory(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, len(reminders), len(n.reminderMap))

	for _, r := range reminders {
		val, ok := n.reminderMap[r.TgID]
		assert.True(t, ok)

		assert.Equal(t, r.Name, val.Name)
		assert.Equal(t, r.Date, val.Date)
		assert.Equal(t, r.Time, val.Time)
		assert.Equal(t, r.Type, val.Type)

	}
}

func TestSaveMemory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminders := random.Reminders(5)

	for i := 0; i < len(reminders); i++ {
		reminders[i].TgID = int64(random.Int(2, 1000))
	}

	reminderEditor := mock_reminder.NewMockreminderEditor(ctrl)

	n := New(reminderEditor)

	for _, r := range reminders {
		n.reminderMap[r.TgID] = &r
	}

	for _, r := range n.reminderMap {
		reminderEditor.EXPECT().SaveMemory(gomock.Any(), gomock.Any()).Return(nil).Do(func(ctx interface{}, reminder *model.Reminder) {
			assert.Equal(t, r, reminder)
		})
	}

	err := n.SaveMemory(context.Background())
	assert.NoError(t, err)
}
