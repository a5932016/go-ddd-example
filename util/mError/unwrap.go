package mError

import "errors"

func Unwrap(err error) error {
	currentErr := err
	for errors.Unwrap(currentErr) != nil {
		currentErr = errors.Unwrap(currentErr)
	}
	return currentErr
}
