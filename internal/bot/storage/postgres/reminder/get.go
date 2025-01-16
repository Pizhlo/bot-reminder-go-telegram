package reminder

import (
	"context"
	"fmt"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/google/uuid"
)

// GetAllByUserID достает из базы все напоминания пользователя по ID
func (db *ReminderRepo) GetAllByUserID(ctx context.Context, userID int64) ([]model.Reminder, error) {
	reminders := make([]model.Reminder, 0)

	rows, err := db.db.QueryContext(ctx, `select reminders.reminders.id, reminder_number, tg_id, text, created, date, time, name as type, reminders.jobs.job_id
	from reminders.reminders
		join reminders.types on reminders.types.id = reminders.reminders.type_id
		join users.users on users.id = reminders.user_id
		join reminders.reminders_view on reminders.reminders_view.id = reminders.reminders.id
		join reminders.jobs on reminders.jobs.reminder_id = reminders.reminders.id
		where reminders.reminders.user_id = (select id from users.users where tg_id = $1)
		order by created ASC`, userID)

	if err != nil {
		return nil, fmt.Errorf("error while getting all reminders from DB by user ID %d: %w", userID, err)
	}
	defer rows.Close()

	for rows.Next() {
		reminder := model.Reminder{}

		err := rows.Scan(&reminder.ID, &reminder.ViewID, &reminder.TgID, &reminder.Name, &reminder.Created, &reminder.Date, &reminder.Time, &reminder.Type, &reminder.Job.ID)
		if err != nil {
			return nil, fmt.Errorf("error while scanning reminder: %w", err)
		}

		reminders = append(reminders, reminder)
	}

	if len(reminders) == 0 {
		return nil, api_errors.ErrRemindersNotFound
	}

	return reminders, nil
}

func (db *ReminderRepo) GetByViewID(ctx context.Context, userID int64, viewID int) (*model.Reminder, error) {
	reminder := &model.Reminder{}

	row := db.db.QueryRowContext(ctx, `select reminders.reminders.id, reminder_number, tg_id, text, created, date, time, name as type, reminders.jobs.job_id
	from reminders.reminders
		join reminders.types on reminders.types.id = reminders.reminders.type_id
		join users.users on users.id = reminders.user_id
		join reminders.reminders_view on reminders.reminders_view.id = reminders.reminders.id
		join reminders.jobs on reminders.jobs.reminder_id = reminders.reminders.id
		where reminders.reminders.user_id = (select id from users.users where tg_id = $1)
		and reminders.reminders_view.reminder_number = $2
		order by created ASC;`, userID, viewID)

	err := row.Scan(&reminder.ID, &reminder.ViewID, &reminder.TgID, &reminder.Name, &reminder.Created, &reminder.Date, &reminder.Time, &reminder.Type, &reminder.Job.ID)
	if err != nil {
		return &model.Reminder{}, fmt.Errorf("error while scanning reminder while getting by view ID: %w", err)
	}

	return reminder, nil
}

func (db *ReminderRepo) GetAllJobs(ctx context.Context, userID int64) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, 0)

	rows, err := db.db.QueryContext(ctx, `select job_id from reminders.jobs 
	where reminder_id in (select id from reminders.reminders 
	where user_id = (select id from users.users where tg_id = $1))`, userID)

	if err != nil {
		return nil, fmt.Errorf("error while getting all jobs from DB by user ID %d while getting all jobs: %w", userID, err)
	}
	defer rows.Close()

	for rows.Next() {
		id := uuid.UUID{}

		err := rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("error while scanning reminder: %w", err)
		}

		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return nil, api_errors.ErrRemindersNotFound
	}

	return ids, nil
}

func (db *ReminderRepo) GetJobID(ctx context.Context, reminderID uuid.UUID) (uuid.UUID, error) {
	id := uuid.UUID{}

	row := db.db.QueryRowContext(ctx, `select job_id from reminders.jobs where reminder_id = $1`, reminderID)

	err := row.Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error while scanning job ID: %w", err)
	}

	return id, nil
}

func (db *ReminderRepo) GetReminderID(ctx context.Context, userID int64, viewID int) (uuid.UUID, error) {
	id := uuid.UUID{}

	row := db.db.QueryRowContext(ctx, `select id from reminders.reminders_view where user_id = (select id from users.users where tg_id = $1) and reminder_number = $2`, userID, viewID)

	err := row.Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error while scanning reminder ID: %w", err)
	}

	return id, nil
}

func (db *ReminderRepo) GetMemory(ctx context.Context) ([]model.Reminder, error) {
	rows, err := db.db.QueryContext(ctx, `select 
	reminders.memory_reminders.id, tg_id, text, created, date, time, name as type from reminders.memory_reminders 
	left join reminders.types on reminders.types.id = reminders.memory_reminders.type_id 
	join users.users on users.id = memory_reminders.user_id;`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var res []model.Reminder
	for rows.Next() {
		var r model.Reminder

		err := rows.Scan(&r.ID, &r.TgID, &r.Name, &r.Created, &r.Date, &r.Time, &r.TypeString)
		if err != nil {
			return nil, fmt.Errorf("error while scanning memory reminder: %v", err)
		}

		if r.TypeString.Valid {
			r.Type = model.ReminderType(r.TypeString.String)
		}

		res = append(res, r)
	}

	return res, nil
}

// GetAllByUserIDFromSharedSpaces возвращает все напоминания пользователя, созданные в совместных пространствах
func (db *ReminderRepo) GetAllByUserIDFromSharedSpaces(ctx context.Context, userID int64) ([]model.Reminder, error) {
	reminders := make([]model.Reminder, 0)

	rows, err := db.db.QueryContext(ctx, `select shared_spaces.reminders.id, reminder_number, tg_id, text, created, date, time, name as type, shared_spaces.jobs.job_id
	from shared_spaces.reminders
		join reminders.types on reminders.types.id = shared_spaces.reminders.type_id
		join users.users on users.id = reminders.user_id
		join shared_spaces.reminders_view on shared_spaces.reminders_view.id = shared_spaces.reminders.id
		join shared_spaces.jobs on shared_spaces.jobs.reminder_id = shared_spaces.reminders.id
		where shared_spaces.reminders.user_id = (select id from users.users where tg_id = $1)
		order by created ASC;`, userID)

	if err != nil {
		return nil, fmt.Errorf("error while getting all reminders from DB by user ID %d from shared spaces: %w", userID, err)
	}
	defer rows.Close()

	for rows.Next() {
		reminder := model.Reminder{}

		err := rows.Scan(&reminder.ID, &reminder.ViewID, &reminder.TgID, &reminder.Name, &reminder.Created, &reminder.Date, &reminder.Time, &reminder.Type, &reminder.Job.ID)
		if err != nil {
			return nil, fmt.Errorf("error while scanning reminder from shared space: %w", err)
		}

		reminders = append(reminders, reminder)
	}

	if len(reminders) == 0 {
		return nil, api_errors.ErrRemindersNotFound
	}

	return reminders, nil
}
