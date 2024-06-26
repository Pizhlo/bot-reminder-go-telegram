package reminder

import (
	"context"
	"errors"
	"testing"
	"time"

	mock_reminder "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/mocks"
	"github.com/Pizhlo/bot-reminder-go-telegram/pkg/random"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSave(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userID := int64(1)
	reminder := random.Reminder()

	reminderEditor := mock_reminder.NewMockreminderEditor(ctrl)
	n := New(reminderEditor)

	// создаем напоминание в памяти
	n.SaveName(userID, reminder.Name)

	err := n.SaveDate(userID, reminder.Date)
	require.NoError(t, err)

	err = n.SaveType(userID, reminder.Type)
	require.NoError(t, err)

	err = n.SaveCreatedField(userID, reminder.Created.Location())
	require.NoError(t, err)

	err = n.ParseTime(userID, reminder.Time)
	require.NoError(t, err)

	err = n.SaveTime(userID, reminder.Time)
	require.NoError(t, err)

	reminderEditor.EXPECT().Save(gomock.Any(), gomock.Any()).Return(reminder.ID, nil)

	// сохраняем напоминание
	id, err := n.Save(context.Background(), userID, reminder)
	require.NoError(t, err)

	// проверяем, что сохраненное напоминание совпадает с переданным
	result, ok := n.reminderMap[userID]
	assert.Equal(t, true, ok)

	assert.Equal(t, id, result.ID)
	assert.Equal(t, reminder.TgID, result.TgID)
	assert.Equal(t, reminder.Name, result.Name)
	assert.Equal(t, reminder.Created.Location(), result.Created.Location())
	assert.Equal(t, reminder.Date, result.Date)
	assert.Equal(t, reminder.Type, result.Type)
}

func TestSave_DBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userID := int64(1)
	reminder := random.Reminder()

	reminderEditor := mock_reminder.NewMockreminderEditor(ctrl)
	n := New(reminderEditor)

	// создаем напоминание в памяти
	n.SaveName(userID, reminder.Name)

	err := n.SaveType(userID, reminder.Type)
	assert.NoError(t, err)

	err = n.SaveDate(userID, reminder.Date)
	assert.NoError(t, err)

	err = n.SaveTime(userID, reminder.Time)
	assert.NoError(t, err)

	err = n.SaveType(userID, reminder.Type)
	assert.NoError(t, err)

	err = n.SaveCreatedField(userID, time.Local)
	assert.NoError(t, err)

	testErr := errors.New("test error")

	reminderEditor.EXPECT().Save(gomock.Any(), gomock.Any()).Return(uuid.New(), testErr)

	_, err = n.Save(context.Background(), userID, reminder)
	assert.EqualError(t, err, testErr.Error())
}

func TestSaveJobID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userID := int64(1)
	reminder := random.Reminder()

	reminderEditor := mock_reminder.NewMockreminderEditor(ctrl)
	n := New(reminderEditor)

	// создаем напоминание в памяти
	n.SaveName(userID, reminder.Name)

	err := n.SaveDate(userID, reminder.Date)
	require.NoError(t, err)

	err = n.SaveType(userID, reminder.Type)
	require.NoError(t, err)

	err = n.SaveCreatedField(userID, reminder.Created.Location())
	require.NoError(t, err)

	err = n.ParseTime(userID, reminder.Time)
	require.NoError(t, err)

	reminderEditor.EXPECT().SaveJob(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// сохраняем job ID
	err = n.SaveJobID(context.Background(), uuid.New(), uuid.New())
	require.NoError(t, err)
}

func TestSaveJobID_DBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userID := int64(1)
	reminder := random.Reminder()

	reminderEditor := mock_reminder.NewMockreminderEditor(ctrl)
	n := New(reminderEditor)

	// создаем напоминание в памяти
	n.SaveName(userID, reminder.Name)

	err := n.SaveDate(userID, reminder.Date)
	require.NoError(t, err)

	err = n.SaveType(userID, reminder.Type)
	require.NoError(t, err)

	err = n.SaveCreatedField(userID, reminder.Created.Location())
	require.NoError(t, err)

	err = n.ParseTime(userID, reminder.Time)
	require.NoError(t, err)

	testErr := errors.New("test error")

	reminderEditor.EXPECT().SaveJob(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)

	// сохраняем job ID
	err = n.SaveJobID(context.Background(), uuid.New(), uuid.New())
	require.EqualError(t, err, testErr.Error())
}
