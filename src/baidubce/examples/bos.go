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

var bosClient bos.Client = bos.Client{
	bos.Config{
		Credentials: credentials,
	},
}

func GetBucketLocation() {
	body, err := bosClient.GetBucketLocation("baidubce-sdk-go", nil)

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
