package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	bce "baidubce"
	"baidubce/bos"
)

var bosClient = bos.DefaultClient

func getBucketLocation() {
	option := &bce.SignOption{
		//Timestamp:                 "2015-11-20T10:00:05Z",
		ExpirationPeriodInSeconds: 1200,
		Headers: map[string]string{
			"host":                "bj.bcebos.com",
			"other":               "other",
			"x-bce-meta-data":     "meta data",
			"x-bce-meta-data-tag": "meta data tag",
			//"x-bce-date":          "2015-11-20T07:49:05Z",
			//"date": "2015-11-20T10:00:05Z",
		},
		HeadersToSign: []string{"host", "date", "other", "x-bce-meta-data", "x-bce-meta-data-tag"},
	}

	location, err := bosClient.GetBucketLocation("baidubce-sdk-go", option)

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(location.LocationConstraint)
}

func listBuckets() {
	bucketSummary, err := bosClient.ListBuckets(nil)

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(bucketSummary.Buckets)
}

func createBucket() {
	err := bosClient.CreateBucket("baidubce-sdk-go-create-bucket-example", nil)

	if err != nil {
		log.Println(err)
	}
}

func doesBucketExist() {
	// exists, err := bosClient.DoesBucketExist("baidubce-sdk-go-create-bucket-example", nil)
	exists, err := bosClient.DoesBucketExist("guoyao11122", nil)

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(exists)
}

func deleteBucket() {
	bucketName := "baidubce-sdk-go-delete-bucket-example"
	err := bosClient.CreateBucket(bucketName, nil)

	if err != nil {
		log.Println(err)
	} else {
		err := bosClient.DeleteBucket(bucketName, nil)
		if err != nil {
			log.Println(err)
		}
	}
}

func setBucketPrivate() {
	bucketName := "baidubce-sdk-go"
	err := bosClient.SetBucketPrivate(bucketName, nil)

	if err != nil {
		log.Println(err)
	}
}

func setBucketPublicRead() {
	bucketName := "baidubce-sdk-go"
	err := bosClient.SetBucketPublicRead(bucketName, nil)

	if err != nil {
		log.Println(err)
	}
}

func setBucketPublicReadWrite() {
	bucketName := "baidubce-sdk-go"
	err := bosClient.SetBucketPublicReadWrite(bucketName, nil)

	if err != nil {
		log.Println(err)
	}
}

func getBucketAcl() {
	bucketName := "baidubce-sdk-go"
	bucketAcl, err := bosClient.GetBucketAcl(bucketName, nil)

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(bucketAcl.Owner)

	for _, accessControl := range bucketAcl.AccessControlList {
		for _, grantee := range accessControl.Grantee {
			fmt.Println(grantee.Id)
		}
		for _, permission := range accessControl.Permission {
			fmt.Println(permission)
		}
	}
}

func setBucketAcl() {
	bucketName := "baidubce-sdk-go"
	bucketAcl := bos.BucketAcl{
		AccessControlList: []bos.Grant{
			bos.Grant{
				Grantee: []bos.BucketGrantee{
					bos.BucketGrantee{Id: "ef5a4b19192f4931adcf0e12f82795e2"},
				},
				Permission: []string{"FULL_CONTROL"},
			},
		},
	}

	err := bosClient.SetBucketAcl(bucketName, bucketAcl, nil)

	if err != nil {
		log.Println(err)
	}
}

func putObject() {
	bucketName := "baidubce-sdk-go"

	objectKey := "put-object-from-string.txt"
	str := "Hello World 你好"

	option := new(bce.SignOption)
	metadata := new(bos.ObjectMetadata)
	metadata.AddUserMetadata("x-bce-meta-name", "guoyao")

	putObjectResponse, bceError := bosClient.PutObject(bucketName, objectKey, str, metadata, option)

	if bceError != nil {
		log.Println(bceError)
	} else {
		fmt.Println(putObjectResponse.GetETag())
	}

	pwd, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	filePath := path.Join(pwd, "baidubce", "examples", "baidubce-sdk-go-test.pdf")

	objectKey = "pdf/put-object-from-bytes.pdf"
	byteArray, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Println(err)
	} else {
		putObjectResponse, bceError = bosClient.PutObject(bucketName, objectKey, byteArray, nil, nil)

		if bceError != nil {
			log.Println(bceError)
		} else {
			fmt.Println(putObjectResponse.GetETag())
		}
	}

	objectKey = "pdf/put-object-from-file.pdf"
	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		log.Println(err)
	} else {
		putObjectResponse, bceError = bosClient.PutObject(bucketName, objectKey, file, nil, nil)

		if bceError != nil {
			log.Println(bceError)
		} else {
			fmt.Println(putObjectResponse.GetETag())
		}
	}
}

func deleteObject() {
	bucketName := "baidubce-sdk-go"

	objectKey := "put-object-from-string.txt"
	str := "Hello World 你好"

	option := new(bce.SignOption)
	metadata := new(bos.ObjectMetadata)
	metadata.AddUserMetadata("x-bce-meta-name", "guoyao")

	putObjectResponse, bceError := bosClient.PutObject(bucketName, objectKey, str, metadata, option)

	if bceError != nil {
		log.Println(bceError)
	} else {
		fmt.Println(putObjectResponse.GetETag())
	}

	bceError = bosClient.DeleteObject(bucketName, objectKey, nil)

	if bceError != nil {
		log.Println(bceError)
	}
}

func listObjects() {
	bucketName := "baidubce-sdk-go"
	params := map[string]string{
		"prefix":    "pdf/",
		"delimiter": "/",
		"marker":    "",
		//"marker":    "pdf/put-object-from-bytes.pdf",
		"maxKeys": "1000",
	}

	listObjectsResponse, bceError := bosClient.ListObjects(bucketName, params, nil)

	if bceError != nil {
		log.Println(bceError)
	} else {
		for _, objectSummary := range listObjectsResponse.Contents {
			fmt.Println(objectSummary.Key)
		}

		for _, prefix := range listObjectsResponse.GetCommonPrefixes() {
			fmt.Println(prefix)
		}
	}
}

func copyObject() {
	srcBucketName := "baidubce-sdk-go"
	srcKey := "baidubce-sdk-go-test.pdf"
	destBucketName := "baidubce-sdk-go"
	destKey := "pdf/baidubce-sdk-go-test-copy.pdf"

	copyObjectResponse, bceError := bosClient.CopyObject(srcBucketName, srcKey, destBucketName, destKey, nil)

	if bceError != nil {
		log.Println(bceError)
	} else {
		fmt.Println(copyObjectResponse.ETag)
		fmt.Println(copyObjectResponse.LastModified)
	}
}

func main() {
	copyObject()
	return
	listObjects()
	deleteObject()
	putObject()
	getBucketAcl()
	setBucketAcl()
	getBucketLocation()
	listBuckets()
	createBucket()
	doesBucketExist()
	deleteBucket()
	setBucketPublicReadWrite()
	setBucketPublicRead()
	setBucketPrivate()
}
