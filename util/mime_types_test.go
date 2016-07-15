package util

import (
	"testing"

	"github.com/guoyao/baidubce-sdk-go/test"
)

func TestGuessMimeType(t *testing.T) {
	expected := "image/png"
	result := GuessMimeType("examples/test.png")

	if expected != result {
		t.Error(test.Format("GuessMimeType", result, expected))
	}

	expected = "application/octet-stream"
	result = GuessMimeType("examples/test")

	if expected != result {
		t.Error(test.Format("GuessMimeType", result, expected))
	}
}
