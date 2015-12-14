package bos

import (
	"strconv"
	"testing"
	"time"

	"baidubce/test"
)

var bosClient = DefaultClient

func TestGetBucketLocation(t *testing.T) {
	expected := "bj"
	location, _ := bosClient.GetBucketLocation("baidubce-sdk-go", nil)
	if location.LocationConstraint != expected {
		t.Error(test.Format("GetBucketLocation", location.LocationConstraint, expected))
	}
}

func TestListBuckets(t *testing.T) {
	expected := "baidubce-sdk-go"
	bucketSummary, _ := bosClient.ListBuckets(nil)
	bucket := bucketSummary.Buckets[0]
	if bucket.Name != expected {
		t.Error(test.Format("ListBuckets", bucket.Name, expected))
	}
}

func TestCreateBucket(t *testing.T) {
	bucketName := "baidubce-sdk-go-create-bucket-test-" + strconv.Itoa(int(time.Now().Unix()))
	err := bosClient.CreateBucket(bucketName, nil)

	if err != nil {
		t.Error(test.Format("CreateBucket", err.Error(), "nil"))
	}
}

func TestDoesBucketExist(t *testing.T) {
	expected := true
	bucketName := "baidubce-sdk-go"
	exists, err := bosClient.DoesBucketExist(bucketName, nil)

	if err != nil {
		t.Error(test.Format("DoesBucketExist", err.Error(), strconv.FormatBool(expected)))
	} else if exists != expected {
		t.Error(test.Format("DoesBucketExist", strconv.FormatBool(exists), strconv.FormatBool(expected)))
	}
}
