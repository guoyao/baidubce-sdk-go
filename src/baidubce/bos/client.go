package bos

import (
	bce "baidubce"
	"io/ioutil"
	"log"
	"net/http"
)

type Client struct {
	bce.Config
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

	return c.sendRequest(req, option)
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

	return c.sendRequest(req, option)
}

func (c *Client) sendRequest(req *bce.Request, option *bce.SignOption) (string, error) {
	if option == nil {
		option = bce.NewSignOption("", bce.EXPIRATION_PERIOD_IN_SECONDS)
	}

	authorization := bce.GenerateAuthorization(c.Credentials, *req, option)
	req.Header.Add("Authorization", authorization)
	httpClient := http.Client{}
	res, err := httpClient.Do((*http.Request)(req))

	defer res.Body.Close()

	if err != nil {
		log.Println(err)
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(body), nil
}
