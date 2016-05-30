package bos

import (
	"bytes"
	"encoding/json"
	"fmt"
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

// GetBucketName returns the actual name of bucket.
func (c *Client) GetBucketName(bucketName string) string {
	if c.Endpoint != "" && !util.MapContains(bce.Region, func(key, value string) bool {
		return strings.ToLower(value) == strings.ToLower(c.Endpoint)
	}) {
		bucketName = ""
	}

	return bucketName
}

// GetBucketLocation returns the location of a bucket.
func (c *Client) GetBucketLocation(bucketName string, option *bce.SignOption) (*Location, *bce.Error) {
	bucketName = c.GetBucketName(bucketName)

	req, err := bce.NewRequest("GET", "/"+bucketName, c.Endpoint, map[string]string{"location": ""}, nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	res, bceError := c.SendRequest(req, option)

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
	req, err := bce.NewRequest("GET", fmt.Sprintf("/%s/", c.APIVersion), c.Endpoint, nil, nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	res, bceError := c.SendRequest(req, option)

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
	option = bce.AddDateToSignOption(option)
	req, err := bce.NewRequest("PUT", fmt.Sprintf("/%s/%s", c.APIVersion, bucketName), c.Endpoint, nil, nil)

	if err != nil {
		return bce.NewError(err)
	}

	_, bceError := c.SendRequest(req, option)

	return bceError
}

func (c *Client) DoesBucketExist(bucketName string, option *bce.SignOption) (bool, *bce.Error) {
	req, err := bce.NewRequest("HEAD", fmt.Sprintf("/%s/%s", c.APIVersion, bucketName), c.Endpoint, nil, nil)

	if err != nil {
		return false, bce.NewError(err)
	}

	res, bceError := c.SendRequest(req, option)

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
	req, err := bce.NewRequest("DELETE", fmt.Sprintf("/%s/%s", c.APIVersion, bucketName), c.Endpoint, nil, nil)

	if err != nil {
		return bce.NewError(err)
	}

	_, bceError := c.SendRequest(req, option)

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
	req, err := bce.NewRequest("GET", fmt.Sprintf("/%s/%s", c.APIVersion, bucketName), c.Endpoint, params, nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	res, bceError := c.SendRequest(req, option)

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
	option = bce.AddDateToSignOption(option)
	byteArray, err := json.Marshal(bucketAcl)

	if err != nil {
		return bce.NewError(err)
	}

	m, err := util.JsonToMap(byteArray, "accessControlList")

	if err != nil {
		return bce.NewError(err)
	}

	byteArray, err = json.Marshal(m)

	if err != nil {
		return bce.NewError(err)
	}

	params := map[string]string{"acl": ""}
	req, err := bce.NewRequest("PUT", fmt.Sprintf("/%s/%s", c.APIVersion, bucketName), c.Endpoint, params, bytes.NewReader(byteArray))

	if err != nil {
		return bce.NewError(err)
	}

	_, bceError := c.SendRequest(req, option)

	return bceError
}

func (c *Client) setBucketAclFromString(bucketName, acl string, option *bce.SignOption) *bce.Error {
	option = bce.AddDateToSignOption(option)
	params := map[string]string{"acl": ""}
	req, err := bce.NewRequest("PUT", fmt.Sprintf("/%s/%s", c.APIVersion, bucketName), c.Endpoint, params, nil)

	if err != nil {
		return bce.NewError(err)
	}

	headers := map[string]string{"x-bce-acl": acl}
	req.AddHeaders(headers)
	_, bceError := c.SendRequest(req, option)

	return bceError
}
