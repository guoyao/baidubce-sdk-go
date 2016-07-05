package bos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	bce "baidubce"
	"baidubce/util"
)

// Client is the client for bos.
type Client struct {
	bce.Client
}

// DefaultClient provided a default `bos.Client` instance.
var DefaultClient = Client{bce.Client{bce.DefaultConfig}}

// NewClient returns an instance of type `bos.Client`.
func NewClient(config bce.Config) Client {
	bceClient := bce.Client{config}
	return Client{bceClient}
}

func checkBucketName(bucketName string) {
	if bucketName == "" {
		panic("bucket name should not be empty.")
	}

	if strings.Index(bucketName, "/") == 0 {
		panic("bucket name should not be start with '/'")
	}
}

func checkObjectKey(objectKey string) {
	if objectKey == "" {
		panic("object key should not be empty.")
	}

	if strings.Index(objectKey, "/") == 0 {
		panic("object key should not be start with '/'")
	}
}

// GetBucketName returns the actual name of bucket.
func (c *Client) GetBucketName(bucketName string) string {
	if c.Endpoint != "" && !util.MapContains(bce.Region, func(key, value string) bool {
		return strings.ToLower(value) == strings.ToLower(c.Endpoint)
	}) {
		bucketName = ""
	}

	return bucketName
}

func (c *Client) GetBucketEndpoint(bucketName string) string {
	endpoint := c.Endpoint

	if endpoint == "" {
		endpoint = bce.Region["bj"]
	}

	if bucketName != "" {
		endpoint = bucketName + "." + endpoint
	}

	return endpoint
}

