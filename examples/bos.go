package examples

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"sync"

	"github.com/guoyao/baidubce-sdk-go/bce"
	"github.com/guoyao/baidubce-sdk-go/bos"
	//"github.com/guoyao/baidubce-sdk-go/util"
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

	filePath := path.Join(pwd, "baidubce", "examples", "test.tgz")

	objectKey = "compressed/put-object-from-bytes.tgz"
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

	objectKey = "compressed/put-object-from-file.tgz"
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

func deleteMultipleObjects() {
	bucketName := "baidubce-sdk-go"

	objects := []string{
		"examples/delete-multiple-objects/put-object-from-string.txt",
		"examples/delete-multiple-objects/put-object-from-string-2.txt",
		"examples/delete-multiple-objects/put-object-from-string-3.txt",
	}
	str := "Hello World 你好"

	for _, value := range objects {
		putObjectResponse, bceError := bosClient.PutObject(bucketName, value, str, nil, nil)

		if bceError != nil {
			log.Fatal(bceError)
		}
		//util.CheckError(bceError)

		fmt.Println(putObjectResponse.GetETag())
	}

	deleteMultipleObjectsResponse, bceError := bosClient.DeleteMultipleObjects(bucketName, objects, nil)

	if bceError != nil {
		log.Println(bceError)
	} else if deleteMultipleObjectsResponse != nil {
		for _, deleteMultipleObjectsError := range deleteMultipleObjectsResponse.Errors {
			log.Println(deleteMultipleObjectsError.Error())
		}
	}
}

func listObjects() {
	bucketName := "baidubce-sdk-go"

	listObjectsResponse, bceError := bosClient.ListObjects(bucketName, nil)

	if bceError != nil {
		log.Println(bceError)
	} else {
		for _, objectSummary := range listObjectsResponse.Contents {
			fmt.Println(objectSummary.Key, objectSummary.ETag)
		}

		for _, prefix := range listObjectsResponse.GetCommonPrefixes() {
			fmt.Println(prefix)
		}
	}
}

func listObjectsFromRequest() {
	listObjectsRequest := bos.ListObjectsRequest{
		BucketName: "baidubce-sdk-go",
		Delimiter:  "/",
		//Marker:    "compressed/put-object-from-bytes.tgz",
		//Prefix:    "compressed/",
		MaxKeys: 100,
	}

	listObjectsResponse, bceError := bosClient.ListObjectsFromRequest(listObjectsRequest, nil)

	if bceError != nil {
		log.Println(bceError)
	} else {
		for _, objectSummary := range listObjectsResponse.Contents {
			fmt.Println(objectSummary.Key, objectSummary.ETag)
		}

		for _, prefix := range listObjectsResponse.GetCommonPrefixes() {
			fmt.Println(prefix)
		}
	}
}

func copyObject() {
	srcBucketName := "baidubce-sdk-go"
	srcKey := "test.tgz"
	destBucketName := "baidubce-sdk-go"
	destKey := "compressed/test-copy.tgz"

	copyObjectResponse, bceError := bosClient.CopyObject(srcBucketName, srcKey, destBucketName, destKey, nil)

	if bceError != nil {
		log.Println(bceError)
	} else {
		fmt.Println(copyObjectResponse.ETag, copyObjectResponse.LastModified)
	}
}

func copyObjectFromRequest() {
	etag := "fa412a6ca6d415208be69bc4a00f4103"

	copyObjectRequest := bos.CopyObjectRequest{
		SrcBucketName:  "baidubce-sdk-go",
		SrcKey:         "test.tgz",
		DestBucketName: "baidubce-sdk-go",
		DestKey:        "compressed/test-copy.tgz",
		//SourceMatch:    etag,
		SourceNoneMatch: etag,
		//SourceModifiedSince:   "2016-05-28T22:32:00Z",
		//SourceUnmodifiedSince: "2016-05-28T22:32:00Z",
		ObjectMetadata: &bos.ObjectMetadata{
			CacheControl: "no-cache",
			UserMetadata: map[string]string{
				"test-user-metadata": "test user metadata",
				"x-bce-meta-name":    "x-bce-meta-name",
			},
		},
	}

	copyObjectResponse, bceError := bosClient.CopyObjectFromRequest(copyObjectRequest, nil)

	if bceError != nil {
		log.Println(bceError)
	} else {
		fmt.Println(copyObjectResponse.ETag, copyObjectResponse.LastModified)
	}
}

