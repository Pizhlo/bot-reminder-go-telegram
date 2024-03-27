CREATE TABLE IF NOT EXISTS users.states(
    id SERIAL NOT NULL,
    user_id INT NOT NULL unique,
    state text not null,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users.users(id) ON DELETE CASCADE
);