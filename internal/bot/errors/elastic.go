package errors

import "errors"

type ElasticError struct {
	message string
}

func (e ElasticError) Error() string {
	return e.message
}

var ErrRecordsNotFound = errors.New(`records not found in elastic`)

// func NewNotFound(format string, a ...any) error {
// 	// e := ErrRecordsNotFound
// 	// e.message = fmt.Sprintf(format, a...)
// 	// return e
// 	return nil
// }
