package baidubce

import (
	"encoding/json"
)

type Error struct {
	Code, Message, RequestId string
	Raw                      error
}

func (err *Error) Error() string {
	if err.Message != "" {
		return err.Message
	}

	return err.Code
}

func NewErrorFromRaw(err error) *Error {
	return &Error{
		Message: err.Error(),
		Raw:     err,
	}
}

func NewErrorFromJson(bytes []byte) *Error {
	var err *Error
	rawError := json.Unmarshal(bytes, &err)
	if rawError != nil {
		return NewErrorFromRaw(rawError)
	}

	return err
}
