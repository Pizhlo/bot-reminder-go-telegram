package controller

func (c *Controller) saveUser(telegramID int64) error {
	id, err := c.srv.UserEditor.SaveUser(telegramID)
	if err != nil {
		return err
	}

	c.srv.UserCacheEditor.SaveUser(id, telegramID)

	return nil
}
