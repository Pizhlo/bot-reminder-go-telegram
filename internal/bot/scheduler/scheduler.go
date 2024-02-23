package gocron

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"gopkg.in/telebot.v3"
)

// task - функция, которую будут вызывать в момент срабатывания напоминания
type task func(ctx telebot.Context, reminder model.Reminder) error

// Scheduler управляет отложенными вызовами
type Scheduler struct {
	gocron.Scheduler
}

func New() (*Scheduler, error) {
	sch, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	sch.Start()

	return &Scheduler{sch}, nil
}

type FuncParams struct {
	Ctx      telebot.Context
	Reminder model.Reminder
}

type NextRun struct {
	JobID   uuid.UUID
	NextRun time.Time
}

// CreateEverydayJob создает ежедневные вызовы в указанное время
func (s *Scheduler) CreateEverydayJob(userTime string, task task, params FuncParams) (NextRun, error) {
	cronTime, err := s.makeTime(userTime)
	if err != nil {
		return NextRun{}, fmt.Errorf("error while creating cron time: %w", err)
	}

	job := gocron.NewTask(task, params.Ctx, params.Reminder)

	dailyJob := gocron.DailyJob(uint(1), cronTime)

	j, err := s.NewJob(dailyJob, job)
	if err != nil {
		return NextRun{}, fmt.Errorf("error while creating new job: %w", err)
	}

	run, err := j.NextRun()
	if err != nil {
		return NextRun{}, fmt.Errorf("error while getting next run: %w", err)
	}

	result := NextRun{
		JobID:   j.ID(),
		NextRun: run,
	}

	return result, nil
}

// DeleteJob удаляет задачу
func (s *Scheduler) DeleteJob(id uuid.UUID) error {
	return s.RemoveJob(id)
}

// CreateMinutesReminder создает напоминание один раз в несколько минут
func (s *Scheduler) CreateMinutesReminder(minutes string, task task, params FuncParams) (NextRun, error) {
	job := gocron.NewTask(task, params.Ctx, params.Reminder)

	minutesInt, err := strconv.Atoi(minutes)
	if err != nil {
		return NextRun{}, err
	}

	d := time.Minute * time.Duration(minutesInt)

	j, err := s.NewJob(gocron.DurationJob(d), job)
	if err != nil {
		return NextRun{}, fmt.Errorf("error while creating new job: %w", err)
	}

	run, err := j.NextRun()
	if err != nil {
		return NextRun{}, fmt.Errorf("error while getting next run: %w", err)
	}

	result := NextRun{
		JobID:   j.ID(),
		NextRun: run,
	}

	return result, nil

}

// CreateHoursReminder создает напоминание один раз в несколько часов
func (s *Scheduler) CreateHoursReminder(hours string, task task, params FuncParams) (NextRun, error) {
	job := gocron.NewTask(task, params.Ctx, params.Reminder)

	hoursInt, err := strconv.Atoi(hours)
	if err != nil {
		return NextRun{}, err
	}

	d := time.Hour * time.Duration(hoursInt)

	j, err := s.NewJob(gocron.DurationJob(d), job)
	if err != nil {
		return NextRun{}, fmt.Errorf("error while creating new job: %w", err)
	}

	run, err := j.NextRun()
	if err != nil {
		return NextRun{}, fmt.Errorf("error while getting next run: %w", err)
	}

	result := NextRun{
		JobID:   j.ID(),
		NextRun: run,
	}

	return result, nil
}

// CreateEveryWeekReminder создает напоминание еженедельное напоминание
func (s *Scheduler) CreateEveryWeekReminder(weekDay time.Weekday, userTime string, task task, params FuncParams) (NextRun, error) {
	job := gocron.NewTask(task, params.Ctx, params.Reminder)

	cronTime, err := s.makeTime(userTime)
	if err != nil {
		return NextRun{}, fmt.Errorf("error while creating cron time: %w", err)
	}

	j, err := s.NewJob(gocron.WeeklyJob(0, gocron.NewWeekdays(weekDay), cronTime), job)
	if err != nil {
		return NextRun{}, fmt.Errorf("error while creating new job: %w", err)
	}

	run, err := j.NextRun()
	if err != nil {
		return NextRun{}, fmt.Errorf("error while getting next run: %w", err)
	}

	result := NextRun{
		JobID:   j.ID(),
		NextRun: run,
	}

	return result, nil
}

// CreateSeveralDaysReminder создает напоминание раз в несколько дней
func (s *Scheduler) CreateSeveralDaysReminder(days string, userTime string, task task, params FuncParams) (NextRun, error) {
	job := gocron.NewTask(task, params.Ctx, params.Reminder)

	cronTime, err := s.makeTime(userTime)
	if err != nil {
		return NextRun{}, fmt.Errorf("error while creating cron time: %w", err)
	}

	daysInt, err := strconv.Atoi(days)
	if err != nil {
		return NextRun{}, err
	}

	j, err := s.NewJob(gocron.DailyJob(uint(daysInt), cronTime), job)
	if err != nil {
		return NextRun{}, err
	}

	run, err := j.NextRun()
	if err != nil {
		return NextRun{}, fmt.Errorf("error while getting next run: %w", err)
	}

	result := NextRun{
		JobID:   j.ID(),
		NextRun: run,
	}

	return result, nil
}

// CreateMonthlyReminder создает напоминание раз в месяц
func (s *Scheduler) CreateMonthlyReminder(days string, userTime string, task task, params FuncParams) (NextRun, error) {
	day, err := strconv.Atoi(days)
	if err != nil {
		return NextRun{}, err
	}

	cronTime, err := s.makeTime(userTime)
	if err != nil {
		return NextRun{}, fmt.Errorf("error while creating cron time: %w", err)
	}

	job := gocron.NewTask(task, params.Ctx, params.Reminder)

	j, err := s.NewJob(gocron.MonthlyJob(uint(0), gocron.NewDaysOfTheMonth(day), cronTime), job)
	if err != nil {
		return NextRun{}, err
	}

	run, err := j.NextRun()
	if err != nil {
		return NextRun{}, fmt.Errorf("error while getting next run: %w", err)
	}

	result := NextRun{
		JobID:   j.ID(),
		NextRun: run,
	}

	return result, nil
}

func (s *Scheduler) CreateOnceInYearReminder(date, userTime string, task task, params FuncParams) (NextRun, error) {
	job := gocron.NewTask(task, params.Ctx, params.Reminder)

	cronTab := s.makeCronTab(date, userTime)

	j, err := s.NewJob(gocron.CronJob(cronTab, false), job)
	if err != nil {
		return NextRun{}, fmt.Errorf("error while creating job: %w", err)
	}

	run, err := j.NextRun()
	if err != nil {
		return NextRun{}, fmt.Errorf("error while getting next run: %w", err)
	}

	result := NextRun{
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

	return gocron.NewAtTimes(gocron.NewAtTime(uint(t.Hour()), uint(t.Minute()), 0)), nil
}
