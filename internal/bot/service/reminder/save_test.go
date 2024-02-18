package reminder

import (
	"context"
	"database/sql"
	"testing"

	mock_reminder "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/reminder/mocks"
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

	err = n.ProcessTime(userID, reminder.Time)
	require.NoError(t, err)

	reminderEditor.EXPECT().Save(gomock.Any(), gomock.Any()).Return(reminder.ID, nil)

	// сохраняем напоминание
	err = n.Save(context.Background(), userID)
	require.NoError(t, err)

	// проверяем, что сохраненное напоминание совпадает с переданным
	result, ok := n.reminderMap[userID]
	assert.Equal(t, true, ok)

	assert.Equal(t, reminder.ID, result.ID)
	assert.Equal(t, reminder.TgID, result.TgID)
	assert.Equal(t, reminder.Name, result.Name)
	assert.Equal(t, reminder.Created.Location(), result.Created.Location())
	assert.Equal(t, reminder.Date, result.Date)
	assert.Equal(t, reminder.Type, result.Type)
}

func TestSave_NotFoundInMap(t *testing.T) {
	userID := int64(1)
	n := New(nil)

	err := n.Save(context.Background(), userID)
	assert.EqualError(t, err, "error while getting reminder by user ID: reminder not found")
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

	sqlErr := sql.ErrNoRows

	reminderEditor.EXPECT().Save(gomock.Any(), gomock.Any()).Return(int64(0), sqlErr)

	err := n.Save(context.Background(), userID)
	assert.EqualError(t, err, sqlErr.Error())
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

	err = n.ProcessTime(userID, reminder.Time)
	require.NoError(t, err)

	reminderEditor.EXPECT().SaveJob(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// сохраняем job ID
	err = n.SaveJobID(context.Background(), uuid.New(), userID)
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

	err = n.ProcessTime(userID, reminder.Time)
	require.NoError(t, err)

	sqlErr := sql.ErrNoRows

	reminderEditor.EXPECT().SaveJob(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(sqlErr)

	// сохраняем job ID
	err = n.SaveJobID(context.Background(), uuid.New(), userID)
	require.EqualError(t, err, sqlErr.Error())
}

func TestSaveJobID_NotFoundInMap(t *testing.T) {
	userID := int64(1)

	n := New(nil)

	// сохраняем job ID
	err := n.SaveJobID(context.Background(), uuid.New(), userID)
	require.EqualError(t, err, "error while getting reminder by user ID: reminder not found")
}
