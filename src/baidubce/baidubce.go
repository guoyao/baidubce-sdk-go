package main

import (
	//"baidubce/util"
	"fmt"
	//"net/url"
	"baidubce/core"
	"net/http"
)

var credentials core.Credentials = core.Credentials{
	AccessKeyId:     "0b0f67dfb88244b289b72b142befad0c",
	SecretAccessKey: "bad522c2126a4618a8125f4b6cf6356f",
}

var signOption core.SignOption = core.SignOption{
	Timestamp:                 "2015-04-27T08:23:49Z",
	ExpirationPeriodInSeconds: 1800,
}

var request core.Request = core.Request{
	HttpMethod:  "PUT",
	URI:         "/v1/test/myfolder/readme.txt",
	QueryString: "partNumber=9&uploadId=VXBsb2FkIElpZS5tMnRzIHVwbG9hZA",
	Header:      getHttpHeader(),
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

func main() {
	signature := core.Sign(credentials, request, signOption)
	fmt.Println(signature)
}
