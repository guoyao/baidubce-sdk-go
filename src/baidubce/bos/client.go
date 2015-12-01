package bos

import (
	bce "baidubce"
	"encoding/json"
	"fmt"
)

type Client struct {
	bce.Client
}

var DefaultClient Client = Client{bce.Client{bce.DefaultConfig}}

func NewClient(config bce.Config) Client {
	bceClient := bce.Client{config}
	return Client{bceClient}
}

func (c *Client) GetBucketLocation(bucketName string, option *bce.SignOption) (*Location, error) {
	bucketName = c.GetBucketName(bucketName)

	req, err := bce.NewRequest(
		"GET",
		"/"+bucketName,
		c.Endpoint,
		map[string]string{"location": ""},
		nil,
	)

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

func (c *Client) ListBuckets(option *bce.SignOption) (*BucketSummary, error) {
	req, err := bce.NewRequest("GET", fmt.Sprintf("/%s/", c.ApiVersion), c.Endpoint, nil, nil)

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

func (c *Client) CreateBucket(bucketName string, option *bce.SignOption) error {
	option = &bce.SignOption{
		HeadersToSign: []string{"date"},
	}

	req, err := bce.NewRequest("PUT", fmt.Sprintf("/%s/%s", c.ApiVersion, bucketName), c.Endpoint, nil, nil)

	if err != nil {
		return err
	}

	_, err = c.SendRequest(req, option)

	if err != nil {
		return err
	}

	return nil
}
