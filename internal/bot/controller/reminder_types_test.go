package controller

import (
	"context"
	"fmt"
	"testing"
	"time"

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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tele "gopkg.in/telebot.v3"
)

func TestToday(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	tzEditor := mocks.NewMocktimezoneEditor(ctrl)
	userEditor := mocks.NewMockuserEditor(ctrl)
	tz := tz_cache.New()
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := &tele.Chat{ID: int64(1)}

	telectx.EXPECT().Chat().Return(chat).Times(3)

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

	controller := New(userSrv, nil, nil, reminderSrv)

	randomReminder := random.Reminder()
	reminderSrv.SaveName(chat.ID, randomReminder.Name)

	expectedTxt := messages.ReminderTimeMessage
	expectedOpts := &tele.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.BackToReminderMenuBtns(),
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(msg interface{}, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedOpts, sendOpts)
		assert.Equal(t, expectedTxt, msg)
	}).Return(nil)

	err := controller.Today(context.Background(), telectx)
	assert.NoError(t, err)

	r, err := reminderSrv.GetFromMemory(chat.ID)
	assert.NoError(t, err)

	assert.Equal(t, model.DateType, r.Type)

	layout := "02.01.2006"

	loc, err := time.LoadLocation("Europe/Moscow")
	require.NoError(t, err)

	date := time.Now().In(loc).Format(layout)

	assert.Equal(t, date, r.Date)
}

func TestTomorrow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	tzEditor := mocks.NewMocktimezoneEditor(ctrl)
	userEditor := mocks.NewMockuserEditor(ctrl)
	tz := tz_cache.New()
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := &tele.Chat{ID: int64(1)}

	telectx.EXPECT().Chat().Return(chat).Times(3)

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

	controller := New(userSrv, nil, nil, reminderSrv)

	randomReminder := random.Reminder()
	reminderSrv.SaveName(chat.ID, randomReminder.Name)

	expectedTxt := messages.ReminderTimeMessage
	expectedOpts := &tele.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.BackToReminderMenuBtns(),
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(msg interface{}, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedOpts, sendOpts)
		assert.Equal(t, expectedTxt, msg)
	}).Return(nil)

	err := controller.Tomorrow(context.Background(), telectx)
	assert.NoError(t, err)

	r, err := reminderSrv.GetFromMemory(chat.ID)
	assert.NoError(t, err)

	assert.Equal(t, model.DateType, r.Type)

	layout := "02.01.2006"

	loc, err := time.LoadLocation("Europe/Moscow")
	require.NoError(t, err)

	date := time.Now().In(loc).Add(24 * time.Hour).Format(layout)

	assert.Equal(t, date, r.Date)
}

func TestEverydayReminder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	tzEditor := mocks.NewMocktimezoneEditor(ctrl)
	userEditor := mocks.NewMockuserEditor(ctrl)
	tz := tz_cache.New()
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := &tele.Chat{ID: int64(1)}

	telectx.EXPECT().Chat().Return(chat).Times(2)

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

	controller := New(userSrv, nil, nil, reminderSrv)

	randomReminder := random.Reminder()
	reminderSrv.SaveName(chat.ID, randomReminder.Name)

	expectedTxt := messages.ReminderTimeMessage
	expectedOpts := &tele.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.BackToReminderMenuBtns(),
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(msg interface{}, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedOpts, sendOpts)
		assert.Equal(t, expectedTxt, msg)
	}).Return(nil)

	err := controller.EverydayReminder(context.Background(), telectx)
	assert.NoError(t, err)

	r, err := reminderSrv.GetFromMemory(chat.ID)
	assert.NoError(t, err)

	assert.Equal(t, model.EverydayType, r.Type)
	assert.Equal(t, string(model.EverydayType), r.Date)
}

func TestSeveralTimesADayReminder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := &tele.Chat{ID: int64(1)}

	telectx.EXPECT().Chat().Return(chat).Times(2)
	telectx.EXPECT().Message().Return(&tele.Message{Sender: &tele.User{IsBot: true}})

	controller := New(nil, nil, nil, reminderSrv)

	randomReminder := random.Reminder()
	reminderSrv.SaveName(chat.ID, randomReminder.Name)

	expectedTxt := fmt.Sprintf(messages.ChooseMinutesOrHoursMessage, randomReminder.Name, "несколько раз в день")
	expectedOpts := &tele.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.SeveralTimesBtns(),
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(msg interface{}, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedOpts, sendOpts)
		assert.Equal(t, expectedTxt, msg)
	}).Return(nil)

	err := controller.SeveralTimesADayReminder(context.Background(), telectx)
	assert.NoError(t, err)

	r, err := reminderSrv.GetFromMemory(chat.ID)
	assert.NoError(t, err)

	assert.Equal(t, model.SeveralTimesDayType, r.Type)
}

