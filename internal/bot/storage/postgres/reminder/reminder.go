package reminder

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type ReminderRepo struct {
	*pgxpool.Pool
}

func New(conn *pgxpool.Pool) *ReminderRepo {
	return &ReminderRepo{conn}
}
