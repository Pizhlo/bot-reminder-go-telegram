package server

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

func (s *Server) GetAllNotes(ctx context.Context, userID int) ([]model.Note, error) {
	return s.noteEditor.GetAllNotes(ctx, userID)
}

func (s *Server) SaveNote(ctx context.Context, note model.Note) error {
	return s.noteEditor.SaveNote(ctx, note)
}
