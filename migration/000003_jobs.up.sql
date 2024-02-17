CREATE TABLE IF NOT EXISTS reminders.jobs(
    id SERIAL NOT NULL,
    job_id UUID NOT NULL,
    reminder_id serial NOT NULL,
    user_id INT NOT NULL, 
    PRIMARY KEY (id),
    FOREIGN KEY (reminder_id, user_id) REFERENCES reminders.reminders(id, user_id) ON DELETE CASCADE
);