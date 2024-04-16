package gocron

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Task - функция для отложенного вызова
type Task func(ctx context.Context, reminder *model.Reminder) error

// Scheduler управляет отложенными вызовами
type Scheduler struct {
	loc *time.Location // часовой пояс пользователя
	scheduelrInterface
}

//go:generate mockgen -source ./scheduler.go -destination=./mocks/schediler.go
type scheduelrInterface interface {
	gocron.Scheduler
}

func New(loc *time.Location) (*Scheduler, error) {
	locOption := gocron.WithLocation(loc)

	sch, err := gocron.NewScheduler(locOption)
	if err != nil {
		return nil, err
	}

	sch.Start()

	return &Scheduler{loc, sch}, nil
}

type NewJob struct {
	JobID   uuid.UUID
	NextRun time.Time
}

// CreateEverydayJob создает ежедневные вызовы в указанное время
func (s *Scheduler) CreateEverydayJob(userTime string, task Task, params ...any) (NewJob, error) {
	cronTime, err := s.makeTime(userTime)
	if err != nil {
		return NewJob{}, fmt.Errorf("error while creating cron time: %w", err)
	}

	job := makeTask(task, params...)

	dailyJob := gocron.DailyJob(uint(1), cronTime)

	j, err := s.NewJob(dailyJob, job)
	if err != nil {
		return NewJob{}, fmt.Errorf("error while creating new job: %w", err)
	}

	run, err := j.NextRun()
	if err != nil {
		return NewJob{}, fmt.Errorf("error while getting next run: %w", err)
	}

	newJob := NewJob{
		JobID:   j.ID(),
		NextRun: run,
	}

	return newJob, nil
}

// DeleteJob удаляет задачу
func (s *Scheduler) DeleteJob(id uuid.UUID) error {
	return s.RemoveJob(id)
}

func makeTask(task Task, params ...any) gocron.Task {
	return gocron.NewTask(task, params...)
}

// CreateMinutesReminder создает напоминание один раз в несколько минут
func (s *Scheduler) CreateMinutesReminder(minutes string, task Task, params ...any) (NewJob, error) {
	job := makeTask(task, params...)

	minutesInt, err := strconv.Atoi(minutes)
	if err != nil {
		return NewJob{}, err
	}

	d := time.Minute * time.Duration(minutesInt)

	j, err := s.NewJob(gocron.DurationJob(d), job)
	if err != nil {
		return NewJob{}, fmt.Errorf("error while creating new job: %w", err)
	}

	run, err := j.NextRun()
	if err != nil {
		return NewJob{}, fmt.Errorf("error while getting next run: %w", err)
	}

	result := NewJob{
		JobID:   j.ID(),
		NextRun: run,
	}

	return result, nil

}

// CreateHoursReminder создает напоминание один раз в несколько часов
func (s *Scheduler) CreateHoursReminder(hours string, task Task, params ...any) (NewJob, error) {
	job := makeTask(task, params...)

	hoursInt, err := strconv.Atoi(hours)
	if err != nil {
		return NewJob{}, err
	}

	d := time.Hour * time.Duration(hoursInt)

	j, err := s.NewJob(gocron.DurationJob(d), job)
	if err != nil {
		return NewJob{}, fmt.Errorf("error while creating new job: %w", err)
	}

	run, err := j.NextRun()
	if err != nil {
		return NewJob{}, fmt.Errorf("error while getting next run: %w", err)
	}

	result := NewJob{
		JobID:   j.ID(),
		NextRun: run,
	}

	return result, nil
}

// CreateEveryWeekReminder создает напоминание еженедельное напоминание
func (s *Scheduler) CreateEveryWeekReminder(weekDay time.Weekday, userTime string, task Task, params ...any) (NewJob, error) {
	job := makeTask(task, params...)

	cronTime, err := s.makeTime(userTime)
	if err != nil {
		return NewJob{}, fmt.Errorf("error while creating cron time: %w", err)
	}

	logrus.Errorf("weekday: %+v, user time: %s, cronTime: %+v", weekDay, userTime, cronTime)

	j, err := s.NewJob(gocron.WeeklyJob(1, gocron.NewWeekdays(weekDay), cronTime), job)
	if err != nil {
		return NewJob{}, fmt.Errorf("error while creating new job: %w", err)
	}

	run, err := j.NextRun()
	if err != nil {
		return NewJob{}, fmt.Errorf("error while getting next run: %w", err)
	}

	result := NewJob{
		JobID:   j.ID(),
		NextRun: run,
	}

	return result, nil
}

// CreateSeveralDaysReminder создает напоминание раз в несколько дней
func (s *Scheduler) CreateSeveralDaysReminder(days string, userTime string, task Task, params ...any) (NewJob, error) {
	job := makeTask(task, params...)

	cronTime, err := s.makeTime(userTime)
	if err != nil {
		return NewJob{}, fmt.Errorf("error while creating cron time: %w", err)
	}

	daysInt, err := strconv.Atoi(days)
	if err != nil {
		return NewJob{}, err
	}

	j, err := s.NewJob(gocron.DailyJob(uint(daysInt), cronTime), job)
	if err != nil {
		return NewJob{}, err
	}

	run, err := j.NextRun()
	if err != nil {
		return NewJob{}, fmt.Errorf("error while getting next run: %w", err)
	}

	result := NewJob{
		JobID:   j.ID(),
		NextRun: run,
	}

	return result, nil
}

