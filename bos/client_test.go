package bos

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/guoyao/baidubce-sdk-go/bce"
	"github.com/guoyao/baidubce-sdk-go/util"
)

var credentials = bce.NewCredentials(os.Getenv("BAIDU_BCE_AK"), os.Getenv("BAIDU_BCE_SK"))

//var bceConfig = bce.NewConfig(credentials)
var bceConfig = &bce.Config{
	Credentials: credentials,
	Checksum:    true,
}
var bosConfig = NewConfig(bceConfig)
var bosClient = NewClient(bosConfig)

func TestGetBucketLocation(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-get-bucket-location-"
	method := "GetBucketLocation"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		expected := "bj"
		location, _ := bosClient.GetBucketLocation(bucketName, nil)

		if location.LocationConstraint != expected {
			t.Error(util.FormatTest(method, location.LocationConstraint, expected))
		}
	})
}

func TestListBuckets(t *testing.T) {
	_, err := bosClient.ListBuckets(nil)

	if err != nil {
		t.Error(util.FormatTest("ListBuckets", err.Error(), "nil"))
	}
}

func TestCreateBucket(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-create-bucket-"
	method := "CreateBucket"

	around(t, method, bucketNamePrefix, "", nil)
}

func TestDoesBucketExist(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-does-bucket-exist-"
	method := "DoesBucketExist"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		expected := true
		exists, err := bosClient.DoesBucketExist(bucketName, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), strconv.FormatBool(expected)))
		} else if exists != expected {
			t.Error(util.FormatTest(method, strconv.FormatBool(exists), strconv.FormatBool(expected)))
		}
	})

}

func TestDeleteBucket(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-delete-bucket-"
	method := "DeleteBucket"

	around(t, method, bucketNamePrefix, "", nil)
}

func TestSetBucketPrivate(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-set-bucket-private-"
	method := "SetBucketPrivate"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		err := bosClient.SetBucketPrivate(bucketName, nil)
		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		}
	})
}

func TestSetBucketPublicRead(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-set-bucket-public-read-"
	method := "SetBucketPublicRead"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		err := bosClient.SetBucketPublicRead(bucketName, nil)
		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		}
	})
}

func TestSetBucketPublicReadWrite(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-set-bucket-public-rw-"
	method := "SetBucketPublicReadWrite"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		err := bosClient.SetBucketPublicReadWrite(bucketName, nil)
		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		}
	})
}

func TestGetBucketAcl(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-get-bucket-acl-"
	method := "GetBucketAcl"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		_, err := bosClient.GetBucketAcl(bucketName, nil)
		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		}
	})
}

func TestSetBucketAcl(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-set-bucket-acl-"
	method := "SetBucketAcl"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		bucketAcl := BucketAcl{
			AccessControlList: []Grant{
				Grant{
					Grantee: []BucketGrantee{
						BucketGrantee{Id: "ef5a4b19192f4931adcf0e12f82795e2"},
					},
					Permission: []string{"FULL_CONTROL"},
				},
			},
		}
		if err := bosClient.SetBucketAcl(bucketName, bucketAcl, nil); err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		}
	})
}

func TestPubObject(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-put-object-"
	method := "PutObject"
	objectKey := "put-object-from-string.txt"
	str := "Hello World 你好"

	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		_, err := bosClient.PutObject(bucketName, objectKey, str, nil, nil)
		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		}
	})
}

func TestDeleteObject(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-delete-object-"
	method := "DeleteObject"
	objectKey := "put-object-from-string.txt"
	str := "Hello World 你好"

	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		_, err := bosClient.PutObject(bucketName, objectKey, str, nil, nil)
		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		}
	})
}

func TestDeleteMultipleObjects(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-delete-multiple-objects-"
	method := "DeleteMultipleObjects"
	str := "Hello World 你好"

	objects := []string{
		"examples/delete-multiple-objects/put-object-from-string.txt",
		"examples/delete-multiple-objects/put-object-from-string-2.txt",
		"examples/delete-multiple-objects/put-object-from-string-3.txt",
	}

	around(t, method, bucketNamePrefix, objects, func(bucketName string) {
		for _, objectKey := range objects {
			_, err := bosClient.PutObject(bucketName, objectKey, str, nil, nil)

			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
			}
		}
	})
}

