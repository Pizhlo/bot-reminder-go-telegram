CREATE TABLE IF NOT EXISTS reminders.memory_reminders (
    id uuid,
    user_id INT unique not null,
    "text" TEXT,
    type_id int,
    date text,
    time text,
    created TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users.users(id) ON DELETE CASCADE,
    FOREIGN KEY (type_id) REFERENCES reminders.types(id) ON DELETE CASCADE
);