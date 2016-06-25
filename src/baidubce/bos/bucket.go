package bos

import (
	"net/http"
	"strconv"
	"time"

	bce "baidubce"
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
	ContentLength      int
	ContentMD5         string
	ContentType        string
	ContentRange       string
	Expires            string
	ETag               string
	UserMetadata       map[string]string
}

func (metadata *ObjectMetadata) AddUserMetadata(key, value string) {
	if metadata.UserMetadata == nil {
		metadata.UserMetadata = make(map[string]string)
	}

	metadata.UserMetadata[key] = value
}

func (metadata *ObjectMetadata) MergeToSignOption(option *bce.SignOption) {
	if metadata.CacheControl != "" {
		option.AddHeader("Cache-Control", metadata.CacheControl)
	}

	if metadata.ContentDisposition != "" {
		option.AddHeader("Content-Disposition", metadata.ContentDisposition)
	}

	if metadata.ContentLength != 0 {
		option.AddHeader("Content-Length", strconv.Itoa(metadata.ContentLength))
	}

	if metadata.ContentMD5 != "" {
		option.AddHeader("Content-MD5", metadata.ContentMD5)
	}

	if metadata.ContentType != "" {
		option.AddHeader("Content-Type", metadata.ContentType)
	}

	if metadata.ContentRange != "" {
		option.AddHeader("Content-Range", metadata.ContentRange)
	}

	if metadata.Expires != "" {
		option.AddHeader("Expires", metadata.Expires)
	}

	if metadata.ETag != "" {
		option.AddHeader("ETag", metadata.ETag)
	}

	option.AddHeaders(metadata.UserMetadata)
}

type PutObjectResponse http.Header

func NewPutObjectResponse(h http.Header) PutObjectResponse {
	return PutObjectResponse(h)
}

func (res PutObjectResponse) Get(key string) string {
	return http.Header(res).Get(key)
}

func (res PutObjectResponse) GetETag() string {
	return res.Get("Etag")
}

type ObjectSummary struct {
	Key          string
	LastModified string
	ETag         string
	Size         uint
	Owner        BucketOwner
}

type ListObjectsResponse struct {
	Name        string
	Prefix      string
	Delimiter   string
	Marker      string
	MaxKeys     uint
	IsTruncated bool
	Contents    []ObjectSummary
}

var CannedAccessControlList = map[string]string{
	"Private":         "private",
	"PublicRead":      "public-read",
	"PublicReadWrite": "public-read-write",
}
