package examples

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	//"time"

	"github.com/guoyao/baidubce-sdk-go/bce"
	"github.com/guoyao/baidubce-sdk-go/bos"
	"github.com/guoyao/baidubce-sdk-go/util"
)

var credentials = bce.NewCredentials(os.Getenv("BAIDU_BCE_AK"), os.Getenv("BAIDU_BCE_SK"))

//var bceConfig = bce.NewConfig(credentials)
var bceConfig = &bce.Config{
	Credentials: credentials,
	Checksum:    true,
}
var bosConfig = bos.NewConfig(bceConfig)
var bosClient = bos.NewClient(bosConfig)

func init() {
	bosClient.SetDebug(true)

	/*
		bceConfig.Endpoint = "baidubce-sdk-go.bj.bcebos.com"
		bceConfig.ProxyHost = "agent.baidu.com"
		bceConfig.ProxyPort = 8118
		bceConfig.MaxConnections = 6
		bceConfig.Timeout = 6 * time.Second
	*/
}

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
	/*------------------ put object from string --------------------*/
	bucketName := "baidubce-sdk-go"
	objectKey := "examples/put-object-from-string.txt"
	str := "Hello World 你好"

	option := new(bce.SignOption)
	metadata := new(bos.ObjectMetadata)
	metadata.AddUserMetadata("x-bce-meta-name", "guoyao")
	putObjectResponse, err := bosClient.PutObject(bucketName, objectKey, str, metadata, option)

	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(putObjectResponse.GetETag())
	}

	return

	/*------------------ put object from bytes --------------------*/
	objectKey = "examples/put-object-from-bytes"
	byteArray := make([]byte, 1024, 1024)
	putObjectResponse, err = bosClient.PutObject(bucketName, objectKey, byteArray, nil, nil)

	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(putObjectResponse.GetETag())
	}

	/*------------------ put object from file --------------------*/
	file, err := util.TempFileWithSize(1024)

	defer func() {
		if file != nil {
			file.Close()
			os.Remove(file.Name())
		}
	}()

	if err != nil {
		log.Fatal(err)
	}

	objectKey = "examples/put-object-from-file"

	if err != nil {
		log.Println(err)
	} else {
		putObjectResponse, err = bosClient.PutObject(bucketName, objectKey, file, nil, nil)

		if err != nil {
			log.Println(err)
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

	putObjectResponse, err := bosClient.PutObject(bucketName, objectKey, str, metadata, option)

	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(putObjectResponse.GetETag())
	}

	err = bosClient.DeleteObject(bucketName, objectKey, nil)

	if err != nil {
		log.Println(err)
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
		putObjectResponse, err := bosClient.PutObject(bucketName, value, str, nil, nil)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(putObjectResponse.GetETag())
	}

	deleteMultipleObjectsResponse, err := bosClient.DeleteMultipleObjects(bucketName, objects, nil)

	if err != nil {
		log.Println(err)
	} else if deleteMultipleObjectsResponse != nil {
		for _, deleteMultipleObjectsError := range deleteMultipleObjectsResponse.Errors {
			log.Println(deleteMultipleObjectsError.Error())
		}
	}
}

func listObjects() {
	bucketName := "baidubce-sdk-go"

	listObjectsResponse, err := bosClient.ListObjects(bucketName, nil)

	if err != nil {
		log.Println(err)
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
		//Marker:    "examples/put-object-from-bytes",
		//Prefix:    "examples/",
		MaxKeys: 100,
	}

	listObjectsResponse, err := bosClient.ListObjectsFromRequest(listObjectsRequest, nil)

	if err != nil {
		log.Println(err)
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
	srcKey := "examples/test-copy-src"
	destBucketName := "baidubce-sdk-go"
	destKey := "examples/test-copy-dest"

	copyObjectResponse, err := bosClient.CopyObject(srcBucketName, srcKey, destBucketName, destKey, nil)

	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(copyObjectResponse.ETag, copyObjectResponse.LastModified)
	}
}

