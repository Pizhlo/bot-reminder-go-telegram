package gocron

import (
	"fmt"
	"strconv"
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
	layout := "15:04"

	t, err := time.Parse(layout, userTime)
	if err != nil {
		return NextRun{}, fmt.Errorf("error while parsing user's time %s: %w", userTime, err)
	}

	cronTime := gocron.NewAtTimes(gocron.NewAtTime(uint(t.Hour()), uint(t.Minute()), 0))

	job := gocron.NewTask(task, params.Ctx, params.Reminder)

	j, err := s.NewJob(gocron.DailyJob(uint(1), cronTime), job)
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
