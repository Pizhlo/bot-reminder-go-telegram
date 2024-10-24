package elastic

import (
	"context"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/elastic"
)

func (c *client) Delete(ctx context.Context, data elastic.Data) error {
	note, err := data.ValidateNote()
	if err != nil {
		return err
	}

	res, err := c.cl.Delete(data.Index, note.ElasticID).Do(ctx)
	if err != nil {
		return err
	}

	if res.Result.Name != "deleted" || res.Result.String() != "deleted" {
		return fmt.Errorf("result of deletion is not equal `deleted`: %+v", res.Result.Name)
	}

	return nil
}
