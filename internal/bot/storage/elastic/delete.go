package elastic

import (
	"context"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/elastic"
)

func (c *client) Delete(ctx context.Context, data elastic.Data) error {
	note, err := data.ValidateNote()
	if err != nil {
		return fmt.Errorf("error validating note while deleting: %+v", err)
	}

	res, err := c.cl.Delete(data.Index.String(), note.ElasticID).Do(ctx)
	if err != nil {
		return fmt.Errorf("error while deleting notes by ID: %+v", err)
	}

	if res.Result.Name != "deleted" || res.Result.String() != "deleted" {
		return fmt.Errorf("result of deletion is not equal `deleted`: %+v", res.Result.Name)
	}

	return nil
}

func (c *client) DeleteAllByUserID(ctx context.Context, data elastic.Data) error {
	req, err := data.DeleteByQuery()
	if err != nil {
		return fmt.Errorf("error creating query for deleting by query: %+v", err)
	}

	_, err = c.cl.DeleteByQuery(data.Index.String()).Request(req).Do(ctx)
	if err != nil {
		return fmt.Errorf("error while deleting notes by query: %+v", err)
	}

	return nil
}
