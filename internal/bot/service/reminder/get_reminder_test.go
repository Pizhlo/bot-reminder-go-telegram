package reminder

import (
	"context"
	"database/sql"
	"testing"

	mock_reminder "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/mocks"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/Pizhlo/bot-reminder-go-telegram/pkg/random"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userID := int64(1)

	reminderEditor := mock_reminder.NewMockreminderEditor(ctrl)
	n := New(reminderEditor)
	n.SaveUser(userID)

	// ожидаемые напоминания, которые возвращает база
	expectedResult := random.Reminders(5)

	// создаем view для генерации сообщения на основе данных из БД
	view := view.NewReminder()
	// подготавливаем сообщение - сохраняем во view
	view.Message(expectedResult)

	reminderEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Return(expectedResult, nil)

	msg, _, err := n.GetAll(context.Background(), userID)
	assert.NoError(t, err)

	assert.Equal(t, view.First(), msg)
}

func TestGetAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userID := int64(1)

	reminderEditor := mock_reminder.NewMockreminderEditor(ctrl)
	n := New(reminderEditor)
	n.SaveUser(userID)

	sqlErr := sql.ErrNoRows

	reminderEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Return(nil, sqlErr)

	msg, _, err := n.GetAll(context.Background(), userID)
	assert.EqualError(t, err, sqlErr.Error())

	assert.Equal(t, "", msg)
}
