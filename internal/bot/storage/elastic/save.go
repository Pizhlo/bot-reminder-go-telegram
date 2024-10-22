package elastic

import (
	"context"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/elastic"
	"github.com/sirupsen/logrus"
)

func (c *client) Save(ctx context.Context, data elastic.Data) error {
	_, err := data.ValidateNote()
	if err != nil {
		return fmt.Errorf("error validating note while saving: %+v", err)
	}

	_, err = c.cl.Index(data.Index).
		Request(data.Model).
		Do(ctx)
	if err != nil {
		return err
	}

	logrus.Debugf("Elastic: sucecssfully saved user's note")

	return nil
}
