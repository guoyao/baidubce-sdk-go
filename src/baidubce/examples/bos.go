package main

import (
	bce "baidubce"
	"baidubce/bos"
	"log"
	"os"
)

var credentials bce.Credentials = bce.Credentials{
	AccessKeyId:     os.Getenv("BAIDU_BCE_AK"),
	SecretAccessKey: os.Getenv("BAIDU_BCE_SK"),
}

var bosClient bos.Client = bos.NewClient(
	bce.Config{
		Credentials: credentials,
		Endpoint:    "baidubce-sdk-go.tocloud.org",
	},
)

func GetBucketLocation() {
	option := &bce.SignOption{
		//Timestamp:                 "2015-11-19T15:25:05Z",
		ExpirationPeriodInSeconds: 1000,
		Headers: map[string]string{
			"host":                "bj.bcebos.com",
			"other":               "other",
			"x-bce-meta-data":     "meta data",
			"x-bce-meta-data-tag": "meta data tag",
		},
		HeadersToSign: []string{"host", "date", "other", "x-bce-meta-data", "x-bce-meta-data-tag"},
	}

	body, err := bosClient.GetBucketLocation("baidubce-sdk-go", option)

	if err != nil {
		log.Println(err)
	}

	log.Println(body)
}

func ListBucket() {
	body, err := bosClient.ListBucket(nil)

	if err != nil {
		log.Println(err)
	}

	log.Println(body)
}

func main() {
	GetBucketLocation()
	ListBucket()
}
