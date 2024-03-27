CREATE SCHEMA IF NOT EXISTS reminders;

CREATE TABLE IF NOT EXISTS reminders.types (
    id SERIAL NOT NULL,
    name text not null,
    PRIMARY KEY (id)
);

INSERT INTO reminders.types VALUES
(1, 'several_times_day'),
(2, 'everyday'),
(3, 'everyweek'),
(4, 'several_days'),
(5, 'once_month'),
(6, 'once_year'),
(7, 'date');

CREATE TABLE IF NOT EXISTS reminders.reminders (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    user_id INT NOT NULL,
    "text" TEXT,
    type_id int not null,
    date text not null,
    time text not null,
    created TIMESTAMP NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users.users(id) ON DELETE CASCADE,
    FOREIGN KEY (type_id) REFERENCES reminders.types(id) ON DELETE CASCADE
);

CREATE VIEW reminders.reminders_view AS 
SELECT id, user_id, row_number()over(partition by user_id order by id) AS reminder_number
FROM reminders.reminders;