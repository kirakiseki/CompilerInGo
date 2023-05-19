package utils

import (
	"errors"
	"fmt"
)

func NewError(msg string) error {
	return errors.New(msg)
}

func NewErrorf(format string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(format, args...))
}
