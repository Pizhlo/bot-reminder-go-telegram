package note

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/sirupsen/logrus"
)

type NoteService struct {
	noteEditor noteEditor
	logger     *logrus.Logger
}

type noteEditor interface {
	Save(ctx context.Context, note model.Note) error
}

func New(noteEditor noteEditor) *NoteService {
	return &NoteService{noteEditor: noteEditor, logger: logger.New()}
}
