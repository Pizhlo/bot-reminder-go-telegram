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
    id SERIAL NOT NULL,
    user_id INT NOT NULL,
    "text" TEXT,
    type int not null,
    date text not null,
    time text not null,
    created TIMESTAMP NOT NULL,
    PRIMARY KEY (id),
    unique(id, user_id),
    FOREIGN KEY (user_id) REFERENCES users.users(id) ON DELETE CASCADE,
    FOREIGN KEY (type) REFERENCES reminders.types(id) ON DELETE CASCADE
);