func TestSeveralTimesADayReminder_NewName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := &tele.Chat{ID: int64(1)}

	newName := random.String(5)

	telectx.EXPECT().Chat().Return(chat).Times(3)
	telectx.EXPECT().Message().Return(&tele.Message{Sender: &tele.User{IsBot: false}, Text: newName}).Times(2)

	controller := New(nil, nil, nil, reminderSrv)

	randomReminder := random.Reminder()
	reminderSrv.SaveName(chat.ID, randomReminder.Name)

	expectedTxt := fmt.Sprintf(messages.ChooseMinutesOrHoursMessage, newName, "несколько раз в день")
	expectedOpts := &tele.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.SeveralTimesBtns(),
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(msg interface{}, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedOpts, sendOpts)
		assert.Equal(t, expectedTxt, msg)
	}).Return(nil)

	err := controller.SeveralTimesADayReminder(context.Background(), telectx)
	assert.NoError(t, err)

	r, err := reminderSrv.GetFromMemory(chat.ID)
	assert.NoError(t, err)

	assert.Equal(t, model.SeveralTimesDayType, r.Type)
	assert.Equal(t, newName, r.Name)
}

func TestOnceInMinutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := &tele.Chat{ID: int64(1)}

	telectx.EXPECT().Chat().Return(chat)

	controller := New(nil, nil, nil, reminderSrv)

	randomReminder := random.Reminder()
	reminderSrv.SaveName(chat.ID, randomReminder.Name)

	expectedTxt := messages.MinutesDurationMessage
	expectedKb := view.BackToReminderMenuBtns()

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(msg interface{}, kb *tele.ReplyMarkup) {
		assert.Equal(t, expectedKb, kb)
		assert.Equal(t, expectedTxt, msg)
	}).Return(nil)

	err := controller.OnceInMinutes(context.Background(), telectx)
	assert.NoError(t, err)

	r, err := reminderSrv.GetFromMemory(chat.ID)
	assert.NoError(t, err)

	assert.Equal(t, minutesDate, r.Date)
}

func TestOnceInHours(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := &tele.Chat{ID: int64(1)}

	telectx.EXPECT().Chat().Return(chat)

	controller := New(nil, nil, nil, reminderSrv)

	randomReminder := random.Reminder()
	reminderSrv.SaveName(chat.ID, randomReminder.Name)

	expectedTxt := messages.HoursDurationMessage
	expectedKb := view.BackToReminderMenuBtns()

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(msg interface{}, kb *tele.ReplyMarkup) {
		assert.Equal(t, expectedKb, kb)
		assert.Equal(t, expectedTxt, msg)
	}).Return(nil)

	err := controller.OnceInHours(context.Background(), telectx)
	assert.NoError(t, err)

	r, err := reminderSrv.GetFromMemory(chat.ID)
	assert.NoError(t, err)

	assert.Equal(t, hoursDate, r.Date)
}

func TestEveryWeek(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := &tele.Chat{ID: int64(1)}

	telectx.EXPECT().Chat().Return(chat).Times(2)
	telectx.EXPECT().Message().Return(&tele.Message{Sender: &tele.User{IsBot: true}})

	controller := New(nil, nil, nil, reminderSrv)

	randomReminder := random.Reminder()
	reminderSrv.SaveName(chat.ID, randomReminder.Name)

	expectedTxt := fmt.Sprintf(messages.ChooseWeekDayMessage, randomReminder.Name, "раз в неделю")
	expectedOpts := &tele.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.WeekMenu(),
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(msg interface{}, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedOpts, sendOpts)
		assert.Equal(t, expectedTxt, msg)
	}).Return(nil)

	err := controller.EveryWeek(context.Background(), telectx)
	assert.NoError(t, err)

	r, err := reminderSrv.GetFromMemory(chat.ID)
	assert.NoError(t, err)

	assert.Equal(t, model.EveryWeekType, r.Type)
}

