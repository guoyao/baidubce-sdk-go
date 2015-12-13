package bos

import (
	"encoding/json"
	"fmt"

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

// GetBucketLocation returns the location of a bucket.
func (c *Client) GetBucketLocation(bucketName string, option *bce.SignOption) (*Location, error) {
	bucketName = c.GetBucketName(bucketName)

	req, err := bce.NewRequest("GET", "/"+bucketName, c.Endpoint, map[string]string{"location": ""}, nil)

	if err != nil {
		return nil, err
	}

	respBody, err := c.SendRequest(req, option)

	if err != nil {
		return nil, err
	}

	var location *Location
	err = json.Unmarshal(respBody, &location)

	if err != nil {
		return nil, err
	}

	return location, nil
}

// ListBuckets is for getting a collection of bucket.
func (c *Client) ListBuckets(option *bce.SignOption) (*BucketSummary, error) {
	req, err := bce.NewRequest("GET", fmt.Sprintf("/%s/", c.APIVersion), c.Endpoint, nil, nil)

	if err != nil {
		return nil, err
	}

	respBody, err := c.SendRequest(req, option)

	if err != nil {
		return nil, err
	}

	var bucketSummary *BucketSummary
	err = json.Unmarshal(respBody, &bucketSummary)

	if err != nil {
		return nil, err
	}

	return bucketSummary, nil
}

// CreateBucket is for creating a bucket.
func (c *Client) CreateBucket(bucketName string, option *bce.SignOption) error {
	if option == nil {
		option = &bce.SignOption{
			HeadersToSign: []string{"date"},
		}
	} else if option.HeadersToSign == nil {
		option.HeadersToSign = []string{"date"}
	} else if !util.Contains(option.HeadersToSign, "date", true) {
		option.HeadersToSign = append(option.HeadersToSign, "date")
	}

	req, err := bce.NewRequest("PUT", fmt.Sprintf("/%s/%s", c.APIVersion, bucketName), c.Endpoint, nil, nil)

	if err != nil {
		return err
	}

	_, err = c.SendRequest(req, option)

	if err != nil {
		return err
	}

	return nil
}
