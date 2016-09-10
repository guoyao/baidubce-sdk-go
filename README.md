# baidubce-sdk-go

Unofficial Go SDK for Baidu Cloud Engine

1. BOS (Baidu Object Storage)

<http://bce.baidu.com/index.html>

## Install

```
go get github.com/guoyao/baidubce-sdk-go
```
## Run Test

Before run test, you should setup two environment variables: `BAIDU_BCE_AK` and `BAIDU_BCE_SK`

```
go test -v github.com/guoyao/baidubce-sdk-go/...
```

## Usage

```go
import (
	"github.com/guoyao/baidubce-sdk-go/bce"
	"github.com/guoyao/baidubce-sdk-go/bos"
    "github.com/guoyao/baidubce-sdk-go/util"
)

var credentials = bce.NewCredentials("AK", "SK")
var bceConfig = bce.NewConfig(credentials)
var bosConfig = bos.NewConfig(bceConfig)
var bosClient = bos.NewClient(bosConfig)
```

### CreateBucket

```go
func CreateBucket() {
    bucketName := "baidubce-sdk-go"
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
```

### PutObject

```go
func PutObject() {
	/*------------------ put object from string --------------------*/
	bucketName := "baidubce-sdk-go"
	objectKey := "examples/put-object-from-string.txt"
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

	/*------------------ put object from bytes --------------------*/
	objectKey = "examples/put-object-from-bytes"
	byteArray := make([]byte, 1024, 1024)
	putObjectResponse, bceError = bosClient.PutObject(bucketName, objectKey, byteArray, nil, nil)

	if bceError != nil {
		log.Println(bceError)
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
		putObjectResponse, bceError = bosClient.PutObject(bucketName, objectKey, file, nil, nil)

		if bceError != nil {
			log.Println(bceError)
		} else {
			fmt.Println(putObjectResponse.GetETag())
		}
	}
}
```

### MultipartUpload

```go
func MultipartUpload() {
	bucketName := "baidubce-sdk-go"
	objectKey := "examples/test-multipart-upload"

	initiateMultipartUploadRequest := bos.InitiateMultipartUploadRequest{
		BucketName: bucketName,
		ObjectKey:  objectKey,
	}

	initiateMultipartUploadResponse, bceError := bosClient.InitiateMultipartUpload(initiateMultipartUploadRequest, nil)

	if bceError != nil {
		panic(bceError)
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

		uploadPartResponse, bceError := bosClient.UploadPart(uploadPartRequest, nil)

		if bceError != nil {
			panic(bceError)
		}

		parts[partNumber-1].ETag = uploadPartResponse.GetETag()
	}

	completeMultipartUploadRequest := bos.CompleteMultipartUploadRequest{
		BucketName: bucketName,
		ObjectKey:  objectKey,
		UploadId:   uploadId,
		Parts:      parts,
	}

	completeMultipartUploadResponse, bceError := bosClient.CompleteMultipartUpload(
		completeMultipartUploadRequest, nil)

	if bceError != nil {
		panic(bceError)
	}

	fmt.Println(completeMultipartUploadResponse.ETag)
}
```

### MultipartUploadFromFile

```go
func MultipartUploadFromFile() {
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

	completeMultipartUploadResponse, bceError := bosClient.MultipartUploadFromFile(bucketName,
		objectKey, file.Name(), partSize)

	if bceError != nil {
		log.Println(bceError)
	} else {
		fmt.Println(completeMultipartUploadResponse.ETag)
	}
}
```

### Others

More api usages please refer

* [examples/bos.go](examples/bos.go)
* [bos/client_test.go](bos/client_test.go)

## Authors

**Guoyao Wu**

+ [http://guoyao.me](http://guoyao.me)
+ [http://github.com/guoyao](http://github.com/guoyao)
