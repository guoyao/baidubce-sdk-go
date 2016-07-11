package bos

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"testing"
	"time"

	"baidubce/test"
	"baidubce/util"
)

var bosClient = DefaultClient

func TestGetBucketLocation(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-get-bucket-location-"
	method := "GetBucketLocation"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		expected := "bj"
		location, _ := bosClient.GetBucketLocation(bucketName, nil)

		if location.LocationConstraint != expected {
			t.Error(test.Format(method, location.LocationConstraint, expected))
		}
	})
}

func TestListBuckets(t *testing.T) {
	_, err := bosClient.ListBuckets(nil)

	if err != nil {
		t.Error(test.Format("ListBuckets", err.Error(), "nil"))
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
			t.Error(test.Format(method, err.Error(), strconv.FormatBool(expected)))
		} else if exists != expected {
			t.Error(test.Format(method, strconv.FormatBool(exists), strconv.FormatBool(expected)))
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
			t.Error(test.Format(method, err.Error(), "nil"))
		}
	})
}

func TestSetBucketPublicRead(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-set-bucket-public-read-"
	method := "SetBucketPublicRead"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		err := bosClient.SetBucketPublicRead(bucketName, nil)
		if err != nil {
			t.Error(test.Format(method, err.Error(), "nil"))
		}
	})
}

func TestSetBucketPublicReadWrite(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-set-bucket-public-rw-"
	method := "SetBucketPublicReadWrite"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		err := bosClient.SetBucketPublicReadWrite(bucketName, nil)
		if err != nil {
			t.Error(test.Format(method, err.Error(), "nil"))
		}
	})
}

