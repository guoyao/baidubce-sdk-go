package util

import (
	"baidubce/test"
	"strings"
	"testing"
	"time"
)

const URI string = "http://bos.cn-n1.baidubce.com/v1/example/测试"

func TestGetUriPath(t *testing.T) {
	expected := "/v1/example/测试"
	path := GetUriPath(URI)
	if path != expected {
		t.Error(test.Format("GetUriPath", path, expected))
	}
}

func TestUriEncodeExceptSlash(t *testing.T) {
	expected := "/v1/example/%E6%B5%8B%E8%AF%95"
	path := GetUriPath(URI)
	path = UriEncodeExceptSlash(path)
	if path != expected {
		t.Error(test.Format("UriEncodeExceptSlash", path, expected))
	}
}

func TestTimeToUTCString(t *testing.T) {
	expected := "2015-11-16T07:33:15Z"
	datetime, _ := time.Parse(time.RFC1123, "Mon, 16 Nov 2015 15:33:15 CST")
	utc := TimeToUTCString(datetime)
	if utc != expected {
		t.Error(test.Format("TimeToUTCString", utc, expected))
	}
}

func TestHostToUrl(t *testing.T) {
	expected := "http://bj.bcebos.com"
	host := "bj.bcebos.com"
	url := HostToUrl(host)
	if url != expected {
		t.Error(test.Format("HostToUrl", url, expected))
	}

	host = "http://bj.bcebos.com"
	url = HostToUrl(host)
	if url != expected {
		t.Error(test.Format("HostToUrl", url, expected))
	}
}

func TestToCanonicalQueryString(t *testing.T) {
	const expected = "text10=test&text1=%E6%B5%8B%E8%AF%95&text="
	params := map[string]string{
		"text":   "",
		"text1":  "测试",
		"text10": "test",
	}
	encodedQueryString := ToCanonicalQueryString(params)

	if encodedQueryString != expected {
		t.Error(test.Format("ToCanonicalQueryString", encodedQueryString, expected))
	}
}

func TestToCanonicalHeaderString(t *testing.T) {
	expected := strings.Join([]string{
		"content-length:8",
		"content-md5:0a52730597fb4ffa01fc117d9e71e3a9",
		"content-type:text%2Fplain",
		"host:bj.bcebos.com",
		"x-bce-date:2015-04-27T08%3A23%3A49Z",
	}, "\n")

	canonicalHeader := ToCanonicalHeaderString(getHeaders())

	if canonicalHeader != expected {
		t.Error(test.Format("ToCanonicalHeaderString", canonicalHeader, expected))
	}
}

func getHeaders() map[string]string {
	header := map[string]string{
		"Host":           "bj.bcebos.com",
		"Content-Type":   "text/plain",
		"Content-Length": "8",
		"Content-Md5":    "0a52730597fb4ffa01fc117d9e71e3a9",
		"x-bce-date":     "2015-04-27T08:23:49Z",
	}

	return header
}
