package bos

import (
	bce "baidubce"
	"encoding/json"
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
	req, err := bce.NewRequest(
		"GET",
		"/v1/",
		c.Endpoint,
		nil,
		nil,
	)

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
