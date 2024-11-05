package sharedaccess

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

func (s *SharedSpace) Save(ctx context.Context, space model.SharedSpace) error {
	return s.storage.Save(ctx, space)
}
