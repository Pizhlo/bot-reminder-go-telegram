package errors

import "fmt"

type ElasticError struct {
	message string
}

func (e *ElasticError) Error() string {
	return e.message
}

var ErrRecordsNotFound ElasticError

func NewNotFound(format string, a ...any) *ElasticError {
	e := ErrRecordsNotFound
	e.message = fmt.Sprintf(format, a...)
	return &e
}
