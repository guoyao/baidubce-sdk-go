package util

import (
	"testing"
)

const URI string = "http://bos.cn-n1.baidubce.com/v1/example/测试"

func TestGetUriPath(t *testing.T) {
	expected := "/v1/example/测试"
	path := GetUriPath(URI)
	if path != expected {
		t.Error(ToTestError("GetUriPath", path, expected))
	}
}

func TestUriEncodeExceptSlash(t *testing.T) {
	expected := "/v1/example/%E6%B5%8B%E8%AF%95"
	path := GetUriPath(URI)
	path = UriEncodeExceptSlash(path)
	if path != expected {
		t.Error(ToTestError("UriEncodeExceptSlash", path, expected))
	}
}
