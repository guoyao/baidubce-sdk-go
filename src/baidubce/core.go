/**
 * Copyright (c) 2015 Guoyao Wu, All Rights Reserved
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * @file core.go
 * @author guoyao
 */

// Package baidubce define a set of core data structure, and implements a set of core functions
package baidubce

import (
	"baidubce/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	EXPIRATION_PERIOD_IN_SECONDS = 1800
)

var Region map[string]string = map[string]string{
	"bj": "bj.bcebos.com",
	"gz": "gz.bcebos.com",
}

type Credentials struct {
	AccessKeyId     string
	SecretAccessKey string
}

func NewCredentials(accessKeyId, secretAccessKey string) *Credentials {
	return &Credentials{accessKeyId, secretAccessKey}
}

var DefaultCredentials Credentials = Credentials{
	os.Getenv("BAIDU_BCE_AK"),
	os.Getenv("BAIDU_BCE_SK"),
}

type Config struct {
	Credentials
	Endpoint   string
	ApiVersion string
}

var DefaultConfig Config = Config{DefaultCredentials, "", "v1"}

type SignOption struct {
	Timestamp                 string
	ExpirationPeriodInSeconds int
	Headers                   map[string]string
	HeadersToSign             []string
	headersToSignSpecified    bool
}

func (option *SignOption) signedHeadersToString() string {
	var result string
	length := len(option.HeadersToSign)
	if option.headersToSignSpecified && length > 0 {
		headers := make([]string, 0, length)
		headers = append(headers, option.HeadersToSign...)
		sort.Strings(headers)
		result = strings.Join(headers, ";")
	}

	return result
}

func NewSignOption(timestamp string, expirationPeriodInSeconds int,
	headers map[string]string, headersToSign []string) *SignOption {

	return &SignOption{timestamp, expirationPeriodInSeconds,
		headers, headersToSign, len(headersToSign) > 0}
}

func GenerateAuthorization(credentials Credentials, req Request, option *SignOption) string {
	if option == nil {
		option = &SignOption{}
	}
	option.init()

	authorization := "bce-auth-v1/" + credentials.AccessKeyId
	authorization += "/" + option.Timestamp
	authorization += "/" + strconv.Itoa(option.ExpirationPeriodInSeconds)
	signature := sign(credentials, req, option)
	authorization += "/" + option.signedHeadersToString() + "/" + signature

	req.addHeader("Authorization", authorization)

	return authorization
}

type Client struct {
	Config
}

func (c *Client) GetBucketName(bucketName string) string {
	if c.Endpoint != "" && !util.MapContains(Region, func(key, value string) bool {
		return strings.ToLower(value) == strings.ToLower(c.Endpoint)
	}) {
		bucketName = ""
	}

	return bucketName
}

func (c *Client) SendRequest(req *Request, option *SignOption) ([]byte, error) {
	GenerateAuthorization(c.Credentials, *req, option)
	httpClient := http.Client{}
	res, err := httpClient.Do(req.raw())

	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return body, NewErrorFromJson(body)
	}

	return body, nil
}

func (option *SignOption) init() {
	if option.Timestamp == "" {
		option.Timestamp = util.TimeToUTCString(time.Now())
	}

	if option.ExpirationPeriodInSeconds <= 0 {
		option.ExpirationPeriodInSeconds = EXPIRATION_PERIOD_IN_SECONDS
	}

	if option.Headers == nil {
		option.Headers = make(map[string]string, 3)
	} else {
		util.MapKeyToLower(option.Headers)
	}

	option.headersToSignSpecified = len(option.HeadersToSign) > 0
	util.SliceToLower(option.HeadersToSign)

	if !util.Contains(option.HeadersToSign, "host", true) {
		option.HeadersToSign = append(option.HeadersToSign, "host")
	}

	if util.Contains(option.HeadersToSign, "date", true) {
		if !util.MapContains(option.Headers, generateHeaderValidCompareFunc("date")) {
			option.Headers["date"] = time.Now().Format(time.RFC1123)
		} else {
			option.Headers["date"] = util.TimeStringToRFC1123(util.GetMapValue(option.Headers, "date", true))
		}
	} else if util.Contains(option.HeadersToSign, "x-bce-date", true) {
		if !util.MapContains(option.Headers, generateHeaderValidCompareFunc("x-bce-date")) {
			option.Headers["x-bce-date"] = option.Timestamp
		}
	} else {
		option.HeadersToSign = append(option.HeadersToSign, "x-bce-date")
		option.Headers["x-bce-date"] = option.Timestamp
	}
}

func generateHeaderValidCompareFunc(headerKey string) func(string, string) bool {
	return func(key, value string) bool {
		return strings.ToLower(key) == strings.ToLower(headerKey) && value != ""
	}
}

// generate signature
func sign(credentials Credentials, req Request, option *SignOption) string {
	signingKey := getSigningKey(credentials, option)
	req.prepareHeaders(option)
	canonicalRequest := req.canonical(option)
	signature := util.HmacSha256Hex(signingKey, canonicalRequest)

	return signature
}

func getSigningKey(credentials Credentials, option *SignOption) string {
	var authStringPrefix = fmt.Sprintf("bce-auth-v1/%s", credentials.AccessKeyId)
	authStringPrefix += "/" + option.Timestamp
	authStringPrefix += "/" + strconv.Itoa(option.ExpirationPeriodInSeconds)

	return util.HmacSha256Hex(credentials.SecretAccessKey, authStringPrefix)
}
