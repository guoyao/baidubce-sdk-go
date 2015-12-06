package baidubce

import (
	"encoding/json"
)

// Error is a implementation of error.
type Error struct {
	Code, Message, RequestID string
	Raw                      error
}

func (err *Error) Error() string {
	if err.Message != "" {
		return err.Message
	}

	return err.Code
}

// NewErrorFromRaw returns a `Error` instance from another error instance.
func NewErrorFromRaw(err error) *Error {
	return &Error{
		Message: err.Error(),
		Raw:     err,
	}
}

// NewErrorFromJSON returns a `Error` instance from bytes.
func NewErrorFromJSON(bytes []byte) *Error {
	var err *Error
	rawError := json.Unmarshal(bytes, &err)
	if rawError != nil {
		return NewErrorFromRaw(rawError)
	}

	return err
}
