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
	unique(id, user_id),
	primary key(id, user_id),
	foreign key (user_id) references users.users(id) on delete cascade,
);

CREATE SCHEMA notes;

CREATE TABLE IF NOT EXISTS notes.notes (
    id SERIAL NOT NULL,
    user_id INT NOT NULL,
    "text" TEXT,
    created TIMESTAMP NOT NULL,
    PRIMARY KEY (id, user_id),
    FOREIGN KEY (user_id) REFERENCES users.users(id) ON DELETE CASCADE
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