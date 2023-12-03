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
	id serial not null primary key,
	user_id int not null,
	"text" text,
	foreign key (user_id) references users.users(id) on delete cascade,
	created timestamp not null,
	unique(id, user_id)
);

-- Создание триггера для заполнения таблицы "users.timezones" при внесении записи в "users.users"
-- CREATE OR REPLACE FUNCTION insert_into_timezone() RETURNS TRIGGER AS $$
-- BEGIN
--     IF NEW.timezone <> '' THEN
--         INSERT INTO users.timezones (user_id, timezone) VALUES (NEW.user_id, NEW.timezone);
--     END IF;
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;