// GetBucketLocation returns the location of a bucket.
func (c *Client) GetBucketLocation(bucketName string, option *bce.SignOption) (*Location, *bce.Error) {
	bucketName = c.GetBucketName(bucketName)
	params := map[string]string{"location": ""}

	req, err := bce.NewRequest("GET", c.GetUriPath(""), c.GetBucketEndpoint(bucketName), params, nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	res, bceError := c.SendRequest(req, option, true)

	if bceError != nil {
		return nil, bceError
	}

	var location *Location
	err = json.Unmarshal(res.Body, &location)

	if err != nil {
		return nil, bce.NewError(err)
	}

	return location, nil
}

// ListBuckets is for getting a collection of bucket.
func (c *Client) ListBuckets(option *bce.SignOption) (*BucketSummary, *bce.Error) {
	req, err := bce.NewRequest("GET", c.GetUriPath(""), c.Endpoint, nil, nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	res, bceError := c.SendRequest(req, option, true)

	if bceError != nil {
		return nil, bceError
	}

	var bucketSummary *BucketSummary
	err = json.Unmarshal(res.Body, &bucketSummary)

	if err != nil {
		return nil, bce.NewError(err)
	}

	return bucketSummary, nil
}

// CreateBucket is for creating a bucket.
func (c *Client) CreateBucket(bucketName string, option *bce.SignOption) *bce.Error {
	req, err := bce.NewRequest("PUT", c.GetUriPath(""), c.GetBucketEndpoint(bucketName), nil, nil)

	if err != nil {
		return bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")

	_, bceError := c.SendRequest(req, option, true)

	return bceError
}

func (c *Client) DoesBucketExist(bucketName string, option *bce.SignOption) (bool, *bce.Error) {
	req, err := bce.NewRequest("HEAD", c.GetUriPath(""), c.GetBucketEndpoint(bucketName), nil, nil)

	if err != nil {
		return false, bce.NewError(err)
	}

	res, bceError := c.SendRequest(req, option, true)

	if res != nil {
		switch {
		case res.StatusCode < 400 || res.StatusCode == 403:
			return true, nil
		case res.StatusCode == 404:
			return false, nil
		}
	}

	return false, bceError
}

func (c *Client) DeleteBucket(bucketName string, option *bce.SignOption) *bce.Error {
	req, err := bce.NewRequest("DELETE", c.GetUriPath(""), c.GetBucketEndpoint(bucketName), nil, nil)

	if err != nil {
		return bce.NewError(err)
	}

	_, bceError := c.SendRequest(req, option, true)

	return bceError
}

func (c *Client) SetBucketPrivate(bucketName string, option *bce.SignOption) *bce.Error {
	return c.setBucketAclFromString(bucketName, CannedAccessControlList["Private"], option)
}

func (c *Client) SetBucketPublicRead(bucketName string, option *bce.SignOption) *bce.Error {
	return c.setBucketAclFromString(bucketName, CannedAccessControlList["PublicRead"], option)
}

func (c *Client) SetBucketPublicReadWrite(bucketName string, option *bce.SignOption) *bce.Error {
	return c.setBucketAclFromString(bucketName, CannedAccessControlList["PublicReadWrite"], option)
}

func (c *Client) GetBucketAcl(bucketName string, option *bce.SignOption) (*BucketAcl, *bce.Error) {
	params := map[string]string{"acl": ""}
	req, err := bce.NewRequest("GET", c.GetUriPath(""), c.GetBucketEndpoint(bucketName), params, nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	res, bceError := c.SendRequest(req, option, true)

	if bceError != nil {
		return nil, bceError
	}

	var bucketAcl *BucketAcl
	err = json.Unmarshal(res.Body, &bucketAcl)

	if err != nil {
		return nil, bce.NewError(err)
	}

	return bucketAcl, nil
}

func (c *Client) SetBucketAcl(bucketName string, bucketAcl BucketAcl, option *bce.SignOption) *bce.Error {
	byteArray, err := util.ToJson(bucketAcl, "accessControlList")

	if err != nil {
		return bce.NewError(err)
	}

	params := map[string]string{"acl": ""}
	req, err := bce.NewRequest("PUT", c.GetUriPath(""), c.GetBucketEndpoint(bucketName), params, bytes.NewReader(byteArray))

	if err != nil {
		return bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")

	_, bceError := c.SendRequest(req, option, true)

	return bceError
}

func (c *Client) PutObject(bucketName, objectKey string, data interface{}, metadata *ObjectMetadata, option *bce.SignOption) (PutObjectResponse, *bce.Error) {
	checkObjectKey(objectKey)

	var reader io.Reader

	if str, ok := data.(string); ok {
		reader = strings.NewReader(str)
	} else if byteArray, ok := data.([]byte); ok {
		reader = bytes.NewReader(byteArray)
	} else if r, ok := data.(io.Reader); ok {
		byteArray, err := ioutil.ReadAll(r)

		if err != nil {
			return nil, bce.NewError(err)
		}

		reader = bytes.NewReader(byteArray)
	} else {
		panic("data type should be string or []byte or io.Reader.")
	}

	req, err := bce.NewRequest("PUT", c.GetUriPath(objectKey), c.GetBucketEndpoint(bucketName), nil, reader)

	if err != nil {
		return nil, bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")

	if metadata != nil {
		metadata.MergeToSignOption(option)
	}

	res, bceError := c.SendRequest(req, option, true)

	if bceError != nil {
		return nil, bceError
	}

	putObjectResponse := NewPutObjectResponse(res.Header)

	return putObjectResponse, nil
}

func (c *Client) DeleteObject(bucketName, objectKey string, option *bce.SignOption) *bce.Error {
	checkObjectKey(objectKey)

	req, err := bce.NewRequest("DELETE", c.GetUriPath(objectKey), c.GetBucketEndpoint(bucketName), nil, nil)

	if err != nil {
		return bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")

	_, bceError := c.SendRequest(req, option, true)

	if bceError != nil {
		return bceError
	}

	return nil
}

func (c *Client) ListObjects(bucketName string, params map[string]string, option *bce.SignOption) (*ListObjectsResponse, *bce.Error) {
	req, err := bce.NewRequest("GET", c.GetUriPath(""), c.GetBucketEndpoint(bucketName), params, nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")

	res, bceError := c.SendRequest(req, option, true)

	if bceError != nil {
		return nil, bceError
	}

	var listObjectsResponse *ListObjectsResponse
	err = json.Unmarshal(res.Body, &listObjectsResponse)

	if err != nil {
		return nil, bce.NewError(err)
	}

	return listObjectsResponse, nil
}

func (c *Client) CopyObject(srcBucketName, srcKey, destBucketName, destKey string, option *bce.SignOption) (*CopyObjectResponse, *bce.Error) {
	checkBucketName(srcBucketName)
	checkBucketName(destBucketName)
	checkObjectKey(srcKey)
	checkObjectKey(destKey)

	req, err := bce.NewRequest("PUT", c.GetUriPath(destKey), c.GetBucketEndpoint(destBucketName), nil, nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")
	option.AddHeader("x-bce-copy-source", fmt.Sprintf("/%s/%s", srcBucketName, srcKey))

	res, bceError := c.SendRequest(req, option, true)

	if bceError != nil {
		return nil, bceError
	}

	var copyObjectResponse *CopyObjectResponse
	err = json.Unmarshal(res.Body, &copyObjectResponse)

	if err != nil {
		return nil, bce.NewError(err)
	}

	return copyObjectResponse, nil
}

func (c *Client) CopyObjectFromRequest(copyObjectRequest *CopyObjectRequest, option *bce.SignOption) (*CopyObjectResponse, *bce.Error) {
	checkBucketName(copyObjectRequest.SrcBucketName)
	checkBucketName(copyObjectRequest.DestBucketName)
	checkObjectKey(copyObjectRequest.SrcKey)
	checkObjectKey(copyObjectRequest.DestKey)

	req, err := bce.NewRequest("PUT", c.GetUriPath(copyObjectRequest.DestKey), c.GetBucketEndpoint(copyObjectRequest.DestBucketName), nil, nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")

	source := util.URIEncodeExceptSlash(fmt.Sprintf("/%s/%s", copyObjectRequest.SrcBucketName, copyObjectRequest.SrcKey))
	option.AddHeader("x-bce-copy-source", source)
	copyObjectRequest.MergeToSignOption(option)

	res, bceError := c.SendRequest(req, option, true)

	if bceError != nil {
		return nil, bceError
	}

	var copyObjectResponse *CopyObjectResponse
	err = json.Unmarshal(res.Body, &copyObjectResponse)

	if err != nil {
		return nil, bce.NewError(err)
	}

	return copyObjectResponse, nil
}

func (c *Client) GetObject(bucketName, objectKey string, option *bce.SignOption) (*Object, *bce.Error) {
	checkBucketName(bucketName)
	checkObjectKey(objectKey)

	req, err := bce.NewRequest("GET", c.GetUriPath(objectKey), c.GetBucketEndpoint(bucketName), nil, nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")

	res, bceError := c.SendRequest(req, option, false)

	if bceError != nil {
		return nil, bceError
	}

	object := &Object{
		ObjectMetadata: NewObjectMetadataFromHeader(res.Header),
		ObjectContent:  res.Response.Body,
	}

	return object, nil
}

func (c *Client) GetObjectFromRequest(getObjectRequest *GetObjectRequest, option *bce.SignOption) (*Object, *bce.Error) {
	checkBucketName(getObjectRequest.BucketName)
	checkObjectKey(getObjectRequest.ObjectKey)

	req, err := bce.NewRequest(
		"GET",
		c.GetUriPath(getObjectRequest.ObjectKey),
		c.GetBucketEndpoint(getObjectRequest.BucketName),
		nil,
		nil,
	)

	if err != nil {
		return nil, bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")

	getObjectRequest.MergeToSignOption(option)

	res, bceError := c.SendRequest(req, option, false)

	if bceError != nil {
		return nil, bceError
	}

	object := &Object{
		ObjectMetadata: NewObjectMetadataFromHeader(res.Header),
		ObjectContent:  res.Response.Body,
	}

	return object, nil
}

func (c *Client) GetObjectToFile(getObjectRequest *GetObjectRequest, file *os.File, option *bce.SignOption) (*ObjectMetadata, *bce.Error) {
	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	checkBucketName(getObjectRequest.BucketName)
	checkObjectKey(getObjectRequest.ObjectKey)

	req, err := bce.NewRequest(
		"GET",
		c.GetUriPath(getObjectRequest.ObjectKey),
		c.GetBucketEndpoint(getObjectRequest.BucketName),
		nil,
		nil,
	)

	if err != nil {
		return nil, bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")

	getObjectRequest.MergeToSignOption(option)

	res, bceError := c.SendRequest(req, option, true)

	if bceError != nil {
		return nil, bceError
	}

	objectMetadata := NewObjectMetadataFromHeader(res.Header)

	_, err = file.Write(res.Body)

	if err != nil {
		return objectMetadata, bce.NewError(err)
	}

	return objectMetadata, nil
}

func (c *Client) GetObjectMetadata(bucketName, objectKey string, option *bce.SignOption) (*ObjectMetadata, *bce.Error) {
	checkBucketName(bucketName)
	checkObjectKey(objectKey)

	req, err := bce.NewRequest("HEAD", c.GetUriPath(objectKey), c.GetBucketEndpoint(bucketName), nil, nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	res, bceError := c.SendRequest(req, option, false)

	if bceError != nil {
		return nil, bceError
	}

	objectMetadata := NewObjectMetadataFromHeader(res.Header)

	return objectMetadata, nil
}

func (c *Client) setBucketAclFromString(bucketName, acl string, option *bce.SignOption) *bce.Error {
	params := map[string]string{"acl": ""}
	req, err := bce.NewRequest("PUT", c.GetUriPath(""), c.GetBucketEndpoint(bucketName), params, nil)

	if err != nil {
		return bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")

	headers := map[string]string{"x-bce-acl": acl}
	option.AddHeaders(headers)

	_, bceError := c.SendRequest(req, option, true)

	return bceError
}
