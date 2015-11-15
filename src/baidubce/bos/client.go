package bos

import (
	bce "baidubce"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Config struct {
	bce.Credentials
	Endpoint string
}

type Client struct {
	Config
}

func (c *Client) GetBucketLocation(bucketName string, option *bce.SignOption) (string, error) {
	bceRequest := bce.Request{
		HttpMethod: "GET",
		URI:        "/" + bucketName,
		Params:     map[string]string{"location": ""},
	}

	return c.sendRequest(bceRequest, option)
}

func (c *Client) ListBucket(option *bce.SignOption) (string, error) {
	bceRequest := bce.Request{
		HttpMethod: "GET",
		URI:        "/v1/",
	}

	return c.sendRequest(bceRequest, option)
}

func (c *Client) sendRequest(bceRequest bce.Request, option *bce.SignOption) (string, error) {
	if option == nil {
		option = bce.NewSignOption("", bce.EXPIRATION_PERIOD_IN_SECONDS)
	}

	bceRequest.Header = getHttpHeader()
	authorization := bce.GenerateAuthorization(c.Credentials, bceRequest, option)
	URI := bce.Region["bj"]

	if c.Endpoint != "" {
		URI = c.Endpoint
	}

	URI += bceRequest.URI
	queryString := bceRequest.ParamsToCanonicalQueryString()
	if queryString != "" {
		URI += "?" + queryString
	}

	req, err := http.NewRequest(bceRequest.HttpMethod, URI, nil)
	if err != nil {
		log.Println(err)
		return "", err
	}

	req.Header = bceRequest.Header
	req.Header.Add("Authorization", authorization)
	httpClient := http.Client{}
	res, err := httpClient.Do(req)

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

func getHttpHeader() http.Header {
	var header http.Header = http.Header{}

	header.Add("Host", "bj.bcebos.com")
	header.Add("Date", time.Now().Format(time.RFC1123))

	return header
}