func TestListObjects(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-list-objects-"
	method := "ListObjects"
	objectKey := "put-object-from-string.txt"
	str := "Hello World 你好"

	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		_, err := bosClient.PutObject(bucketName, objectKey, str, nil, nil)
		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else {
			listObjectsResponse, err := bosClient.ListObjects(bucketName, nil)
			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
			} else if length := len(listObjectsResponse.Contents); length != 1 {
				t.Error(util.FormatTest(method, strconv.Itoa(length), "1"))
			}
		}
	})
}

func TestCopyObject(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-copy-object-"
	method := "CopyObject"
	objectKey := "put-object-from-string.txt"
	str := "Hello World 你好"

	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		_, err := bosClient.PutObject(bucketName, objectKey, str, nil, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else {
			destKey := "put-object-from-string-copy.txt"
			_, err := bosClient.CopyObject(bucketName, objectKey, bucketName, destKey, nil)

			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
			} else {
				listObjectsResponse, err := bosClient.ListObjects(bucketName, nil)

				if err != nil {
					t.Error(util.FormatTest(method, err.Error(), "nil"))
				} else if length := len(listObjectsResponse.Contents); length != 2 {
					t.Error(util.FormatTest(method, strconv.Itoa(length), "2"))
				} else {
					err = bosClient.DeleteObject(bucketName, destKey, nil)

					if err != nil {
						t.Error(util.FormatTest(method+" at deleting object", err.Error(), "nil"))
					}
				}
			}
		}
	})
}

func TestCopyObjectFromRequest(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-copy-object-from-request-"
	method := "CopyObjectFromRequest"
	objectKey := "put-object-from-string.txt"
	str := "Hello World 你好"

	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		_, err := bosClient.PutObject(bucketName, objectKey, str, nil, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else {
			destKey := "put-object-from-string-copy.txt"

			copyObjectRequest := CopyObjectRequest{
				SrcBucketName:  bucketName,
				SrcKey:         objectKey,
				DestBucketName: bucketName,
				DestKey:        destKey,
				ObjectMetadata: &ObjectMetadata{
					CacheControl: "no-cache",
					UserMetadata: map[string]string{
						"test-user-metadata": "test user metadata",
						"x-bce-meta-name":    "x-bce-meta-name",
					},
				},
			}

			_, err := bosClient.CopyObjectFromRequest(copyObjectRequest, nil)

			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
			} else {
				listObjectsResponse, err := bosClient.ListObjects(bucketName, nil)

				if err != nil {
					t.Error(util.FormatTest(method, err.Error(), "nil"))
				} else if length := len(listObjectsResponse.Contents); length != 2 {
					t.Error(util.FormatTest(method, strconv.Itoa(length), "2"))
				} else {
					err = bosClient.DeleteObject(bucketName, destKey, nil)

					if err != nil {
						t.Error(util.FormatTest(method+" at deleting object", err.Error(), "nil"))
					}
				}
			}
		}
	})
}

func TestGetObject(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-get-object-"
	method := "GetObject"
	objectKey := "put-object-from-string.txt"
	str := "Hello World 你好"

	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		_, err := bosClient.PutObject(bucketName, objectKey, str, nil, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else {
			object, err := bosClient.GetObject(bucketName, objectKey, nil)

			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
			} else if object.ObjectMetadata.ETag == "" {
				t.Error(util.FormatTest(method, "etag is empty", "non empty etag"))
			} else {
				byteArray, err := ioutil.ReadAll(object.ObjectContent)

				if err != nil {
					t.Error(util.FormatTest(method, err.Error(), "nil"))
				} else if len(byteArray) == 0 {
					t.Error(util.FormatTest(method, "body is empty", "non empty body"))
				}
			}
		}
	})
}

