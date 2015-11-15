package baidubce

import (
	"baidubce/util"
	"net/http"
	"strings"
	"testing"
)

var credentials Credentials = Credentials{
	AccessKeyId:     "0b0f67dfb88244b289b72b142befad0c",
	SecretAccessKey: "bad522c2126a4618a8125f4b6cf6356f",
}

var signOption *SignOption = NewSignOption("2015-04-27T08:23:49Z", EXPIRATION_PERIOD_IN_SECONDS)

var request Request = Request{
	HttpMethod: "PUT",
	URI:        "/v1/test/myfolder/readme.txt",
	Params: map[string]string{
		"partNumber": "9",
		"uploadId":   "VXBsb2FkIElpZS5tMnRzIHVwbG9hZA",
	},
	Header: getHttpHeader(),
}

func TestGetSigningKey(t *testing.T) {
	const expected = "d9f35aaba8a5f3efa654851917114b6f22cd831116fd7d8431e08af22dcff24c"
	signingKey := getSigningKey(credentials, signOption)

	if signingKey != expected {
		t.Error(util.ToTestError("getSigningKey", signingKey, expected))
	}
}

func TestGetCanonicalQueryString(t *testing.T) {
	const expected = "text10=test&text1=%E6%B5%8B%E8%AF%95&text="
	params := map[string]string{
		"text":   "",
		"text1":  "测试",
		"text10": "test",
	}
	encodedQueryString := getCanonicalQueryString(params)

	if encodedQueryString != expected {
		t.Error(util.ToTestError("getCanonicalQueryString", encodedQueryString, expected))
	}
}

func TestGetCanonicalHeader(t *testing.T) {
	expected := strings.Join([]string{
		"content-length:8",
		"content-md5:0a52730597fb4ffa01fc117d9e71e3a9",
		"content-type:text%2Fplain",
		"host:bj.bcebos.com",
		"x-bce-date:2015-04-27T08%3A23%3A49Z",
	}, "\n")

	canonicalHeader := getCanonicalHeader(getHttpHeader())

	if canonicalHeader != expected {
		t.Error(util.ToTestError("getCanonicalHeaders", canonicalHeader, expected))
	}
}

func TestSign(t *testing.T) {
	expected := "a19e6386e990691aca1114a20357c83713f1cb4be3d74942bb4ed37469ecdacf"
	signature := sign(credentials, request, signOption)

	if signature != expected {
		t.Error(util.ToTestError("sign", signature, expected))
	}
}

func getHttpHeader() http.Header {
	var header http.Header = http.Header{}

	header.Add("host", "bj.bcebos.com")
	header.Add("Date", "Mon, 27 Apr 2015 16:23:49 +0800")
	header.Add("Content-Type", "text/plain")
	header.Add("Content-Length", "8")
	header.Add("Content-Md5", "0a52730597fb4ffa01fc117d9e71e3a9")
	header.Add("x-bce-date", "2015-04-27T08:23:49Z")

	return header
}
