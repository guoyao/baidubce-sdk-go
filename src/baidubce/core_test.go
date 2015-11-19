package baidubce

import (
	"baidubce/util"
	"testing"
)

var credentials Credentials = Credentials{
	AccessKeyId:     "0b0f67dfb88244b289b72b142befad0c",
	SecretAccessKey: "bad522c2126a4618a8125f4b6cf6356f",
}

var defaultSignOption *SignOption = NewSignOption(
	"2015-04-27T08:23:49Z",
	EXPIRATION_PERIOD_IN_SECONDS,
	getHeaders(),
	nil,
)

func TestGetSigningKey(t *testing.T) {
	const expected = "d9f35aaba8a5f3efa654851917114b6f22cd831116fd7d8431e08af22dcff24c"
	signingKey := getSigningKey(credentials, defaultSignOption)

	if signingKey != expected {
		t.Error(util.ToTestError("getSigningKey", signingKey, expected))
	}
}

func TestSign(t *testing.T) {
	expected := "a19e6386e990691aca1114a20357c83713f1cb4be3d74942bb4ed37469ecdacf"
	req := getRequest()
	signature := sign(credentials, *req, defaultSignOption)

	if signature != expected {
		t.Error(util.ToTestError("sign", signature, expected))
	}
}

func TestGenerateAuthorization(t *testing.T) {
	expected := "bce-auth-v1/0b0f67dfb88244b289b72b142befad0c/2015-04-27T08:23:49Z/1800//a19e6386e990691aca1114a20357c83713f1cb4be3d74942bb4ed37469ecdacf"
	req := getRequest()
	authorization := GenerateAuthorization(credentials, *req, defaultSignOption)
	if authorization != expected {
		t.Error(util.ToTestError("GenerateAuthorization", authorization, expected))
	}
}

func getRequest() *Request {
	request, _ := NewRequest(
		"PUT",
		"/v1/test/myfolder/readme.txt",
		"",
		map[string]string{
			"partNumber": "9",
			"uploadId":   "VXBsb2FkIElpZS5tMnRzIHVwbG9hZA",
		},
		nil,
	)

	return request
}

func getHeaders() map[string]string {
	headers := map[string]string{
		"Host":           Region["bj"],
		"Date":           "Mon, 27 Apr 2015 16:23:49 +0800",
		"Content-Type":   "text/plain",
		"Content-Length": "8",
		"Content-Md5":    "0a52730597fb4ffa01fc117d9e71e3a9",
		"x-bce-date":     "2015-04-27T08:23:49Z",
	}

	return headers
}
