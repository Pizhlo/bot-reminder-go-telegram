CREATE SCHEMA users;

CREATE TABLE IF NOT EXISTS users.users (
	id serial not null primary key,
	tg_id int not null unique
);

create table if not exists users.timezones (
	id serial not null primary key,
	user_id int not null unique,
	timezone text not null,
	foreign key (user_id) references users.users(id) on delete cascade
);

CREATE SCHEMA notes;

create table if not exists notes.notes (
	id serial not null,
	user_id int not null,
	"text" text,	
	created timestamp not null,
	primary key(id, user_id),
	foreign key (user_id) references users.users(id) on delete cascade
);

CREATE OR REPLACE FUNCTION set_note_id()
RETURNS TRIGGER AS $f$
DECLARE
max_identifier INTEGER;
BEGIN
  SELECT MAX(id)+1 INTO max_identifier
  FROM notes.notes
  WHERE user_id = NEW.user_id;
IF max_identifier IS NULL THEN
    max_identifier := 1;
END IF;
NEW.id := max_identifier;
RETURN NEW;
END;
$f$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_trigger_notes
BEFORE INSERT ON notes.notes
FOR EACH ROW
EXECUTE FUNCTION set_note_id();