func TestGetBucketAcl(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-get-bucket-acl-"
	method := "GetBucketAcl"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		_, err := bosClient.GetBucketAcl(bucketName, nil)
		if err != nil {
			t.Error(test.Format(method, err.Error(), "nil"))
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
			t.Error(test.Format(method, err.Error(), "nil"))
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
			t.Error(test.Format(method, err.Error(), "nil"))
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
			t.Error(test.Format(method, err.Error(), "nil"))
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
				t.Error(test.Format(method, err.Error(), "nil"))
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
			t.Error(test.Format(method, err.Error(), "nil"))
		} else {
			listObjectsResponse, err := bosClient.ListObjects(bucketName, nil, nil)
			if err != nil {
				t.Error(test.Format(method, err.Error(), "nil"))
			} else if length := len(listObjectsResponse.Contents); length != 1 {
				t.Error(test.Format(method, strconv.Itoa(length), "1"))
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
			t.Error(test.Format(method, err.Error(), "nil"))
		} else {
			destKey := "put-object-from-string-copy.txt"
			_, err := bosClient.CopyObject(bucketName, objectKey, bucketName, destKey, nil)

			if err != nil {
				t.Error(test.Format(method, err.Error(), "nil"))
			} else {
				listObjectsResponse, err := bosClient.ListObjects(bucketName, nil, nil)

				if err != nil {
					t.Error(test.Format(method, err.Error(), "nil"))
				} else if length := len(listObjectsResponse.Contents); length != 2 {
					t.Error(test.Format(method, strconv.Itoa(length), "2"))
				} else {
					err = bosClient.DeleteObject(bucketName, destKey, nil)

					if err != nil {
						t.Error(test.Format(method+" at deleting object", err.Error(), "nil"))
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
			t.Error(test.Format(method, err.Error(), "nil"))
		} else {
			destKey := "put-object-from-string-copy.txt"

			copyObjectRequest := &CopyObjectRequest{
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
				t.Error(test.Format(method, err.Error(), "nil"))
			} else {
				listObjectsResponse, err := bosClient.ListObjects(bucketName, nil, nil)

				if err != nil {
					t.Error(test.Format(method, err.Error(), "nil"))
				} else if length := len(listObjectsResponse.Contents); length != 2 {
					t.Error(test.Format(method, strconv.Itoa(length), "2"))
				} else {
					err = bosClient.DeleteObject(bucketName, destKey, nil)

					if err != nil {
						t.Error(test.Format(method+" at deleting object", err.Error(), "nil"))
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
			t.Error(test.Format(method, err.Error(), "nil"))
		} else {
			object, err := bosClient.GetObject(bucketName, objectKey, nil)

			if err != nil {
				t.Error(test.Format(method, err.Error(), "nil"))
			} else if object.ObjectMetadata.ETag == "" {
				t.Error(test.Format(method, "etag is empty", "non empty etag"))
			} else {
				byteArray, err := ioutil.ReadAll(object.ObjectContent)

				if err != nil {
					t.Error(test.Format(method, err.Error(), "nil"))
				} else if len(byteArray) == 0 {
					t.Error(test.Format(method, "body is empty", "non empty body"))
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
			t.Error(test.Format(method, err.Error(), "nil"))
		} else {
			getObjectRequest := &GetObjectRequest{
				BucketName: bucketName,
				ObjectKey:  objectKey,
			}
			getObjectRequest.SetRange(0, 1000)
			object, err := bosClient.GetObjectFromRequest(getObjectRequest, nil)

			if err != nil {
				t.Error(test.Format(method, err.Error(), "nil"))
			} else if object.ObjectMetadata.ETag == "" {
				t.Error(test.Format(method, "etag is empty", "non empty etag"))
			} else {
				byteArray, err := ioutil.ReadAll(object.ObjectContent)

				if err != nil {
					t.Error(test.Format(method, err.Error(), "nil"))
				} else if len(byteArray) == 0 {
					t.Error(test.Format(method, "body is empty", "non empty body"))
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
			t.Error(test.Format(method, err.Error(), "nil"))
		} else {
			getObjectRequest := &GetObjectRequest{
				BucketName: bucketName,
				ObjectKey:  objectKey,
			}
			getObjectRequest.SetRange(0, 1000)

			file, err := os.OpenFile(objectKey, os.O_WRONLY|os.O_CREATE, 0666)

			if err != nil {
				t.Error(test.Format(method, err.Error(), "nil"))
			} else {
				objectMetadata, err := bosClient.GetObjectToFile(getObjectRequest, file, nil)

				if err != nil {
					t.Error(test.Format(method, err.Error(), "nil"))
				} else if objectMetadata.ETag == "" {
					t.Error(test.Format(method, "etag is empty", "non empty etag"))
				} else if !util.CheckFileExists(objectKey) {
					t.Error(test.Format(method, "file is not saved to local", "file saved to local"))
				} else {
					err := os.Remove(objectKey)

					if err != nil {
						t.Error(test.Format(method, err.Error(), "nil"))
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
			t.Error(test.Format(method, err.Error(), "nil"))
		} else {
			objectMetadata, err := bosClient.GetObjectMetadata(bucketName, objectKey, nil)

			if err != nil {
				t.Error(test.Format(method, err.Error(), "nil"))
			} else if objectMetadata.ETag == "" {
				t.Error(test.Format(method, "etag is empty", "non empty etag"))
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
			t.Error(test.Format(method, err.Error(), "nil"))
		} else {
			url, err := bosClient.GeneratePresignedUrl(bucketName, objectKey, nil)

			if err != nil {
				t.Error(test.Format(method, err.Error(), "nil"))
			} else {
				req, err := http.NewRequest("GET", url, nil)

				if err != nil {
					t.Error(test.Format(method, err.Error(), "nil"))
				} else {
					httpClient := http.Client{}
					res, err := httpClient.Do(req)

					if err != nil {
						t.Error(test.Format(method, err.Error(), "nil"))
					} else if res.StatusCode != 200 {
						t.Error(test.Format(method, fmt.Sprintf("status code: %v", res.StatusCode), "status code: 200"))
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
			t.Error(test.Format(method, err.Error(), "nil"))
		} else if appendObjectResponse.GetETag() == "" || appendObjectResponse.GetNextAppendOffset() == "" {
			t.Error(test.Format(method, "etag and next append offset are not exists", "etag and next append offset"))
		} else {
			length, err := strconv.Atoi(appendObjectResponse.GetNextAppendOffset())

			if err != nil {
				t.Error(test.Format(method, err.Error(), "nil"))
			} else {
				offset = length
				appendObjectResponse, err := bosClient.AppendObject(bucketName, objectKey, offset, str, nil, nil)

				if err != nil {
					t.Error(test.Format(method, err.Error(), "nil"))
				} else if appendObjectResponse.GetETag() == "" || appendObjectResponse.GetNextAppendOffset() == "" {
					t.Error(test.Format(method, "etag and next append offset are not exists", "etag and next append offset"))
				} else {
					length, err := strconv.Atoi(appendObjectResponse.GetNextAppendOffset())

					if err != nil {
						t.Error(test.Format(method, err.Error(), "nil"))
					} else if length != offset*2 {
						t.Error(test.Format(method, strconv.Itoa(length), strconv.Itoa(offset*2)))
					}
				}
			}
		}
	})
}

func TestMultipartUploadFromFile(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-multipart-upload-from-file-"
	method := "MultipartUploadFromFile"
	objectKey := "test-multipart-upload.zip"
	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		pwd, err := os.Getwd()

		if err != nil {
			t.Error(test.Format(method, err.Error(), "nil"))
			return
		}

		filePath := path.Join(pwd, "../", "examples", "test-multipart-upload.zip")
		var partSize int64 = 1024 * 1024 * 2

		completeMultipartUploadResponse, bceError := bosClient.MultipartUploadFromFile(bucketName,
			objectKey, filePath, partSize)

		if bceError != nil {
			t.Error(test.Format(method, bceError.Error(), "nil"))
		} else if completeMultipartUploadResponse.ETag == "" {
			t.Error(test.Format(method, "etag is not exists", "etag"))
		}
	})
}

func TestAbortMultipartUpload(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-abort-multipart-upload-"
	method := "AbortMultipartUpload"
	objectKey := "test-multipart-upload.zip"
	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		initiateMultipartUploadRequest := InitiateMultipartUploadRequest{
			BucketName: bucketName,
			ObjectKey:  objectKey,
		}

		initiateMultipartUploadResponse, err := bosClient.InitiateMultipartUpload(initiateMultipartUploadRequest, nil)

		if err != nil {
			t.Error(test.Format(method, err.Error(), "nil"))
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
			t.Error(test.Format(method, err.Error(), "nil"))
			return
		}
	})
}

func around(t *testing.T, method, bucketNamePrefix string, objectKey interface{}, f func(string)) {
	bucketName := bucketNamePrefix + strconv.Itoa(int(time.Now().Unix()))
	err := bosClient.CreateBucket(bucketName, nil)

	if err != nil {
		t.Error(test.Format(method+" at creating bucket", err.Error(), "nil"))
	} else {
		if f != nil {
			f(bucketName)

			if key, ok := objectKey.(string); ok {
				if key != "" {
					err = bosClient.DeleteObject(bucketName, key, nil)

					if err != nil {
						t.Error(test.Format(method+" at deleting object", err.Error(), "nil"))
					}
				}
			} else if keys, ok := objectKey.([]string); ok {
				if len(keys) > 0 {
					deleteMultipleObjectsResponse, err := bosClient.DeleteMultipleObjects(bucketName, keys, nil)

					if err != nil {
						t.Error(test.Format(method, err.Error(), "nil"))
					} else if deleteMultipleObjectsResponse != nil {
						str := ""

						for _, deleteMultipleObjectsError := range deleteMultipleObjectsResponse.Errors {
							str += deleteMultipleObjectsError.Error()
						}

						t.Error(test.Format(method, str, "empty string"))
					}
				}
			} else {
				t.Error(test.Format(method, "objectKey is not valid", "objectKey should be string or []string"))
			}
		}

		err = bosClient.DeleteBucket(bucketName, nil)

		if err != nil {
			t.Error(test.Format(method+" at deleting bucket", err.Error(), "nil"))
		}
	}
}