func TestGetObjectFromRequest(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-get-object-from-request-"
	method := "GetObjectFromRequest"
	objectKey := "put-object-from-string.txt"
	str := "Hello World 你好"

	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		_, err := bosClient.PutObject(bucketName, objectKey, str, nil, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else {
			getObjectRequest := GetObjectRequest{
				BucketName: bucketName,
				ObjectKey:  objectKey,
			}
			getObjectRequest.SetRange(0, 1000)
			object, err := bosClient.GetObjectFromRequest(getObjectRequest, nil)

			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
			} else if object.ObjectMetadata.ETag == "" {
				t.Error(util.FormatTest(method, "etag is empty", "non empty etag"))
			} else {
				byteArray, err := ioutil.ReadAll(object.ObjectContent)

				if err != nil {
					t.Error(util.FormatTest(method, err.Error(), "nil"))
				} else if len(byteArray) == 0 {
					t.Error(util.FormatTest(method, "body is empty", "non empty body"))
				}
			}
		}
	})
}

func TestGetObjectToFile(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-get-object-to-file-"
	method := "GetObjectToFile"
	objectKey := "put-object-from-string.txt"
	str := "Hello World 你好"

	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		_, err := bosClient.PutObject(bucketName, objectKey, str, nil, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else {
			getObjectRequest := &GetObjectRequest{
				BucketName: bucketName,
				ObjectKey:  objectKey,
			}
			getObjectRequest.SetRange(0, 1000)

			file, err := os.OpenFile(objectKey, os.O_WRONLY|os.O_CREATE, 0666)

			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
			} else {
				objectMetadata, err := bosClient.GetObjectToFile(getObjectRequest, file, nil)

				if err != nil {
					t.Error(util.FormatTest(method, err.Error(), "nil"))
				} else if objectMetadata.ETag == "" {
					t.Error(util.FormatTest(method, "etag is empty", "non empty etag"))
				} else if !util.CheckFileExists(objectKey) {
					t.Error(util.FormatTest(method, "file is not saved to local", "file saved to local"))
				} else {
					err := os.Remove(objectKey)

					if err != nil {
						t.Error(util.FormatTest(method, err.Error(), "nil"))
					}
				}
			}
		}
	})
}

func TestGetObjectMetadata(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-get-object-metadata-"
	method := "GetObjectMetadata"
	objectKey := "put-object-from-string.txt"
	str := "Hello World 你好"

	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		_, err := bosClient.PutObject(bucketName, objectKey, str, nil, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else {
			objectMetadata, err := bosClient.GetObjectMetadata(bucketName, objectKey, nil)

			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
			} else if objectMetadata.ETag == "" {
				t.Error(util.FormatTest(method, "etag is empty", "non empty etag"))
			}
		}
	})
}

func TestGeneratePresignedUrl(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-generate-presigned-url-"
	method := "GeneratePresignedUrl"
	objectKey := "put-object-from-string.txt"
	str := "Hello World 你好"

	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		_, err := bosClient.PutObject(bucketName, objectKey, str, nil, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else {
			url, err := bosClient.GeneratePresignedUrl(bucketName, objectKey, nil)

			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
			} else {
				req, err := http.NewRequest("GET", url, nil)

				if err != nil {
					t.Error(util.FormatTest(method, err.Error(), "nil"))
				} else {
					httpClient := http.Client{}
					res, err := httpClient.Do(req)

					if err != nil {
						t.Error(util.FormatTest(method, err.Error(), "nil"))
					} else if res.StatusCode != 200 {
						t.Error(util.FormatTest(method, fmt.Sprintf("status code: %v", res.StatusCode), "status code: 200"))
					}
				}
			}
		}
	})
}