func TestEveryWeek_NewName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := &tele.Chat{ID: int64(1)}

	newName := random.String(5)

	telectx.EXPECT().Chat().Return(chat).Times(3)
	telectx.EXPECT().Message().Return(&tele.Message{Sender: &tele.User{IsBot: false}, Text: newName}).Times(2)

	controller := New(nil, nil, nil, reminderSrv)

	randomReminder := random.Reminder()
	reminderSrv.SaveName(chat.ID, randomReminder.Name)

	expectedTxt := fmt.Sprintf(messages.ChooseWeekDayMessage, newName, "раз в неделю")
	expectedOpts := &tele.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.WeekMenu(),
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(msg interface{}, sendOpts *tele.SendOptions) {
		assert.Equal(t, expectedOpts, sendOpts)
		assert.Equal(t, expectedTxt, msg)
	}).Return(nil)

	err := controller.EveryWeek(context.Background(), telectx)
	assert.NoError(t, err)

	r, err := reminderSrv.GetFromMemory(chat.ID)
	assert.NoError(t, err)

	assert.Equal(t, model.EveryWeekType, r.Type)
	assert.Equal(t, newName, r.Name)
}

func TestWeekDay(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	controller := New(nil, nil, nil, reminderSrv)

	chat := &tele.Chat{ID: int64(1)}

	weekdays := []string{
		"monday",
		"tuesday",
		"wednesday",
		"thursday",
		"friday",
		"saturday",
		"sunday",
	}

	for _, wd := range weekdays {
		telectx.EXPECT().Chat().Return(chat)
		telectx.EXPECT().Message().Return(&tele.Message{Sender: &tele.User{IsBot: true}})

		randomReminder := random.Reminder()
		reminderSrv.SaveName(chat.ID, randomReminder.Name)

		expectedTxt := messages.ReminderTimeMessage
		expectedKb := view.BackToReminderMenuBtns()

		telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(msg interface{}, kb *tele.ReplyMarkup) {
			assert.Equal(t, expectedKb, kb)
			assert.Equal(t, expectedTxt, msg)
		}).Return(nil)

		telectx.EXPECT().Callback().Return(&tele.Callback{Unique: wd})

		err := controller.WeekDay(context.Background(), telectx)
		assert.NoError(t, err)

		r, err := reminderSrv.GetFromMemory(chat.ID)
		assert.NoError(t, err)

		assert.Equal(t, wd, r.Date)
	}
}

func TestWeekDay_NewName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	controller := New(nil, nil, nil, reminderSrv)

	chat := &tele.Chat{ID: int64(1)}

	newName := random.String(5)

	telectx.EXPECT().Chat().Return(chat)
	telectx.EXPECT().Message().Return(&tele.Message{Sender: &tele.User{IsBot: false}, Text: newName}).Times(3)

	randomReminder := random.Reminder()
	reminderSrv.SaveName(chat.ID, randomReminder.Name)

	expectedTxt := fmt.Sprintf(messages.ChooseWeekDayMessage, newName, "раз в неделю")
	expectedOpts := &tele.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.WeekMenu(),
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(msg interface{}, opts *tele.SendOptions) {
		assert.Equal(t, expectedOpts, opts)
		assert.Equal(t, expectedTxt, msg)
	}).Return(nil)

	err := controller.WeekDay(context.Background(), telectx)
	assert.NoError(t, err)

	r, err := reminderSrv.GetFromMemory(chat.ID)
	assert.NoError(t, err)

	assert.Equal(t, newName, r.Name)
}

func TestSeveralDays(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	controller := New(nil, nil, nil, reminderSrv)

	chat := &tele.Chat{ID: int64(1)}

	telectx.EXPECT().Chat().Return(chat).Times(2)
	telectx.EXPECT().Message().Return(&tele.Message{Sender: &tele.User{IsBot: true}})

	randomReminder := random.Reminder()
	reminderSrv.SaveName(chat.ID, randomReminder.Name)

	expectedTxt := fmt.Sprintf(messages.DaysDurationMessage, randomReminder.Name, "раз в несколько дней")
	expectedOpts := &tele.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.BackToReminderMenuBtns(),
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(msg interface{}, opts *tele.SendOptions) {
		assert.Equal(t, expectedOpts, opts)
		assert.Equal(t, expectedTxt, msg)
	}).Return(nil)

	err := controller.SeveralDays(context.Background(), telectx)
	assert.NoError(t, err)

	r, err := reminderSrv.GetFromMemory(chat.ID)
	assert.NoError(t, err)

	assert.Equal(t, model.SeveralDaysType, r.Type)
}

