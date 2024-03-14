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

CREATE VIEW notes.notes_view AS 
SELECT id, user_id, row_number()over(partition by user_id order by id) AS note_number
FROM notes.notes;