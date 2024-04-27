package controller

import (
	"context"
	"testing"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/mocks"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	model_user "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/reminder"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/user"
	tz_cache "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/storage/cache/timezone"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/Pizhlo/bot-reminder-go-telegram/pkg/random"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	tele "gopkg.in/telebot.v3"
)

func TestMinutesDuration_Valid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := &tele.Chat{ID: int64(1)}
	message := &tele.Message{Text: "50"}

	randomReminder := random.Reminder()
	randomReminder.Type = model.SeveralTimesDayType
	randomReminder.Date = "minutes"
	randomReminder.Time = "50"
	reminderSrv.SaveName(chat.ID, randomReminder.Name)
	reminderSrv.SaveDate(chat.ID, randomReminder.Date)
	reminderSrv.SaveType(chat.ID, randomReminder.Type)

	reminderEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Return(nil, api_errors.ErrRemindersNotFound).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	})

	reminderEditor.EXPECT().Save(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, reminder *model.Reminder) {
		assert.Equal(t, randomReminder.TgID, reminder.TgID)
		assert.Equal(t, randomReminder.Date, reminder.Date)
		assert.Equal(t, randomReminder.Time, reminder.Time)
		assert.Equal(t, randomReminder.Type, reminder.Type)
		assert.Equal(t, randomReminder.Name, reminder.Name)
	}).Return(randomReminder.ID, nil)

	reminderEditor.EXPECT().SaveJob(gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx interface{}, reminderID uuid.UUID, jobID uuid.UUID) {
		assert.Equal(t, randomReminder.ID, reminderID)
	})

	telectx.EXPECT().Chat().Return(chat).Times(6)
	telectx.EXPECT().Message().Return(message)

	tzEditor := mocks.NewMocktimezoneEditor(ctrl)
	userEditor := mocks.NewMockuserEditor(ctrl)
	tz := tz_cache.New()

	// при создании user service
	tzEditor.EXPECT().GetAll(gomock.Any()).Return([]*model_user.User{
		{
			TGID: chat.ID,
			Timezone: model_user.Timezone{
				Name: "Europe/Moscow",
			},
		},
	}, nil)

	userSrv := user.New(context.Background(), userEditor, tz, tzEditor)

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Return(nil)

	controller := New(userSrv, nil, nil, reminderSrv)

	err := controller.MinutesDuration(ctx, telectx)
	assert.NoError(t, err)
}

func TestMinutesDuration_Invalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := &tele.Chat{ID: int64(1)}
	message := &tele.Message{Text: random.String(10)}

	randomReminder := random.Reminder()
	randomReminder.Type = model.SeveralTimesDayType
	randomReminder.Date = "minutes"
	reminderSrv.SaveName(chat.ID, randomReminder.Name)
	reminderSrv.SaveDate(chat.ID, randomReminder.Date)
	reminderSrv.SaveType(chat.ID, randomReminder.Type)

	telectx.EXPECT().Chat().Return(chat)
	telectx.EXPECT().Message().Return(message)

	expectedTxt := messages.InvalidMinutesMessage
	expectedKb := view.BackToReminderMenuBtns()
	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, kb *tele.ReplyMarkup) {
		assert.Equal(t, expectedTxt, msg)
		assert.Equal(t, expectedKb, kb)
	})

	controller := New(nil, nil, nil, reminderSrv)

	err := controller.MinutesDuration(ctx, telectx)
	assert.NoError(t, err)
}

func TestHoursDuration_Valid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := &tele.Chat{ID: int64(1)}
	message := &tele.Message{Text: "2"}

	randomReminder := random.Reminder()
	randomReminder.Type = model.SeveralTimesDayType
	randomReminder.Date = "hours"
	randomReminder.Time = "2"
	reminderSrv.SaveName(chat.ID, randomReminder.Name)
	reminderSrv.SaveDate(chat.ID, randomReminder.Date)
	reminderSrv.SaveType(chat.ID, randomReminder.Type)

	reminderEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Return(nil, api_errors.ErrRemindersNotFound).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	})

	reminderEditor.EXPECT().Save(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, reminder *model.Reminder) {
		assert.Equal(t, randomReminder.TgID, reminder.TgID)
		assert.Equal(t, randomReminder.Date, reminder.Date)
		assert.Equal(t, randomReminder.Time, reminder.Time)
		assert.Equal(t, randomReminder.Type, reminder.Type)
		assert.Equal(t, randomReminder.Name, reminder.Name)
	}).Return(randomReminder.ID, nil)

	reminderEditor.EXPECT().SaveJob(gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx interface{}, reminderID uuid.UUID, jobID uuid.UUID) {
		assert.Equal(t, randomReminder.ID, reminderID)
	})

	telectx.EXPECT().Chat().Return(chat).Times(6)
	telectx.EXPECT().Message().Return(message)

	tzEditor := mocks.NewMocktimezoneEditor(ctrl)
	userEditor := mocks.NewMockuserEditor(ctrl)
	tz := tz_cache.New()

	// при создании user service
	tzEditor.EXPECT().GetAll(gomock.Any()).Return([]*model_user.User{
		{
			TGID: chat.ID,
			Timezone: model_user.Timezone{
				Name: "Europe/Moscow",
			},
		},
	}, nil)

	userSrv := user.New(context.Background(), userEditor, tz, tzEditor)

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Return(nil)

	controller := New(userSrv, nil, nil, reminderSrv)

	err := controller.HoursDuration(ctx, telectx)
	assert.NoError(t, err)
}

