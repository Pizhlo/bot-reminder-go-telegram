package note

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/sirupsen/logrus"
)

// NoteService отвечает за заметки: удаление, создание, выдача
type NoteService struct {
	noteEditor noteEditor
	logger     *logrus.Logger
	viewsMap   map[int64]*view.View
}

type noteEditor interface {
	// Save сохраняет заметку в базе данных. Для сохранения требуется: ID пользователя, содержимое заметки, дата создания
	Save(ctx context.Context, note model.Note) error

	// GetAllByUserID достает из базы все заметки пользователя по ID
	GetAllByUserID(ctx context.Context, userID int64) ([]model.Note, error)

	// DeleteAllByUserID удаляет все заметки пользователя по user ID
	DeleteAllByUserID(ctx context.Context, userID int64) error

	// DeleteNoteByID удаляет одну заметку. Для удаления необходим ID заметки и пользователя
	DeleteNoteByID(ctx context.Context, userID int64, noteID int) error

	// GetByID возвращает заметку с переданным ID. Если такой заметки нет, возвращает ErrNotesNotFound
	GetByID(ctx context.Context, userID int64, noteID int) (*model.Note, error)
}

func New(noteEditor noteEditor) *NoteService {
	return &NoteService{noteEditor: noteEditor, logger: logger.New(), viewsMap: make(map[int64]*view.View)}
}

// SaveUser сохраняет пользователя в мапе view
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