func copyObjectFromRequest() {
	etag := "fa412a6ca6d415208be69bc4a00f4103"

	copyObjectRequest := bos.CopyObjectRequest{
		SrcBucketName:  "baidubce-sdk-go",
		SrcKey:         "examples/test-copy-src",
		DestBucketName: "baidubce-sdk-go",
		DestKey:        "examples/test-copy-dest",
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

	copyObjectResponse, err := bosClient.CopyObjectFromRequest(copyObjectRequest, nil)

	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(copyObjectResponse.ETag, copyObjectResponse.LastModified)
	}
}

func getObject() {
	bucketName := "baidubce-sdk-go"
	objectKey := "examples/test-get-object"

	object, err := bosClient.GetObject(bucketName, objectKey, nil)

	if err != nil {
		log.Println(err)
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
	objectKey := "examples/test-get-object-from-request"

	getObjectRequest := bos.GetObjectRequest{
		BucketName: bucketName,
		ObjectKey:  objectKey,
	}
	getObjectRequest.SetRange(0, 1000)

	object, err := bosClient.GetObjectFromRequest(getObjectRequest, nil)

	if err != nil {
		log.Println(err)
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
	objectKey := "examples/test-get-object-to-file"

	getObjectRequest := &bos.GetObjectRequest{
		BucketName: bucketName,
		ObjectKey:  objectKey,
	}
	getObjectRequest.SetRange(0, 1000)

	file, err := os.OpenFile(objectKey, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		log.Println(err)
	} else {
		objectMetadata, err := bosClient.GetObjectToFile(getObjectRequest, file, nil)

		if err != nil {
			log.Println(err)
		} else {
			fmt.Println(objectMetadata)
		}
	}
}

func getObjectMetadata() {
	bucketName := "baidubce-sdk-go"
	objectKey := "examples/test-get-object-metedata"

	objectMetadata, err := bosClient.GetObjectMetadata(bucketName, objectKey, nil)

	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(objectMetadata)
	}
}

func generatePresignedUrl() {
	bucketName := "baidubce-sdk-go"
	objectKey := "examples/test-generate-presigned-url"

	option := &bce.SignOption{
		ExpirationPeriodInSeconds: 300,
	}

	url, err := bosClient.GeneratePresignedUrl(bucketName, objectKey, option)

	if err != nil {
		log.Println(err)
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

	appendObjectResponse, err := bosClient.AppendObject(bucketName, objectKey, offset, str, metadata, option)

	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(appendObjectResponse.GetETag(), appendObjectResponse.GetMD5(),
			appendObjectResponse.GetNextAppendOffset())
	}
}

func multipartUpload() {
	bucketName := "baidubce-sdk-go"
	objectKey := "examples/test-multipart-upload"

	initiateMultipartUploadRequest := bos.InitiateMultipartUploadRequest{
		BucketName: bucketName,
		ObjectKey:  objectKey,
	}

	initiateMultipartUploadResponse, err := bosClient.InitiateMultipartUpload(initiateMultipartUploadRequest, nil)

	if err != nil {
		panic(err)
	}

	uploadId := initiateMultipartUploadResponse.UploadId

	files := make([]*os.File, 0)
	file, err := util.TempFileWithSize(1024 * 1024 * 6)
	files = append(files, file)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		for _, f := range files {
			f.Close()
			os.Remove(f.Name())
		}
	}()

	fileInfo, err := file.Stat()

	if err != nil {
		log.Fatal(err)
	}

	var partSize int64 = 1024 * 1024 * 5
	var totalSize int64 = fileInfo.Size()
	var partCount int = int(math.Ceil(float64(totalSize) / float64(partSize)))

	parts := make([]bos.PartSummary, 0, partCount)

	for i := 0; i < partCount; i++ {
		var skipBytes int64 = partSize * int64(i)
		var size int64 = int64(math.Min(float64(totalSize-skipBytes), float64(partSize)))

		tempFile, err := util.TempFile(nil, "", "")
		files = append(files, tempFile)

		if err != nil {
			panic(err)
		}

		limitReader := io.LimitReader(file, size)
		_, err = io.Copy(tempFile, limitReader)

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
			PartData:   tempFile,
		}

		parts = append(parts, bos.PartSummary{PartNumber: partNumber})

		uploadPartResponse, err := bosClient.UploadPart(uploadPartRequest, nil)

		if err != nil {
			panic(err)
		}

		parts[partNumber-1].ETag = uploadPartResponse.GetETag()
	}

	completeMultipartUploadRequest := bos.CompleteMultipartUploadRequest{
		BucketName: bucketName,
		ObjectKey:  objectKey,
		UploadId:   uploadId,
		Parts:      parts,
	}

	completeMultipartUploadResponse, err := bosClient.CompleteMultipartUpload(
		completeMultipartUploadRequest, nil)

	if err != nil {
		panic(err)
	}

	fmt.Println(completeMultipartUploadResponse.ETag)
}

