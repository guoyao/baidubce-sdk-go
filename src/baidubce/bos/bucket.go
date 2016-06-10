package bos

import (
	"net/http"
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
	Owner             BucketOwner `json:"owner"`
	AccessControlList []Grant     `json:"accessControlList"`
}

type Grant struct {
	Grantee    []BucketGrantee `json:"grantee"`
	Permission []string        `json:"permission"`
}

type BucketGrantee struct {
	Id string `json:"id"`
}

type ObjectMetadata struct {
	CacheControl       string
	ContentDisposition string
	ContentLength      uint
	ContentMD5         string
	ContentType        string
	Expires            string
	UserMetadata       map[string]string
}

type PutObjectResponse http.Header

func NewPutObjectResponse(h http.Header) PutObjectResponse {
	return PutObjectResponse(h)
}

func (res PutObjectResponse) GetETag() string {
	return res["Etag"][0]
}

var CannedAccessControlList = map[string]string{
	"Private":         "private",
	"PublicRead":      "public-read",
	"PublicReadWrite": "public-read-write",
}
