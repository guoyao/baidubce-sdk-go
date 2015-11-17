package util

import (
	"testing"
	"time"
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

func TestTimeToUTCString(t *testing.T) {
	expected := "2015-11-16T07:33:15Z"
	datetime, _ := time.Parse(time.RFC1123, "Mon, 16 Nov 2015 15:33:15 CST")
	utc := TimeToUTCString(datetime)
	if utc != expected {
		t.Error(ToTestError("TimeToUTCString", utc, expected))
	}
}

func TestHostToUrl(t *testing.T) {
	expected := "http://bj.bcebos.com"
	host := "bj.bcebos.com"
	url := HostToUrl(host)
	if url != expected {
		t.Error(ToTestError("HostToUrl", url, expected))
	}

	host = "http://bj.bcebos.com"
	url = HostToUrl(host)
	if url != expected {
		t.Error(ToTestError("HostToUrl", url, expected))
	}
}