func TestAppendObject(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-append-object-"
	method := "AppendObject"
	objectKey := "append-object-from-string.txt"
	str := "Hello World 你好"
	offset := 0

	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		appendObjectResponse, err := bosClient.AppendObject(bucketName, objectKey, offset, str, nil, nil)
		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else if appendObjectResponse.GetETag() == "" || appendObjectResponse.GetNextAppendOffset() == "" {
			t.Error(util.FormatTest(method, "etag and next append offset are not exists", "etag and next append offset"))
		} else {
			length, err := strconv.Atoi(appendObjectResponse.GetNextAppendOffset())

			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
			} else {
				offset = length
				appendObjectResponse, err := bosClient.AppendObject(bucketName, objectKey, offset, str, nil, nil)

				if err != nil {
					t.Error(util.FormatTest(method, err.Error(), "nil"))
				} else if appendObjectResponse.GetETag() == "" || appendObjectResponse.GetNextAppendOffset() == "" {
					t.Error(util.FormatTest(method, "etag and next append offset are not exists", "etag and next append offset"))
				} else {
					length, err := strconv.Atoi(appendObjectResponse.GetNextAppendOffset())

					if err != nil {
						t.Error(util.FormatTest(method, err.Error(), "nil"))
					} else if length != offset*2 {
						t.Error(util.FormatTest(method, strconv.Itoa(length), strconv.Itoa(offset*2)))
					}
				}
			}
		}
	})
}

func TestMultipartUploadFromFile(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-multipart-upload-from-file-"
	method := "MultipartUploadFromFile"
	objectKey := "test-multipart-upload"

	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		file, err := util.TempFileWithSize(1024 * 1024 * 6)

		defer func() {
			if file != nil {
				file.Close()
				os.Remove(file.Name())
			}
		}()

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
			return
		}

		var partSize int64 = 1024 * 1024 * 2

		completeMultipartUploadResponse, err := bosClient.MultipartUploadFromFile(bucketName,
			objectKey, file.Name(), partSize)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else if completeMultipartUploadResponse.ETag == "" {
			t.Error(util.FormatTest(method, "etag is not exists", "etag"))
		}
	})
}

func TestAbortMultipartUpload(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-abort-multipart-upload-"
	method := "AbortMultipartUpload"
	objectKey := "test-multipart-upload"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		initiateMultipartUploadRequest := InitiateMultipartUploadRequest{
			BucketName: bucketName,
			ObjectKey:  objectKey,
		}

		initiateMultipartUploadResponse, err := bosClient.InitiateMultipartUpload(initiateMultipartUploadRequest, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
			return
		}

		uploadId := initiateMultipartUploadResponse.UploadId

		abortMultipartUploadRequest := AbortMultipartUploadRequest{
			BucketName: bucketName,
			ObjectKey:  objectKey,
			UploadId:   uploadId,
		}

		err = bosClient.AbortMultipartUpload(abortMultipartUploadRequest, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
			return
		}
	})
}

func TestListMultipartUploads(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-list-multipart-uploads-"
	objectKey := "test-multipart-upload"
	method := "ListMultipartUploads"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		initiateMultipartUploadRequest := InitiateMultipartUploadRequest{
			BucketName: bucketName,
			ObjectKey:  objectKey,
		}

		initiateMultipartUploadResponse, err := bosClient.InitiateMultipartUpload(initiateMultipartUploadRequest, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
			return
		}

		defer func() {
			if initiateMultipartUploadResponse != nil {
				abortMultipartUploadRequest := AbortMultipartUploadRequest{
					BucketName: bucketName,
					ObjectKey:  objectKey,
					UploadId:   initiateMultipartUploadResponse.UploadId,
				}

				err = bosClient.AbortMultipartUpload(abortMultipartUploadRequest, nil)

				if err != nil {
					t.Error(util.FormatTest(method, err.Error(), "nil"))
				}
			}
		}()

		listMultipartUploadsResponse, err := bosClient.ListMultipartUploads(bucketName, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
			return
		}

		partCount := len(listMultipartUploadsResponse.Uploads)

		if partCount != 1 {
			t.Error(util.FormatTest(method, fmt.Sprintf("part count is %d", partCount), "part count should be 1"))
		}
	})
}