func TestHoursDuration_Invalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := &tele.Chat{ID: int64(1)}
	message := &tele.Message{Text: random.String(10)}

	randomReminder := random.Reminder()
	randomReminder.Type = model.SeveralTimesDayType
	randomReminder.Date = "hours"
	reminderSrv.SaveName(chat.ID, randomReminder.Name)
	reminderSrv.SaveDate(chat.ID, randomReminder.Date)
	reminderSrv.SaveType(chat.ID, randomReminder.Type)

	telectx.EXPECT().Chat().Return(chat)
	telectx.EXPECT().Message().Return(message)

	expectedTxt := messages.InvalidHoursMessage
	expectedKb := view.BackToReminderMenuBtns()
	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, kb *tele.ReplyMarkup) {
		assert.Equal(t, expectedTxt, msg)
		assert.Equal(t, expectedKb, kb)
	})

	controller := New(nil, nil, nil, reminderSrv)

	err := controller.HoursDuration(ctx, telectx)
	assert.NoError(t, err)
}

func TestTimes_Valid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := &tele.Chat{ID: int64(1)}
	messages := []tele.Message{
		{Text: "19:30 20:30"},
		{Text: "19:30, 20:30"},
	}

	reminderEditor.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Return(nil, api_errors.ErrRemindersNotFound).Do(func(ctx interface{}, userID int64) {
		assert.Equal(t, chat.ID, userID)
	})

	for _, m := range messages {
		randomReminder := random.Reminder()
		randomReminder.Type = model.SeveralTimesDayType
		randomReminder.Date = "times_reminder"
		randomReminder.Time = m.Text
		reminderSrv.SaveName(chat.ID, randomReminder.Name)
		reminderSrv.SaveDate(chat.ID, randomReminder.Date)
		reminderSrv.SaveType(chat.ID, randomReminder.Type)

		reminderEditor.EXPECT().Save(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, reminder *model.Reminder) {
			assert.Equal(t, randomReminder.TgID, reminder.TgID)
			assert.Equal(t, randomReminder.Date, reminder.Date)
			assert.Equal(t, randomReminder.Time, reminder.Time)
			assert.Equal(t, randomReminder.Type, reminder.Type)
			assert.Equal(t, randomReminder.Name, reminder.Name)
		}).Return(randomReminder.ID, nil)

		reminderEditor.EXPECT().SaveJob(gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx interface{}, reminderID uuid.UUID, jobID uuid.UUID) {
			assert.Equal(t, randomReminder.ID, reminderID)
		})

		telectx.EXPECT().Chat().Return(chat).Times(6)
		telectx.EXPECT().Message().Return(&m)

		tzEditor := mocks.NewMocktimezoneEditor(ctrl)
		userEditor := mocks.NewMockuserEditor(ctrl)
		tz := tz_cache.New()

		// при создании user service
		tzEditor.EXPECT().GetAll(gomock.Any()).Return([]*model_user.User{
			{
				TGID: chat.ID,
				Timezone: model_user.Timezone{
					Name: "Europe/Moscow",
				},
			},
		}, nil)

		userSrv := user.New(context.Background(), userEditor, tz, tzEditor)

		telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Return(nil)

		controller := New(userSrv, nil, nil, reminderSrv)

		err := controller.Times(ctx, telectx)
		assert.NoError(t, err)
	}

}

