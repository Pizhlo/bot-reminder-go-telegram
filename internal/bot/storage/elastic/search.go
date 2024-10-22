package elastic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/elastic"
	"github.com/google/uuid"
)

// Search производит поиск по переданным данным. Возвращает ID подходящих записей
func (c *client) Search(ctx context.Context, data elastic.Data) ([]uuid.UUID, error) {
	query, err := data.SearchNoteQuery()
	if err != nil {
		return nil, fmt.Errorf("error while creating query for search note: %+v", err)
	}

	res, err := c.cl.Search().
		Index(data.Index).
		Request(query).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("error searching note: %+v", err)
	}

	var ids []uuid.UUID

	for _, val := range res.Hits.Hits {
		bytesJSON, err := val.Source_.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("error marshalling JSON while searching notes: %+v", err)
		}

		var note elastic.Note
		err = json.Unmarshal(bytesJSON, &note)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON while searching notes: %+v", err)
		}

		ids = append(ids, note.ID)
	}

	return ids, nil
}