func TestSeveralDays_NewName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	controller := New(nil, nil, nil, reminderSrv)

	chat := &tele.Chat{ID: int64(1)}

	newName := random.String(5)

	telectx.EXPECT().Chat().Return(chat).Times(3)
	telectx.EXPECT().Message().Return(&tele.Message{Sender: &tele.User{IsBot: false}, Text: newName}).Times(2)

	randomReminder := random.Reminder()
	reminderSrv.SaveName(chat.ID, randomReminder.Name)

	expectedTxt := fmt.Sprintf(messages.DaysDurationMessage, newName, "раз в несколько дней")
	expectedOpts := &tele.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.BackToReminderMenuBtns(),
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(msg interface{}, opts *tele.SendOptions) {
		assert.Equal(t, expectedOpts, opts)
		assert.Equal(t, expectedTxt, msg)
	}).Return(nil)

	err := controller.SeveralDays(context.Background(), telectx)
	assert.NoError(t, err)

	r, err := reminderSrv.GetFromMemory(chat.ID)
	assert.NoError(t, err)

	assert.Equal(t, model.SeveralDaysType, r.Type)
	assert.Equal(t, newName, r.Name)
}

func TestMonth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	controller := New(nil, nil, nil, reminderSrv)

	chat := &tele.Chat{ID: int64(1)}

	telectx.EXPECT().Chat().Return(chat).Times(2)
	telectx.EXPECT().Message().Return(&tele.Message{Sender: &tele.User{IsBot: true}})

	randomReminder := random.Reminder()
	reminderSrv.SaveName(chat.ID, randomReminder.Name)

	expectedTxt := fmt.Sprintf(messages.MonthDayMessage, randomReminder.Name, "раз в месяц")
	expectedOpts := &tele.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.BackToReminderMenuBtns(),
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(msg interface{}, opts *tele.SendOptions) {
		assert.Equal(t, expectedOpts, opts)
		assert.Equal(t, expectedTxt, msg)
	}).Return(nil)

	err := controller.Month(context.Background(), telectx)
	assert.NoError(t, err)

	r, err := reminderSrv.GetFromMemory(chat.ID)
	assert.NoError(t, err)

	assert.Equal(t, model.OnceMonthType, r.Type)
}

func TestMonth_NewName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	controller := New(nil, nil, nil, reminderSrv)

	chat := &tele.Chat{ID: int64(1)}

	newName := random.String(5)

	telectx.EXPECT().Chat().Return(chat).Times(3)
	telectx.EXPECT().Message().Return(&tele.Message{Sender: &tele.User{IsBot: false}, Text: newName}).Times(2)

	randomReminder := random.Reminder()
	reminderSrv.SaveName(chat.ID, randomReminder.Name)

	expectedTxt := fmt.Sprintf(messages.MonthDayMessage, newName, "раз в месяц")
	expectedOpts := &tele.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.BackToReminderMenuBtns(),
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(msg interface{}, opts *tele.SendOptions) {
		assert.Equal(t, expectedOpts, opts)
		assert.Equal(t, expectedTxt, msg)
	}).Return(nil)

	err := controller.Month(context.Background(), telectx)
	assert.NoError(t, err)

	r, err := reminderSrv.GetFromMemory(chat.ID)
	assert.NoError(t, err)

	assert.Equal(t, model.OnceMonthType, r.Type)
	assert.Equal(t, newName, r.Name)
}

func TestYear(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	controller := New(nil, nil, nil, reminderSrv)

	chat := &tele.Chat{ID: int64(1)}

	// для календаря
	reminderSrv.SaveUser(chat.ID)

	telectx.EXPECT().Chat().Return(chat).Times(3)
	telectx.EXPECT().Message().Return(&tele.Message{Sender: &tele.User{IsBot: true}})

	randomReminder := random.Reminder()
	reminderSrv.SaveName(chat.ID, randomReminder.Name)

	reminderView := view.NewReminder()

	expectedTxt := fmt.Sprintf(messages.CalendarMessage, randomReminder.Name, "раз в год")
	expectedOpts := &tele.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: reminderView.Calendar(),
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(msg interface{}, opts *tele.SendOptions) {
		assert.Equal(t, expectedOpts, opts)
		assert.Equal(t, expectedTxt, msg)
	}).Return(nil)

	err := controller.Year(context.Background(), telectx)
	assert.NoError(t, err)

	r, err := reminderSrv.GetFromMemory(chat.ID)
	assert.NoError(t, err)

	assert.Equal(t, model.OnceYearType, r.Type)
}

