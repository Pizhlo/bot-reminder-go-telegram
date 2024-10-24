package elastic

import (
	"context"
	"encoding/json"
	"fmt"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/elastic"
	"github.com/google/uuid"
)

// Search производит поиск по переданным данным. Возвращает ID подходящих записей
func (c *client) SearchByText(ctx context.Context, data elastic.Data) ([]uuid.UUID, error) {
	_, err := data.ValidateNote()
	if err != nil {
		return nil, err
	}

	query, err := data.SearchByTextQuery()
	if err != nil {
		return nil, fmt.Errorf("error while creating query for search note: %+v", err)
	}

	res, err := c.cl.Search().
		Index(data.Index.String()).
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

	if len(ids) == 0 {
		return nil, api_errors.ErrRecordsNotFound
	}

	return ids, nil
}

// SearchByID производит поиск по ID из базы. Возвращает ID из эластика подходящих записей
func (c *client) SearchByID(ctx context.Context, data elastic.Data) ([]string, error) {
	_, err := data.ValidateNote()
	if err != nil {
		return nil, err
	}

	query, err := data.SearchByIDQuery()
	if err != nil {
		return nil, fmt.Errorf("error while creating query for search note: %+v", err)
	}

	res, err := c.cl.Search().
		Index(data.Index.String()).
		Request(query).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("error searching note: %+v", err)
	}

	var ids []string

	for _, val := range res.Hits.Hits {
		ids = append(ids, *val.Id_)
	}

	if len(ids) == 0 {
		return nil, api_errors.ErrRecordsNotFound
	}

	return ids, nil
}
