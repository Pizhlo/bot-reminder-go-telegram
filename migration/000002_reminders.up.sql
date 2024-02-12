CREATE SCHEMA reminders;

CREATE TABLE IF NOT EXISTS reminders.reminders (
    id SERIAL NOT NULL,
    user_id INT NOT NULL,
    "text" TEXT,
    created TIMESTAMP NOT NULL,
    type text not null,
    date text not null,
    time text not null,
    PRIMARY KEY (id, user_id),
    FOREIGN KEY (user_id) REFERENCES users.users(id) ON DELETE CASCADE
);