func multipartUploadFromFile() {
	bucketName := "baidubce-sdk-go"
	objectKey := "examples/test-multipart-upload-from-file"

	file, err := util.TempFileWithSize(1024 * 1024 * 10)

	defer func() {
		if file != nil {
			file.Close()
			os.Remove(file.Name())
		}
	}()

	if err != nil {
		log.Fatal(err)
	}

	var partSize int64 = 1024 * 1024 * 2

	completeMultipartUploadResponse, err := bosClient.MultipartUploadFromFile(bucketName,
		objectKey, file.Name(), partSize)

	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(completeMultipartUploadResponse.ETag)
	}
}

func abortMultipartUpload() {
	bucketName := "baidubce-sdk-go"
	objectKey := "examples/test-multipart-upload"

	initiateMultipartUploadRequest := bos.InitiateMultipartUploadRequest{
		BucketName: bucketName,
		ObjectKey:  objectKey,
	}

	initiateMultipartUploadResponse, err := bosClient.InitiateMultipartUpload(initiateMultipartUploadRequest, nil)

	if err != nil {
		panic(err)
	}

	uploadId := initiateMultipartUploadResponse.UploadId

	abortMultipartUploadRequest := bos.AbortMultipartUploadRequest{
		BucketName: bucketName,
		ObjectKey:  objectKey,
		UploadId:   uploadId,
	}

	err = bosClient.AbortMultipartUpload(abortMultipartUploadRequest, nil)

	if err != nil {
		log.Println(err)
	}
}

func abortAllMultipartUpload(bucketName string) {
	listMultipartUploadsResponse, err := bosClient.ListMultipartUploads(bucketName, nil)

	if err != nil {
		log.Println(err)
		return
	}

	for _, multipartUploadSummary := range listMultipartUploadsResponse.Uploads {
		abortMultipartUploadRequest := bos.AbortMultipartUploadRequest{
			BucketName: bucketName,
			ObjectKey:  multipartUploadSummary.Key,
			UploadId:   multipartUploadSummary.UploadId,
		}

		err = bosClient.AbortMultipartUpload(abortMultipartUploadRequest, nil)

		if err != nil {
			log.Println(err)
		}
	}
}

