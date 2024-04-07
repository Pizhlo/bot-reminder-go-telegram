package reminder

import (
	"context"
	"errors"
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

	reminders, err := n.GetAll(context.Background(), userID)
	assert.NoError(t, err)

	msg, err := n.Message(userID, reminders)
	assert.NoError(t, err)
	assert.Equal(t, view.First(), msg)

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
