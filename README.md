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
go test github.com/guoyao/baidubce-sdk-go/bos
go test github.com/guoyao/baidubce-sdk-go/bce
go test github.com/guoyao/baidubce-sdk-go/util
```

## Usage

```go
import (
	"github.com/guoyao/baidubce-sdk-go/bce"
	"github.com/guoyao/baidubce-sdk-go/bos"
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
    bucketName := "baidubce-sdk-go"

    /* ------ put object from string --------  */
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

    /* ------ put object from bytes --------  */
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

    /* ------ put object from file --------  */
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
```

### MultipartUpload

```go
func MultipartUpload() {
	bucketName := "baidubce-sdk-go"
	objectKey := "test-multipart-upload.zip"

	initiateMultipartUploadRequest := bos.InitiateMultipartUploadRequest{
		BucketName: bucketName,
		ObjectKey:  objectKey,
	}

	initiateMultipartUploadResponse, bceError := bosClient.InitiateMultipartUpload(
        initiateMultipartUploadRequest, nil)

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
```

### MultipartUploadFromFile

```go
func MultipartUploadFromFile() {
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
```

### Others

More api usages please refer

* [examples/bos.go](examples/bos.go)
* [bos/client_test.go](bos/client_test.go)

## Authors

**Guoyao Wu**

+ [http://guoyao.me](http://guoyao.me)
+ [http://github.com/guoyao](http://github.com/guoyao)
