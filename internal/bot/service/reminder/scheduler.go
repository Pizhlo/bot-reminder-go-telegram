package reminder

import (
	"context"
	"errors"
	"fmt"
	"time"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	gocron "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/scheduler"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (c *ReminderService) CreateReminder(ctx context.Context, loc *time.Location, f gocron.Task, r *model.Reminder) (gocron.NewJob, error) {
	var sch *gocron.Scheduler
	var err error
	sch, err = c.getScheduler(r.TgID)
	if err != nil {
		err = c.CreateScheduler(ctx, r.TgID, loc, f)
		if err != nil {
			return gocron.NewJob{}, err
		}

		sch, _ = c.getScheduler(r.TgID)
	}

	logrus.Debugf(wrap(fmt.Sprintf("starting job for user %d. Reminder: %+v", r.TgID, r)))

	switch r.Type {
	case model.EverydayType:
		return sch.CreateEverydayJob(r.Time, f, ctx, r)
	case model.SeveralTimesDayType:
		if r.Date == "minutes" {
			return sch.CreateMinutesReminder(r.Time, f, ctx, r)
		}

		return sch.CreateHoursReminder(r.Time, f, ctx, r)
	case model.EveryWeekType:
		wd, err := view.ParseWeekday(r.Date)
		if err != nil {
			return gocron.NewJob{}, fmt.Errorf("error while parsing week day %s: %w", r.Date, err)
		}

		return sch.CreateEveryWeekReminder(wd, r.Time, f, ctx, r)
	case model.SeveralDaysType:
		return sch.CreateSeveralDaysReminder(r.Date, r.Time, f, ctx, r)
	case model.OnceMonthType:
		return sch.CreateMonthlyReminder(r.Date, r.Time, f, ctx, r)
	case model.OnceYearType:
		return sch.CreateOnceInYearReminder(r.Date, r.Time, f, ctx, r)
	case model.DateType:
		// создаем отложенный вызов отправки напоминания
		newJob, err := sch.CreateCalendarDateReminder(r.Date, r.Time, f, ctx, r)
		if err != nil {
			return gocron.NewJob{}, err
		}

		// создаем отложенный вызов удаления напоминания из БД
		_, err = sch.CreateCalendarDateReminder(r.Date, r.Time, c.DeleteReminder, ctx, r)
		if err != nil {
			return gocron.NewJob{}, err
		}

		return newJob, nil
	default:
		return gocron.NewJob{}, fmt.Errorf("unknown type of reminder: %s", r.Type)
	}
}

// CreateScheduler создает планировщика для конкретного пользователя.
// Запускает также все таски пользователя
func (c *ReminderService) CreateScheduler(ctx context.Context, tgID int64, loc *time.Location, f gocron.Task) error {
	val, ok := c.schedulers[tgID]
	if ok {
		err := val.StopJobs()
		if err != nil {
			return err
		}
	}

	sch, err := gocron.New(loc)
	if err != nil {
		return err
	}

	c.schedulers[tgID] = sch
	return c.StartAllJobs(ctx, tgID, loc, f)
}

func (c *ReminderService) getScheduler(tgID int64) (*gocron.Scheduler, error) {
	if val, ok := c.schedulers[tgID]; ok {
		return val, nil
	}

	return nil, errors.New("no scheduler found for this user")
}

// DeleteJob останавливает и удаляет таску в планировщике
func (c *ReminderService) DeleteJob(tgID int64, jobID uuid.UUID) error {
	if val, ok := c.schedulers[tgID]; ok {
		logrus.Debugf(wrap(fmt.Sprintf("deleting job %v from scheduler", jobID)))
		return val.DeleteJob(jobID)
	}

	return errors.New("no scheduler found for this user")
}

// startAllJobs запускает все таски
func (c *ReminderService) StartAllJobs(ctx context.Context, userID int64, loc *time.Location, f gocron.Task) error {
	reminders, err := c.reminderEditor.GetAllByUserID(ctx, userID)
	if err != nil {
		if !errors.Is(err, api_errors.ErrRemindersNotFound) {
			return err
		}
	}

	for _, r := range reminders {
		newJob, err := c.CreateReminder(ctx, loc, f, &r)
		if err != nil {
			return err
		}

		err = c.SaveJobID(ctx, newJob.JobID, userID, r.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
