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