func listParts() {
	bucketName := "baidubce-sdk-go"
	objectKey := "examples/test-multipart-upload"
	uploadId := "a977803dc94b8c9da2c9f1d8432a1805"

	listPartsResponse, err := bosClient.ListParts(bucketName, objectKey, uploadId, nil)

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
	objectKey := "examples/test-multipart-upload"

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
		objectKey := "examples/test-multipart-upload"

		initiateMultipartUploadRequest := bos.InitiateMultipartUploadRequest{
			BucketName: bucketName,
			ObjectKey:  objectKey,
		}

		initiateMultipartUploadResponse, err := bosClient.InitiateMultipartUpload(initiateMultipartUploadRequest, nil)

		if err != nil {
			log.Println(err)
			log.Println(initiateMultipartUploadResponse.UploadId)
			return
		}
	*/

	listMultipartUploadsRequest := bos.ListMultipartUploadsRequest{
		BucketName: "baidubce-sdk-go",
		//Delimiter:  "/",
		//KeyMarker:  "examples/test-multipart-upload",
		//Prefix:     "examples/",
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

func getBucketCors() {
	bucketName := "baidubce-sdk-go"
	bucketCors, err := bosClient.GetBucketCors(bucketName, nil)

	if err != nil {
		log.Println(err)
	} else {
		for _, bucketCorsItem := range bucketCors.CorsConfiguration {
			fmt.Println(bucketCorsItem.AllowedOrigins)
			fmt.Println(bucketCorsItem.AllowedMethods)
			fmt.Println(bucketCorsItem.AllowedHeaders)
			fmt.Println(bucketCorsItem.AllowedExposeHeaders)
			fmt.Println(bucketCorsItem.MaxAgeSeconds)
		}
	}
}

func setBucketCors() {
	bucketName := "baidubce-sdk-go"
	bucketCors := bos.BucketCors{
		CorsConfiguration: []bos.BucketCorsItem{
			bos.BucketCorsItem{
				AllowedOrigins:       []string{"http://*", "https://*"},
				AllowedMethods:       []string{"GET", "HEAD", "POST", "PUT"},
				AllowedHeaders:       []string{"*"},
				AllowedExposeHeaders: []string{"ETag", "x-bce-request-id", "Content-Type"},
				MaxAgeSeconds:        3600,
			},
			bos.BucketCorsItem{
				AllowedOrigins:       []string{"http://www.example.com", "www.example2.com"},
				AllowedMethods:       []string{"GET", "HEAD", "DELETE"},
				AllowedHeaders:       []string{"Authorization", "x-bce-test", "x-bce-test2"},
				AllowedExposeHeaders: []string{"user-custom-expose-header"},
				MaxAgeSeconds:        2000,
			},
		},
	}

	err := bosClient.SetBucketCors(bucketName, bucketCors, nil)

	if err != nil {
		log.Println(err)
	}
}

func deleteBucketCors() {
	bucketName := "baidubce-sdk-go"
	err := bosClient.DeleteBucketCors(bucketName, nil)

	if err != nil {
		log.Println(err)
	}
}

func optionsObject() {
	bucketName := "baidubce-sdk-go"
	objectKey := "put-object-from-file"

	resp, err := bosClient.OptionsObject(bucketName, objectKey, "http://www.example.com", "GET", "x-bce-test")

	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(resp.Header)
	}
}

func setBucketLogging() {
	bucketName := "baidubce-sdk-go"
	targetBucket := "bucket-logs"
	targetPrefix := "baidubce-sdk-go"
	err := bosClient.SetBucketLogging(bucketName, targetBucket, targetPrefix, nil)

	if err != nil {
		log.Println(err)
	}
}

func getBucketLogging() {
	bucketName := "baidubce-sdk-go"
	bucketLogging, err := bosClient.GetBucketLogging(bucketName, nil)

	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(bucketLogging)
	}
}

func deleteBucketLogging() {
	bucketName := "baidubce-sdk-go"
	err := bosClient.DeleteBucketLogging(bucketName, nil)

	if err != nil {
		log.Println(err)
	}
}

func RunBosExamples() {
	listBuckets()
	return
	//abortAllMultipartUpload("docker-registry-me-test")
	deleteBucketLogging()
	getBucketLogging()
	setBucketLogging()
	optionsObject()
	deleteBucketCors()
	setBucketCors()
	getBucketCors()
	listParts()
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
