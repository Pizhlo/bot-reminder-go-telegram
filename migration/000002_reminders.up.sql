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
    type_id int not null,
    date text not null,
    time text not null,
    created TIMESTAMP NOT NULL,
    PRIMARY KEY (id, user_id),
    FOREIGN KEY (user_id) REFERENCES users.users(id) ON DELETE CASCADE,
    FOREIGN KEY (type_id) REFERENCES reminders.types(id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION set_reminder_id()
RETURNS TRIGGER AS $f$
DECLARE
max_identifier INTEGER;
BEGIN
  SELECT MAX(id)+1 INTO max_identifier
  FROM reminders.reminders
  WHERE user_id = NEW.user_id;
IF max_identifier IS NULL THEN
    max_identifier := 1;
END IF;
NEW.id := max_identifier;
RETURN NEW;
END;
$f$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_trigger_notes
BEFORE INSERT ON reminders.reminders
FOR EACH ROW
EXECUTE FUNCTION set_reminder_id();