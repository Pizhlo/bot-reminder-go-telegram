package elastic

import (
	"context"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/elastic"
)

func (c *client) Update(ctx context.Context, data elastic.Data) error {
	_, err := data.ValidateNote()
	if err != nil {
		return fmt.Errorf("error while validating note while updating: %+v", err)
	}

	elasticID, err := c.getElasticID(ctx, data)
	if err != nil {
		return fmt.Errorf("error while getting elasctid ID: %+v", err)
	}

	data.SetElasticID(elasticID)

	req, err := data.UpdateQuery()
	if err != nil {
		return fmt.Errorf("error creating request while updating: %+v", err)
	}

	_, err = c.cl.Update(data.Index.String(), elasticID).Request(req).Do(ctx)
	return err
}