func getObject() {
	bucketName := "baidubce-sdk-go"
	objectKey := "test.tgz"

	object, bceError := bosClient.GetObject(bucketName, objectKey, nil)

	if bceError != nil {
		log.Println(bceError)
	} else {
		fmt.Println(object.ObjectMetadata)

		byteArray, err := ioutil.ReadAll(object.ObjectContent)

		if err != nil {
			log.Println(err)
		} else {
			err = ioutil.WriteFile(objectKey, byteArray, 0666)

			if err != nil {
				log.Println(err)
			}
		}
	}
}

func getObjectFromRequest() {
	bucketName := "baidubce-sdk-go"
	objectKey := "test.tgz"

	getObjectRequest := bos.GetObjectRequest{
		BucketName: bucketName,
		ObjectKey:  objectKey,
	}
	getObjectRequest.SetRange(0, 1000)

	object, bceError := bosClient.GetObjectFromRequest(getObjectRequest, nil)

	if bceError != nil {
		log.Println(bceError)
	} else {
		fmt.Println(object.ObjectMetadata)

		byteArray, err := ioutil.ReadAll(object.ObjectContent)

		if err != nil {
			log.Println(err)
		} else {
			err = ioutil.WriteFile(objectKey, byteArray, 0666)

			if err != nil {
				log.Println(err)
			}
		}
	}
}

func getObjectToFile() {
	bucketName := "baidubce-sdk-go"
	objectKey := "test.tgz"

	getObjectRequest := &bos.GetObjectRequest{
		BucketName: bucketName,
		ObjectKey:  objectKey,
	}
	getObjectRequest.SetRange(0, 1000)

	file, err := os.OpenFile(objectKey, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		log.Println(err)
	} else {
		objectMetadata, bceError := bosClient.GetObjectToFile(getObjectRequest, file, nil)

		if bceError != nil {
			log.Println(bceError)
		} else {
			fmt.Println(objectMetadata)
		}
	}
}

func getObjectMetadata() {
	bucketName := "baidubce-sdk-go"
	objectKey := "test.tgz"

	objectMetadata, bceError := bosClient.GetObjectMetadata(bucketName, objectKey, nil)

	if bceError != nil {
		log.Println(bceError)
	} else {
		fmt.Println(objectMetadata)
	}
}

func generatePresignedUrl() {
	bucketName := "baidubce-sdk-go"
	objectKey := "test.tgz"

	option := &bce.SignOption{
		ExpirationPeriodInSeconds: 300,
	}

	url, bceError := bosClient.GeneratePresignedUrl(bucketName, objectKey, option)

	if bceError != nil {
		log.Println(bceError)
	} else {
		fmt.Println(url)
	}
}

func appendObject() {
	bucketName := "baidubce-sdk-go"

	objectKey := "append-object-from-string.txt"
	str := "Hello World 你好"
	offset := 0

	option := new(bce.SignOption)
	metadata := new(bos.ObjectMetadata)
	metadata.AddUserMetadata("x-bce-meta-name", "guoyao")

	appendObjectResponse, bceError := bosClient.AppendObject(bucketName, objectKey, offset, str, metadata, option)

	if bceError != nil {
		log.Println(bceError)
	} else {
		fmt.Println(appendObjectResponse.GetETag(), appendObjectResponse.GetMD5(),
			appendObjectResponse.GetNextAppendOffset())
	}
}

