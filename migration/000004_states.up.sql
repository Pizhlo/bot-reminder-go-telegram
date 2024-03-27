create table if not exists users.state_types(
    id SERIAL NOT NULL,
    name text not null,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS users.states(
    id SERIAL NOT NULL,
    user_id INT NOT NULL unique,
    state_id INT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users.users(id) ON DELETE CASCADE,
    FOREIGN KEY (state_id) REFERENCES users.state_types(id) ON DELETE CASCADE
);

INSERT INTO users.state_types VALUES
(1, 'default'),
(2, 'start'),
(3, 'list_note'),
(4, 'create_note'),
(5, 'days_duration'),
(6, 'hours'),
(7, 'list_reminder'),
(8, 'location'),
(9, 'minutes_duration'),
(10, 'month'),
(11, 'date_reminder'),
(12, 'reminder_name'),
(13, 'reminder_time'),
(14, 'search_note_by_text'),
(15, 'search_note_by_date'),
(16, 'search_note_by_two_dates'),
(17, 'several_days'),
(18, 'several_times_a_day'),
(19, 'every_week'),
(20, 'every_year');