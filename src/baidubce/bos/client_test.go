package bos

import (
	"baidubce/test"
	"testing"
)

var bosClient Client = DefaultClient

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