func TestTimes_Invalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := &tele.Chat{ID: int64(1)}
	message := &tele.Message{Text: random.String(10)}

	randomReminder := random.Reminder()
	randomReminder.Type = model.SeveralTimesDayType
	randomReminder.Date = "times_reminder"
	reminderSrv.SaveName(chat.ID, randomReminder.Name)
	reminderSrv.SaveDate(chat.ID, randomReminder.Date)
	reminderSrv.SaveType(chat.ID, randomReminder.Type)

	telectx.EXPECT().Chat().Return(chat)
	telectx.EXPECT().Message().Return(message)

	expectedTxt := messages.InvalidTimesMessage
	expectedKb := view.BackToReminderMenuBtns()
	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, kb *tele.ReplyMarkup) {
		assert.Equal(t, expectedTxt, msg)
		assert.Equal(t, expectedKb, kb)
	})

	controller := New(nil, nil, nil, reminderSrv)

	err := controller.Times(ctx, telectx)
	assert.NoError(t, err)
}

func TestDaysDuration_Valid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := &tele.Chat{ID: int64(1)}
	message := &tele.Message{Text: "2"}

	randomReminder := random.Reminder()
	reminderSrv.SaveName(chat.ID, randomReminder.Name)

	telectx.EXPECT().Chat().Return(chat).Times(2)
	telectx.EXPECT().Message().Return(message).Times(2)

	expectedTxt := messages.ReminderTimeMessage
	expectedKb := view.BackToReminderMenuBtns()
	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, kb *tele.ReplyMarkup) {
		assert.Equal(t, expectedTxt, msg)
		assert.Equal(t, expectedKb, kb)
	})

	controller := New(nil, nil, nil, reminderSrv)

	err := controller.DaysDuration(ctx, telectx)
	assert.NoError(t, err)
}

func TestDaysDuration_NotInt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := &tele.Chat{ID: int64(1)}
	message := &tele.Message{Text: random.String(10)}

	randomReminder := random.Reminder()
	reminderSrv.SaveName(chat.ID, randomReminder.Name)

	telectx.EXPECT().Chat().Return(chat)
	telectx.EXPECT().Message().Return(message)

	controller := New(nil, nil, nil, reminderSrv)

	err := controller.DaysDuration(ctx, telectx)
	assert.EqualError(t, err, api_errors.ErrInvalidDays.Error())
}

func TestDaysDuration_OutOfRange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := &tele.Chat{ID: int64(1)}

	mesages := []tele.Message{
		{Text: "-1"},
		{Text: "0"},
		{Text: "181"},
		{Text: "200"},
	}

	randomReminder := random.Reminder()
	reminderSrv.SaveName(chat.ID, randomReminder.Name)

	controller := New(nil, nil, nil, reminderSrv)

	for _, m := range mesages {
		telectx.EXPECT().Chat().Return(chat)
		telectx.EXPECT().Message().Return(&m)
		err := controller.DaysDuration(ctx, telectx)
		assert.EqualError(t, err, api_errors.ErrInvalidDays.Error())
	}

}

func TestDaysInMonthDurationn_Valid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := &tele.Chat{ID: int64(1)}
	message := &tele.Message{Text: "2"}

	randomReminder := random.Reminder()
	reminderSrv.SaveName(chat.ID, randomReminder.Name)

	telectx.EXPECT().Chat().Return(chat).Times(2)
	telectx.EXPECT().Message().Return(message).Times(2)

	expectedTxt := messages.ReminderTimeMessage
	expectedKb := view.BackToReminderMenuBtns()
	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, kb *tele.ReplyMarkup) {
		assert.Equal(t, expectedTxt, msg)
		assert.Equal(t, expectedKb, kb)
	})

	controller := New(nil, nil, nil, reminderSrv)

	err := controller.DaysInMonthDuration(ctx, telectx)
	assert.NoError(t, err)
}

func TestDaysInMonthDuration_NotInt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := &tele.Chat{ID: int64(1)}
	message := &tele.Message{Text: random.String(10)}

	randomReminder := random.Reminder()
	reminderSrv.SaveName(chat.ID, randomReminder.Name)

	telectx.EXPECT().Chat().Return(chat)
	telectx.EXPECT().Message().Return(message)

	controller := New(nil, nil, nil, reminderSrv)

	err := controller.DaysInMonthDuration(ctx, telectx)
	assert.EqualError(t, err, api_errors.ErrInvalidDays.Error())
}

func TestDaysInMonthDuration_OutOfRange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := &tele.Chat{ID: int64(1)}

	mesages := []tele.Message{
		{Text: "0"},
		{Text: "32"},
		{Text: "500"},
	}

	randomReminder := random.Reminder()
	reminderSrv.SaveName(chat.ID, randomReminder.Name)

	controller := New(nil, nil, nil, reminderSrv)

	for _, m := range mesages {
		telectx.EXPECT().Chat().Return(chat)
		telectx.EXPECT().Message().Return(&m)
		err := controller.DaysInMonthDuration(ctx, telectx)
		assert.EqualError(t, err, api_errors.ErrInvalidDays.Error())
	}

}
