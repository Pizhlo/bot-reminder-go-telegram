CREATE TABLE IF NOT EXISTS reminders.jobs(
    id SERIAL NOT NULL,
    job_id UUID NOT NULL,
    reminder_id serial NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (reminder_id) REFERENCES reminders.jobs(id) ON DELETE CASCADE
);