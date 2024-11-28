package sharedaccess

import (
	"context"
	"fmt"
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

func (s *SharedSpace) Save(ctx context.Context, space model.SharedSpace) error {
	return s.storage.Save(ctx, space)
}

func (s *SharedSpace) SaveNote(ctx context.Context, note model.Note) error {
	spaceID := s.viewsMap[note.Creator.TGID].CurrentSpaceID()

	note.Space.ID = spaceID

	loc, err := s.userSrv.GetLocation(ctx, note.Creator.TGID)
	if err != nil {
		return fmt.Errorf("error whiel getting user's location: %+v", err)
	}

	note.Created = time.Now().In(loc)

	return s.storage.SaveNote(ctx, note)
}
