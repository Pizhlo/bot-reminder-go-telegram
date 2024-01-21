package note

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/sirupsen/logrus"
)

type NoteService struct {
	noteEditor noteEditor
	logger     *logrus.Logger
	viewsMap   map[int64]*view.View
}

type noteEditor interface {
	Save(ctx context.Context, note model.Note) error
	GetAllByUserID(ctx context.Context, userID int64) ([]model.Note, error)
}

func New(noteEditor noteEditor) *NoteService {
	return &NoteService{noteEditor: noteEditor, logger: logger.New(), viewsMap: make(map[int64]*view.View)}
}

func (n *NoteService) SaveUser(userID int64) {
	n.logger.Debugf("Note service: checking if user saved in the views map...\n")
	if _, ok := n.viewsMap[userID]; !ok {
		n.logger.Debugf("Note service: user not found in the views map. Saving...\n")
		n.viewsMap[userID] = view.New()
	} else {
		n.logger.Debugf("Note service: user already saved in the views map.\n")
	}

	n.logger.Debugf("Note service: successfully saved user in the views map.\n")
}
