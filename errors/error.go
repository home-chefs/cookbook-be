package errors

import (
	"errors"
	"fmt"
)

type Error struct {
	Inner   error
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v: %v", e.Message, e.Inner)
}

func (e *Error) Unwrap() error {
	return e.Inner
}

func Wrap(err error, message string) error {
	return &Error{
		Inner:   err,
		Message: message,
	}
}

func New(message string) error {
	return errors.New(message)
}
