package gocron

import (
	"fmt"
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

// CreateEverydayJob создает ежедневные вызовы в указанное время
func (s *Scheduler) CreateEverydayJob(userTime string, task task, params FuncParams) (uuid.UUID, error) {
	layout := "15:04"

	t, err := time.Parse(layout, userTime)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error while parsing user's time %s: %w", userTime, err)
	}

	cronTime := gocron.NewAtTimes(gocron.NewAtTime(uint(t.Hour()), uint(t.Minute()), 0))

	job := gocron.NewTask(task, params.Ctx, params.Reminder)

	j, err := s.NewJob(gocron.DailyJob(uint(1), cronTime), job)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error while creating new job: %w", err)
	}

	return j.ID(), nil
}

// DeleteJob удаляет задачу
func (s *Scheduler) DeleteJob(id uuid.UUID) {
	s.RemoveJob(id)
}
