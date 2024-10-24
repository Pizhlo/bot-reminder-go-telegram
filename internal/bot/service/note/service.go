package note

import (
	"context"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// NoteService отвечает за заметки: удаление, создание, выдача
type NoteService struct {
	noteEditor noteEditor
	viewsMap   map[int64]*view.NoteView
	searchMap  map[int64]model.SearchByTwoDates
}

//go:generate mockgen -source ./service.go -destination=../../mocks/note_srv.go -package=mocks
type noteEditor interface {
	// Save сохраняет заметку в базе данных. Для сохранения требуется: ID пользователя, содержимое заметки, дата создания
	Save(ctx context.Context, note model.Note) error

	// GetAllByUserID достает из базы все заметки пользователя по ID
	GetAllByUserID(ctx context.Context, userID int64) ([]model.Note, error)

	// DeleteAllByUserID удаляет все заметки пользователя по user ID
	DeleteAllByUserID(ctx context.Context, userID int64) error

	// DeleteByID удаляет одну заметку. Для удаления необходим ID заметки
	DeleteByID(ctx context.Context, id uuid.UUID) error

	// GetByID возвращает заметку с переданным ID. Если такой заметки нет, возвращает ErrNotesNotFound
	GetByViewID(ctx context.Context, userID int64, viewID int) (*model.Note, error)

	// SearchByText производит поиск по заметок по тексту. Если таких заметок нет, возвращает ErrNotesNotFound
	SearchByText(ctx context.Context, searchNote model.SearchByText) ([]model.Note, error)

	// SearchByOneDate производит поиск по заметок по выбранной дате. Если таких заметок нет, возвращает ErrNotesNotFound
	SearchByOneDate(ctx context.Context, searchNote model.SearchByOneDate) ([]model.Note, error)

	// SearchByTwoDates производит поиск по заметок по двум датам. Если таких заметок нет, возвращает ErrNotesNotFound
	SearchByTwoDates(ctx context.Context, searchNote *model.SearchByTwoDates) ([]model.Note, error)

	// UpdateNote обновляет заметку на переданную
	UpdateNote(ctx context.Context, note model.EditNote) error
}

func New(noteEditor noteEditor) *NoteService {
	return &NoteService{
		noteEditor: noteEditor,
		viewsMap:   make(map[int64]*view.NoteView),
		searchMap:  make(map[int64]model.SearchByTwoDates),
	}
}

// SaveUser сохраняет пользователя в мапе view
func (n *NoteService) SaveUser(userID int64) {
	if _, ok := n.viewsMap[userID]; !ok {
		logrus.Debugf(wrap(fmt.Sprintf("user %d not found in the views map. Saving...\n", userID)))
		n.viewsMap[userID] = view.NewNote()
	} else {
		logrus.Debugf(wrap(fmt.Sprintf("user %d already saved in the views map.\n", userID)))
	}

}

// SetupCalendar устанавливает месяц и год в календаре на текущие
func (n *NoteService) SetupCalendar(userID int64) {
	n.viewsMap[userID].SetCurMonth()
	n.viewsMap[userID].SetCurYear()
}

func wrap(s string) string {
	return fmt.Sprintf("Note service: %s", s)
}
