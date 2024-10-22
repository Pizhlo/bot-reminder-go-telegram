package elastic

import (
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/elastic"
)

func validateNote(note elastic.Data) error {
	if note.Index != elastic.NoteIndex {
		return fmt.Errorf("index name not equal `notes`: %+v", note.Index)
	}

	_, ok := note.Model.(elastic.Note)
	if !ok {
		return fmt.Errorf("type of field note.Model is not elastic.Note: %+v", note.Model)
	}

	return nil
}
