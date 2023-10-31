package errors

import "errors"

// API ERRORS
var ErrUserNotFound = errors.New(`user not found`)
var ErrUnableCastVariable = errors.New(`unable to cast variable`)

// DB ERROR
var ErrNotesNotFound = errors.New(`notes not found`)
