package bos

import (
	"io/ioutil"
	"os"
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

func around(t *testing.T, method, bucketNamePrefix, objectKey string, f func(string)) {
	bucketName := bucketNamePrefix + strconv.Itoa(int(time.Now().Unix()))
	err := bosClient.CreateBucket(bucketName, nil)

	if err != nil {
		t.Error(test.Format(method+" at creating bucket", err.Error(), "nil"))
	} else {
		if f != nil {
			f(bucketName)

			if objectKey != "" {
				err = bosClient.DeleteObject(bucketName, objectKey, nil)

				if err != nil {
					t.Error(test.Format(method+" at deleting object", err.Error(), "nil"))
				}
			}
		}

		err = bosClient.DeleteBucket(bucketName, nil)
		if err != nil {
			t.Error(test.Format(method+" at deleting bucket", err.Error(), "nil"))
		}
	}
}
