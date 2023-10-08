package user

import "github.com/jackc/pgx/v4/pgxpool"

type UserRepo struct {
	*pgxpool.Pool
}

func New(conn *pgxpool.Pool) *UserRepo {
	return &UserRepo{conn}
}
