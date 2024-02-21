package errors

import "errors"

var ErrRemindersNotFound = errors.New(`reminders not found`)

var ErrInvalidDays = errors.New("must be in within the range from 1 to 180")
