package bos

import (
	"time"
)

// Location is a struct for bucket location info.
type Location struct {
	LocationConstraint string
}

type BucketOwner struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName,omitempty"`
}

// Bucket is a struct for bucket info.
type Bucket struct {
	Name, Location string
	CreationDate   time.Time
}

// BucketSummary is a struct for bucket summary.
type BucketSummary struct {
	Owner   BucketOwner
	Buckets []Bucket
}

type BucketAcl struct {
	Owner             BucketOwner `json:"-"`
	AccessControlList []Grant     `json:"accessControlList"`
}

type Grant struct {
	Grantee    []BucketGrantee `json:"grantee"`
	Permission []string        `json:"permission"`
}

type BucketGrantee struct {
	Id string `json:"id"`
}

var CannedAccessControlList = map[string]string{
	"Private":         "private",
	"PublicRead":      "public-read",
	"PublicReadWrite": "public-read-write",
}
