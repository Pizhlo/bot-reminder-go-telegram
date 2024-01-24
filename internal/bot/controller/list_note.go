package controller

import (
	"context"

	tele "gopkg.in/telebot.v3"
)

func (c *Controller) ListNotes(ctx context.Context, telectx tele.Context) error {
	c.logger.Debugf("Controller: handling /notes command.\n")

	message, kb, err := c.noteSrv.GetAll(ctx, telectx.Chat().ID)
	if err != nil {
		c.logger.Errorf("Error while handling /notes command. User ID: %d. Error: %+v\n", telectx.Chat().ID, err)

		c.handleError(telectx, err)

		return err
	}

	c.logger.Debugf("Successfully got all user's notes. Sending message to user...\n")
	return telectx.Send(message, kb)
}

func (c *Controller) NextPageNotes(ctx context.Context, telectx tele.Context) error {
	c.logger.Debugf("Controller: handling next notes page command.\n")
	next, kb := c.noteSrv.NextPage(telectx.Chat().ID)

	return telectx.Edit(next, kb)
}

func (c *Controller) PrevPageNotes(ctx context.Context, telectx tele.Context) error {
	c.logger.Debugf("Controller: handling previous notes page command.\n")
	next, kb := c.noteSrv.PrevPage(telectx.Chat().ID)

	return telectx.Edit(next, kb)
}

func (c *Controller) LastPageNotes(ctx context.Context, telectx tele.Context) error {
	c.logger.Debugf("Controller: handling last notes page command.\n")
	next, kb := c.noteSrv.LastPage(telectx.Chat().ID)

	return telectx.Edit(next, kb)
}

func (c *Controller) FirstPageNotes(ctx context.Context, telectx tele.Context) error {
	c.logger.Debugf("Controller: handling first notes page command.\n")
	next, kb := c.noteSrv.FirstPage(telectx.Chat().ID)

	return telectx.Edit(next, kb)
}
