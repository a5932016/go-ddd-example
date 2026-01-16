package mError

import (
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestUnwrapError(t *testing.T) {
	customErr := customError{
		Code:     29001,
		HTTPCode: http.StatusNotFound,
		Message:  "Customer id not found.",
	}
	wrapErr := errors.Wrap(customErr, "wrap 1")
	err := Unwrap(wrapErr)
	assert.Equal(t, customErr, err)

	wrap2Err := errors.Wrap(wrapErr, "wrap 2")
	err2 := Unwrap(wrap2Err)
	assert.Equal(t, customErr, err2)
}

type customError struct {
	HTTPCode    int
	Code        int
	Message     string
	ErrorInfo   interface{}
	DeclineCode string
}

func (cErr customError) Error() string {
	return cErr.Message
}
