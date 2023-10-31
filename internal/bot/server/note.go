package server

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

func (s *Server) SearchNotesByText(ctx context.Context, query model.SearchNote) ([]model.Note, error) {
	userID, err := s.GetUserID(ctx, query.TgID)
	if err != nil {
		return nil, err
	}

	query.UserID = userID

	return s.noteEditor.SearchNotesByText(ctx, query)
}

func (s *Server) GetAllNotes(ctx context.Context, tgID int64) ([]model.Note, error) {
	userID, err := s.GetUserID(ctx, tgID)
	if err != nil {
		return nil, err
	}

	return s.noteEditor.GetAllNotes(ctx, userID)
}

func (s *Server) SaveNote(ctx context.Context, note model.Note) error {
	userID, err := s.GetUserID(ctx, note.TgID)
	if err != nil {
		return err
	}

	note.UserID = userID

	return s.noteEditor.SaveNote(ctx, note)
}
