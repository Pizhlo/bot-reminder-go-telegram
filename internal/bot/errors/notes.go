package errors

import "errors"

var ErrSecondDateBeforeFirst = errors.New("second date before first")
var ErrFirstDayFuture = errors.New("the first date has not yet arrived")
var ErrSecondDateFuture = errors.New("the second date has not yet arrived")
