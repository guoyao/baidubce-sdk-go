package bos

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	bce "baidubce"
	"baidubce/util"
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
	Expires            string
	ContentSha256      string

	ContentRange string
	ETag         string
	UserMetadata map[string]string
}

func NewObjectMetadataFromHeader(h http.Header) *ObjectMetadata {
	objectMetadata := &ObjectMetadata{}

	for key, _ := range h {
		key = strings.ToLower(key)

		if key == "cache-control" {
			objectMetadata.CacheControl = h.Get(key)
		} else if key == "content-disposition" {
			objectMetadata.ContentDisposition = h.Get(key)
		} else if key == "content-length" {
			length, err := strconv.Atoi(h.Get(key))

			if err == nil {
				objectMetadata.ContentLength = length
			}
		} else if key == "content-range" {
			objectMetadata.ContentRange = h.Get(key)
		} else if key == "content-type" {
			objectMetadata.ContentType = h.Get(key)
		} else if key == "expires" {
			objectMetadata.Expires = h.Get(key)
		} else if key == "etag" {
			objectMetadata.ETag = h.Get(key)
		} else if IsUserDefinedMetadata(key) {
			objectMetadata.UserMetadata[key] = h.Get(key)
		}
	}

	return objectMetadata
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

	if metadata.Expires != "" {
		option.AddHeader("Expires", metadata.Expires)
	}

	if metadata.ContentSha256 != "" {
		option.AddHeader("x-bce-content-sha256", metadata.ContentSha256)
	}

	for key, value := range metadata.UserMetadata {
		option.AddHeader(ToUserDefinedMetadata(key), value)
	}
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
	Name           string
	Prefix         string
	Delimiter      string
	Marker         string
	NextMarker     string
	MaxKeys        uint
	IsTruncated    bool
	Contents       []ObjectSummary
	CommonPrefixes []map[string]string
}

func (listObjectsResponse *ListObjectsResponse) GetCommonPrefixes() []string {
	prefixes := make([]string, 0, len(listObjectsResponse.CommonPrefixes))

	for _, commonPrefix := range listObjectsResponse.CommonPrefixes {
		prefixes = append(prefixes, commonPrefix["prefix"])
	}

	return prefixes
}

type CopyObjectResponse struct {
	ETag         string
	LastModified time.Time
}

type CopyObjectRequest struct {
	SrcBucketName         string          `json:"-"`
	SrcKey                string          `json:"-"`
	DestBucketName        string          `json:"-"`
	DestKey               string          `json:"-"`
	ObjectMetadata        *ObjectMetadata `json:"-"`
	SourceMatch           string          `json:"x-bce-copy-source-if-match,omitempty"`
	SourceNoneMatch       string          `json:"x-bce-copy-source-if-none-match,omitempty"`
	SourceModifiedSince   string          `json:"x-bce-copy-source-if-modified-since,omitempty"`
	SourceUnmodifiedSince string          `json:"x-bce-copy-source-if-unmodified-since,omitempty"`
}

func (copyObjectRequest *CopyObjectRequest) MergeToSignOption(option *bce.SignOption) {
	m, err := util.ToMap(copyObjectRequest)

	if err != nil {
		return
	}

	headerMap := make(map[string]string)

	for key, value := range m {
		if str, ok := value.(string); ok {
			headerMap[key] = str
		}
	}

	option.AddHeaders(headerMap)

	if copyObjectRequest.ObjectMetadata != nil {
		option.AddHeader("x-bce-metadata-directive", "replace")
		copyObjectRequest.ObjectMetadata.MergeToSignOption(option)
	} else {
		option.AddHeader("x-bce-metadata-directive", "copy")
	}
}

type Object struct {
	ObjectMetadata *ObjectMetadata
	ObjectContent  io.ReadCloser
}

type GetObjectRequest struct {
	BucketName string
	ObjectKey  string
	Range      string
}

func (getObjectRequest *GetObjectRequest) MergeToSignOption(option *bce.SignOption) {
	if getObjectRequest.Range != "" {
		option.AddHeader("Range", "bytes="+getObjectRequest.Range)
	}
}

func (getObjectRequest *GetObjectRequest) SetRange(start uint, end uint) {
	getObjectRequest.Range = fmt.Sprintf("%v-%v", start, end)
}

var UserDefinedMetadataPrefix = "x-bce-meta-"

var CannedAccessControlList = map[string]string{
	"Private":         "private",
	"PublicRead":      "public-read",
	"PublicReadWrite": "public-read-write",
}

func IsUserDefinedMetadata(metadata string) bool {
	return strings.Index(metadata, UserDefinedMetadataPrefix) == 0
}

func ToUserDefinedMetadata(metadata string) string {
	if IsUserDefinedMetadata(metadata) {
		return metadata
	}

	return UserDefinedMetadataPrefix + metadata
}
