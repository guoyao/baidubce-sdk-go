package bce

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Error is a implementation of error.
type Error struct {
	StatusCode               int
	Code, Message, RequestID string
	Raw                      error
}

func (err *Error) Error() string {
	return fmt.Sprintf("Error Message: \"%s\", Error Code: \"%s\", Status Code: %d, Request Id: \"%s\"", err.Message, err.Code, err.StatusCode, err.RequestID)
}

func NewError(err error) *Error {
	if bceError, ok := err.(*Error); ok {
		return bceError
	}

	return newErrorFromRaw(err)
}

// newErrorFromRaw returns a `Error` instance from another error instance.
func newErrorFromRaw(err error) *Error {
	return &Error{
		Message: err.Error(),
		Raw:     err,
	}
}

// NewErrorFromJSON returns a `Error` instance from bytes.
func NewErrorFromJSON(bytes []byte) *Error {
	var err *Error

	if bytes == nil || string(bytes) == "" {
		return newErrorFromRaw(errors.New(""))
	}

	rawError := json.Unmarshal(bytes, &err)

	if rawError != nil {
		return newErrorFromRaw(errors.New(string(bytes)))
	}

	return err
}