func TestListParts(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-list-parts-"
	objectKey := "test-list-parts"
	method := "ListParts"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		initiateMultipartUploadRequest := InitiateMultipartUploadRequest{
			BucketName: bucketName,
			ObjectKey:  objectKey,
		}

		initiateMultipartUploadResponse, err := bosClient.InitiateMultipartUpload(initiateMultipartUploadRequest, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
			return
		}

		defer func() {
			if initiateMultipartUploadResponse != nil {
				abortMultipartUploadRequest := AbortMultipartUploadRequest{
					BucketName: bucketName,
					ObjectKey:  objectKey,
					UploadId:   initiateMultipartUploadResponse.UploadId,
				}

				err := bosClient.AbortMultipartUpload(abortMultipartUploadRequest, nil)

				if err != nil {
					t.Error(util.FormatTest(method, err.Error(), "nil"))
				}
			}
		}()

		files := make([]*os.File, 0)
		file, err := util.TempFileWithSize(1024 * 1024 * 6)
		files = append(files, file)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
			return
		}

		defer func() {
			for _, f := range files {
				f.Close()
				os.Remove(f.Name())
			}
		}()

		fileInfo, err := file.Stat()

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
			return
		}

		var partSize int64 = 1024 * 1024 * 5
		var totalSize int64 = fileInfo.Size()
		var partCount int = int(math.Ceil(float64(totalSize) / float64(partSize)))

		var waitGroup sync.WaitGroup
		parts := make([]PartSummary, 0, partCount)

		for i := 0; i < partCount; i++ {
			var skipBytes int64 = partSize * int64(i)
			var size int64 = int64(math.Min(float64(totalSize-skipBytes), float64(partSize)))

			tempFile, err := util.TempFile(nil, "", "")
			files = append(files, tempFile)

			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
				return
			}

			limitReader := io.LimitReader(file, size)
			_, err = io.Copy(tempFile, limitReader)

			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
				return
			}

			partNumber := i + 1

			uploadPartRequest := UploadPartRequest{
				BucketName: bucketName,
				ObjectKey:  objectKey,
				UploadId:   initiateMultipartUploadResponse.UploadId,
				PartSize:   size,
				PartNumber: partNumber,
				PartData:   tempFile,
			}

			waitGroup.Add(1)

			parts = append(parts, PartSummary{PartNumber: partNumber})

			go func(partNumber int) {
				defer waitGroup.Done()

				uploadPartResponse, err := bosClient.UploadPart(uploadPartRequest, nil)

				if err != nil {
					t.Error(util.FormatTest(method, err.Error(), "nil"))
					return
				}

				parts[partNumber-1].ETag = uploadPartResponse.GetETag()
			}(partNumber)
		}

		waitGroup.Wait()

		listPartsResponse, err := bosClient.ListParts(bucketName, objectKey, initiateMultipartUploadResponse.UploadId, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
			return
		}

		partCount = len(listPartsResponse.Parts)

		if partCount != 2 {
			t.Error(util.FormatTest(method, fmt.Sprintf("part count is %d", partCount), "part count should be 2"))
		}
	})
}

func around(t *testing.T, method, bucketNamePrefix string, objectKey interface{}, f func(string)) {
	bucketName := bucketNamePrefix + strconv.Itoa(int(time.Now().Unix()))
	err := bosClient.CreateBucket(bucketName, nil)

	if err != nil {
		t.Error(util.FormatTest(method+" at creating bucket", err.Error(), "nil"))
	} else {
		if f != nil {
			f(bucketName)

			if key, ok := objectKey.(string); ok {
				if key != "" {
					err = bosClient.DeleteObject(bucketName, key, nil)

					if err != nil {
						t.Error(util.FormatTest(method+" at deleting object", err.Error(), "nil"))
					}
				}
			} else if keys, ok := objectKey.([]string); ok {
				if len(keys) > 0 {
					deleteMultipleObjectsResponse, err := bosClient.DeleteMultipleObjects(bucketName, keys, nil)

					if err != nil {
						t.Error(util.FormatTest(method, err.Error(), "nil"))
					} else if deleteMultipleObjectsResponse != nil {
						str := ""

						for _, deleteMultipleObjectsError := range deleteMultipleObjectsResponse.Errors {
							str += deleteMultipleObjectsError.Error()
						}

						t.Error(util.FormatTest(method, str, "empty string"))
					}
				}
			} else {
				t.Error(util.FormatTest(method, "objectKey is not valid", "objectKey should be string or []string"))
			}
		}

		err = bosClient.DeleteBucket(bucketName, nil)

		if err != nil {
			t.Error(util.FormatTest(method+" at deleting bucket", err.Error(), "nil"))
		}
	}
}
