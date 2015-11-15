package main

import (
	bce "baidubce"
	"fmt"
	"net/http"
	"os"
)

var credentials bce.Credentials = bce.Credentials{
	AccessKeyId:     os.Getenv("BAIDU_BCE_AK"),
	SecretAccessKey: os.Getenv("BAIDU_BCE_SK"),
}

var signOption bce.SignOption = bce.SignOption{
	Timestamp:                 "2015-11-16T08:13:49Z",
	ExpirationPeriodInSeconds: 1800,
}

var request bce.Request = bce.Request{
	HttpMethod:  "GET",
	URI:         "/baidubce-sdk-go",
	QueryString: "location",
	Header:      getHttpHeader(),
}

func getHttpHeader() http.Header {
	var header http.Header = http.Header{}

	header.Add("host", "bj.bcebos.com")
	//header.Add("Date", "Mon, 27 Apr 2015 16:23:49 +0800")
	header.Add("x-bce-date", "2015-11-16T08:13:49Z")

	return header
}

func generateSignature() {
	signature := bce.GenerateAuthorization(credentials, request, signOption)
	fmt.Println(signature)
}

func main() {
	generateSignature()
}
