package bos

import (
	bce "baidubce"
)

type Client struct {
	bce.Client
}

func NewClient(config bce.Config) Client {
	bceClient := bce.Client{config}
	return Client{bceClient}
}

func (c *Client) GetBucketLocation(bucketName string, option *bce.SignOption) (string, error) {
	if c.Endpoint != "" {
		bucketName = ""
	}

	req, err := bce.NewRequest(
		"GET",
		"/"+bucketName,
		c.Endpoint,
		map[string]string{"location": ""},
		nil,
	)

	if err != nil {
		return "", err
	}

	return c.SendRequest(req, option)
}

func (c *Client) ListBucket(option *bce.SignOption) (string, error) {
	req, err := bce.NewRequest(
		"GET",
		"/v1/",
		c.Endpoint,
		nil,
		nil,
	)

	if err != nil {
		return "", err
	}

	return c.SendRequest(req, option)
}