func multipartUpload() {
	bucketName := "baidubce-sdk-go"
	objectKey := "test-multipart-upload.zip"

	initiateMultipartUploadRequest := bos.InitiateMultipartUploadRequest{
		BucketName: bucketName,
		ObjectKey:  objectKey,
	}

	initiateMultipartUploadResponse, bceError := bosClient.InitiateMultipartUpload(initiateMultipartUploadRequest, nil)

	if bceError != nil {
		panic(bceError)
	}

	uploadId := initiateMultipartUploadResponse.UploadId

	pwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	filePath := path.Join(pwd, "baidubce", "examples", objectKey)
	file, err := os.Open(filePath)

	defer file.Close()

	if err != nil {
		panic(err)
	}

	fileInfo, err := file.Stat()

	if err != nil {
		panic(err)
	}

	var partSize int64 = 1024 * 1024 * 5
	var totalSize int64 = fileInfo.Size()
	var partCount int = int(math.Ceil(float64(totalSize) / float64(partSize)))

	var waitGroup sync.WaitGroup
	partETags := make([]bos.PartETag, 0, partCount)

	for i := 0; i < partCount; i++ {
		var skipBytes int64 = partSize * int64(i)
		var size int64 = int64(math.Min(float64(totalSize-skipBytes), float64(partSize)))

		byteArray := make([]byte, size, size)
		_, err := file.Read(byteArray)

		if err != nil {
			panic(err)
		}

		partNumber := i + 1

		uploadPartRequest := bos.UploadPartRequest{
			BucketName: bucketName,
			ObjectKey:  objectKey,
			UploadId:   uploadId,
			PartSize:   size,
			PartNumber: partNumber,
			PartData:   byteArray,
		}

		waitGroup.Add(1)

		partETags = append(partETags, bos.PartETag{PartNumber: partNumber})

		go func(partNumber int) {
			defer waitGroup.Done()

			uploadPartResponse, err := bosClient.UploadPart(uploadPartRequest, nil)

			if err != nil {
				panic(err)
			}

			partETags[partNumber-1].ETag = uploadPartResponse.GetETag()
		}(partNumber)
	}

	waitGroup.Wait()
	waitGroup.Add(1)

	go func() {
		defer waitGroup.Done()

		completeMultipartUploadRequest := bos.CompleteMultipartUploadRequest{
			BucketName: bucketName,
			ObjectKey:  objectKey,
			UploadId:   uploadId,
			Parts:      partETags,
		}

		completeMultipartUploadResponse, err := bosClient.CompleteMultipartUpload(
			completeMultipartUploadRequest, nil)

		if err != nil {
			panic(err)
		}

		fmt.Println(completeMultipartUploadResponse.ETag)
	}()

	waitGroup.Wait()
}

func multipartUploadFromFile() {
	bucketName := "baidubce-sdk-go"
	objectKey := "test-multipart-upload-from-file.zip"

	pwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	filePath := path.Join(pwd, "baidubce", "examples", "test-multipart-upload.zip")
	var partSize int64 = 1024 * 1024 * 2

	completeMultipartUploadResponse, bceError := bosClient.MultipartUploadFromFile(bucketName,
		objectKey, filePath, partSize)

	if bceError != nil {
		log.Println(bceError)
	} else {
		fmt.Println(completeMultipartUploadResponse.ETag)
	}
}

func abortMultipartUpload() {
	bucketName := "baidubce-sdk-go"
	objectKey := "test-multipart-upload.zip"

	initiateMultipartUploadRequest := bos.InitiateMultipartUploadRequest{
		BucketName: bucketName,
		ObjectKey:  objectKey,
	}

	initiateMultipartUploadResponse, bceError := bosClient.InitiateMultipartUpload(initiateMultipartUploadRequest, nil)

	if bceError != nil {
		panic(bceError)
	}

	uploadId := initiateMultipartUploadResponse.UploadId

	abortMultipartUploadRequest := bos.AbortMultipartUploadRequest{
		BucketName: bucketName,
		ObjectKey:  objectKey,
		UploadId:   uploadId,
	}

	bceError = bosClient.AbortMultipartUpload(abortMultipartUploadRequest, nil)

	if bceError != nil {
		log.Println(bceError)
	}
}

