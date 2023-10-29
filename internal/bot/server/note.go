package server

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

func (s *Server) SaveNote(note model.Note) error {
	c, cancel := context.WithCancel(context.TODO()) // тот ли контекст?
	defer cancel()

	if err := s.noteEditor.SaveNote(c, note); err != nil {
		return err
	}
	return nil
}