func TestYear_NewName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	controller := New(nil, nil, nil, reminderSrv)

	chat := &tele.Chat{ID: int64(1)}

	// для календаря
	reminderSrv.SaveUser(chat.ID)

	newName := random.String(5)

	telectx.EXPECT().Chat().Return(chat).Times(4)
	telectx.EXPECT().Message().Return(&tele.Message{Sender: &tele.User{IsBot: false}, Text: newName}).Times(2)

	randomReminder := random.Reminder()
	reminderSrv.SaveName(chat.ID, randomReminder.Name)

	reminderView := view.NewReminder()

	expectedTxt := fmt.Sprintf(messages.CalendarMessage, newName, "раз в год")
	expectedOpts := &tele.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: reminderView.Calendar(),
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(msg interface{}, opts *tele.SendOptions) {
		assert.Equal(t, expectedOpts, opts)
		assert.Equal(t, expectedTxt, msg)
	}).Return(nil)

	err := controller.Year(context.Background(), telectx)
	assert.NoError(t, err)

	r, err := reminderSrv.GetFromMemory(chat.ID)
	assert.NoError(t, err)

	assert.Equal(t, model.OnceYearType, r.Type)
	assert.Equal(t, newName, r.Name)
}

func TestDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	controller := New(nil, nil, nil, reminderSrv)

	chat := &tele.Chat{ID: int64(1)}

	// для календаря
	reminderSrv.SaveUser(chat.ID)

	telectx.EXPECT().Chat().Return(chat).Times(3)
	telectx.EXPECT().Message().Return(&tele.Message{Sender: &tele.User{IsBot: true}})

	randomReminder := random.Reminder()
	reminderSrv.SaveName(chat.ID, randomReminder.Name)

	reminderView := view.NewReminder()

	expectedTxt := fmt.Sprintf(messages.CalendarMessage, randomReminder.Name, "дата")
	expectedOpts := &tele.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: reminderView.Calendar(),
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(msg interface{}, opts *tele.SendOptions) {
		assert.Equal(t, expectedOpts, opts)
		assert.Equal(t, expectedTxt, msg)
	}).Return(nil)

	err := controller.Date(context.Background(), telectx)
	assert.NoError(t, err)

	r, err := reminderSrv.GetFromMemory(chat.ID)
	assert.NoError(t, err)

	assert.Equal(t, model.DateType, r.Type)
}

func TestDate_NewName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reminderEditor := mocks.NewMockreminderEditor(ctrl)
	reminderSrv := reminder.New(reminderEditor)
	telectx := mocks.NewMockteleCtx(ctrl)

	controller := New(nil, nil, nil, reminderSrv)

	chat := &tele.Chat{ID: int64(1)}

	// для календаря
	reminderSrv.SaveUser(chat.ID)

	newName := random.String(5)

	telectx.EXPECT().Chat().Return(chat).Times(4)
	telectx.EXPECT().Message().Return(&tele.Message{Sender: &tele.User{IsBot: false}, Text: newName}).Times(2)

	randomReminder := random.Reminder()
	reminderSrv.SaveName(chat.ID, randomReminder.Name)

	reminderView := view.NewReminder()

	expectedTxt := fmt.Sprintf(messages.CalendarMessage, newName, "дата")
	expectedOpts := &tele.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: reminderView.Calendar(),
	}

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(msg interface{}, opts *tele.SendOptions) {
		assert.Equal(t, expectedOpts, opts)
		assert.Equal(t, expectedTxt, msg)
	}).Return(nil)

	err := controller.Date(context.Background(), telectx)
	assert.NoError(t, err)

	r, err := reminderSrv.GetFromMemory(chat.ID)
	assert.NoError(t, err)

	assert.Equal(t, model.DateType, r.Type)
	assert.Equal(t, newName, r.Name)
}