// CreateMonthlyReminder создает напоминание раз в месяц
func (s *Scheduler) CreateMonthlyReminder(days string, userTime string, task Task, params ...any) (NewJob, error) {
	day, err := strconv.Atoi(days)
	if err != nil {
		return NewJob{}, err
	}

	cronTime, err := s.makeTime(userTime)
	if err != nil {
		return NewJob{}, fmt.Errorf("error while creating cron time: %w", err)
	}

	job := makeTask(task, params...)

	j, err := s.NewJob(gocron.MonthlyJob(uint(1), gocron.NewDaysOfTheMonth(day), cronTime), job)
	if err != nil {
		return NewJob{}, err
	}

	run, err := j.NextRun()
	if err != nil {
		return NewJob{}, fmt.Errorf("error while getting next run: %w", err)
	}

	result := NewJob{
		JobID:   j.ID(),
		NextRun: run,
	}

	return result, nil
}

func (s *Scheduler) CreateOnceInYearReminder(date, userTime string, task Task, params ...any) (NewJob, error) {
	job := makeTask(task, params...)

	cronTab := s.makeCronTab(date, userTime)

	j, err := s.NewJob(gocron.CronJob(cronTab, false), job)
	if err != nil {
		return NewJob{}, fmt.Errorf("error while creating job: %w", err)
	}

	run, err := j.NextRun()
	if err != nil {
		return NewJob{}, fmt.Errorf("error while getting next run: %w", err)
	}

	result := NewJob{
		JobID:   j.ID(),
		NextRun: run,
	}

	return result, nil
}

func (s *Scheduler) CreateCalendarDateReminder(date, userTime string, task Task, params ...any) (NewJob, error) {
	job := makeTask(task, params...)

	dates := strings.Split(date, ".")

	year := dates[2]
	month := dates[1]
	day := dates[0]

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return NewJob{}, fmt.Errorf("error while converting string year %s to int: %w", year, err)
	}

	monthInt, err := strconv.Atoi(month)
	if err != nil {
		return NewJob{}, fmt.Errorf("error while converting string month %s to int: %w", month, err)
	}

	dayInt, err := strconv.Atoi(day)
	if err != nil {
		return NewJob{}, fmt.Errorf("error while converting string day %s to int: %w", day, err)
	}

	minuteHour := strings.Split(userTime, ":")
	minute := minuteHour[1]
	hour := minuteHour[0]

	minuteInt, err := strconv.Atoi(minute)
	if err != nil {
		return NewJob{}, fmt.Errorf("error while converting string minute %s to int: %w", minute, err)
	}

	hourInt, err := strconv.Atoi(hour)
	if err != nil {
		return NewJob{}, fmt.Errorf("error while converting string hour %s to int: %w", hour, err)
	}

	timeDate := time.Date(yearInt, time.Month(monthInt), dayInt, hourInt, minuteInt, 0, 0, s.loc)

	oneTime := gocron.OneTimeJobStartDateTime(timeDate)

	j, err := s.NewJob(gocron.OneTimeJob(oneTime), job)
	if err != nil {
		return NewJob{}, fmt.Errorf("error while creating job: %w", err)
	}

	run, err := j.NextRun()
	if err != nil {
		return NewJob{}, fmt.Errorf("error while getting next run: %w", err)
	}

	result := NewJob{
		JobID:   j.ID(),
		NextRun: run,
	}

	return result, nil
}

func (s *Scheduler) makeCronTab(date, userTime string) string {
	// minute = field(fields[1], minutes)
	// hour = field(fields[2], hours)
	// dayofmonth = field(fields[3], dom)
	// month = field(fields[4], months)
	// dayofweek = field(fields[5], dow)

	minuteHour := strings.Split(userTime, ":") // 11:12

	minute := minuteHour[1]
	hours := minuteHour[0]

	dateSlice := strings.Split(date, ".") // 12.02.2024

	day := dateSlice[0]
	month := dateSlice[1]

	return fmt.Sprintf("%s %s %s %s *", minute, hours, day, month)

}

// makeTime принимает на вход строку вида "13:10" и возвращает gocron.AtTimes
func (s *Scheduler) makeTime(userTime string) (gocron.AtTimes, error) {
	layout := "15:04"

	t, err := time.Parse(layout, userTime)
	if err != nil {
		return nil, fmt.Errorf("error while parsing user's time %s: %w", userTime, err)
	}
	//tLoc := t.In(loc)

	return gocron.NewAtTimes(gocron.NewAtTime(uint(t.Hour()), uint(t.Minute()), 0)), nil
}