func listParts() {
	bucketName := "baidubce-sdk-go"
	objectKey := "test-multipart-upload.zip"

	listPartsResponse, err := bosClient.ListParts(bucketName, objectKey, "4b17efee1a6abfcdab1c856afdc070c2", nil)

	if err != nil {
		log.Println(err)
		return
	}

	for _, partSummary := range listPartsResponse.Parts {
		fmt.Println(partSummary.PartNumber, partSummary.ETag, partSummary.Size, partSummary.LastModified)
	}
}

func listPartsFromRequest() {
	bucketName := "baidubce-sdk-go"
	objectKey := "test-multipart-upload.zip"

	listPartsRequest := bos.ListPartsRequest{
		BucketName: bucketName,
		ObjectKey:  objectKey,
		UploadId:   "4b17efee1a6abfcdab1c856afdc070c2",
		//PartNumberMarker: "1",
		MaxParts: 1,
	}

	listPartsResponse, err := bosClient.ListPartsFromRequest(listPartsRequest, nil)

	if err != nil {
		log.Println(err)
		return
	}

	for _, partSummary := range listPartsResponse.Parts {
		fmt.Println(partSummary.PartNumber, partSummary.ETag, partSummary.Size, partSummary.LastModified)
	}
}

func listMultipartUploads() {
	bucketName := "baidubce-sdk-go"
	listMultipartUploadsResponse, err := bosClient.ListMultipartUploads(bucketName, nil)

	if err != nil {
		log.Println(err)
		return
	}

	for _, multipartUploadSummary := range listMultipartUploadsResponse.Uploads {
		fmt.Println(multipartUploadSummary.Key, multipartUploadSummary.UploadId, multipartUploadSummary.Initiated)
	}

	for _, prefix := range listMultipartUploadsResponse.GetCommonPrefixes() {
		fmt.Println(prefix)
	}
}

func listMultipartUploadsFromRequest() {
	/*
		bucketName := "baidubce-sdk-go"
		objectKey := "compressed/test-multipart-upload.zip"

		initiateMultipartUploadRequest := bos.InitiateMultipartUploadRequest{
			BucketName: bucketName,
			ObjectKey:  objectKey,
		}

		initiateMultipartUploadResponse, bceError := bosClient.InitiateMultipartUpload(initiateMultipartUploadRequest, nil)

		if bceError != nil {
			log.Println(bceError)
			log.Println(initiateMultipartUploadResponse.UploadId)
			return
		}
	*/

	listMultipartUploadsRequest := bos.ListMultipartUploadsRequest{
		BucketName: "baidubce-sdk-go",
		//Delimiter:  "/",
		//KeyMarker:  "compressed/test-multipart-upload.zip",
		//Prefix:     "compressed/",
		MaxUploads: 100,
	}

	listMultipartUploadsResponse, err := bosClient.ListMultipartUploadsFromRequest(listMultipartUploadsRequest, nil)

	if err != nil {
		log.Println(err)
		return
	}

	for _, multipartUploadSummary := range listMultipartUploadsResponse.Uploads {
		fmt.Println(multipartUploadSummary.Key, multipartUploadSummary.UploadId, multipartUploadSummary.Initiated)
	}

	for _, prefix := range listMultipartUploadsResponse.GetCommonPrefixes() {
		fmt.Println(prefix)
	}
}

func RunBosExamples() {
	listParts()
	return
	listPartsFromRequest()
	listMultipartUploads()
	listMultipartUploadsFromRequest()
	abortMultipartUpload()
	multipartUploadFromFile()
	multipartUpload()
	appendObject()
	generatePresignedUrl()
	getObjectMetadata()
	getObjectToFile()
	getObjectFromRequest()
	getObject()
	copyObjectFromRequest()
	copyObject()
	deleteObject()
	deleteMultipleObjects()
	listObjects()
	listObjectsFromRequest()
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
