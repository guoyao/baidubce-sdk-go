package main

import (
	bce "baidubce"
	"fmt"
)

var credentials = bce.DefaultCredentials

var signOption *bce.SignOption = &bce.SignOption{
	Timestamp:                 "2015-11-16T08:13:49Z",
	ExpirationPeriodInSeconds: 1800,
}

func getRequest() *bce.Request {
	req, _ := bce.NewRequest(
		"GET",
		"/baidubce-sdk-go",
		"",
		map[string]string{"location": ""},
		nil,
	)

	return req
}

func generateAuthorization() {
	req := getRequest()
	authorization := bce.GenerateAuthorization(credentials, *req, signOption)
	fmt.Println(authorization)
}

func main() {
	generateAuthorization()
}
