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

// Package baidubce define a set of core data structure, and implements a set of core functions.
package baidubce

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"baidubce/util"
)

const (
	// ExpirationPeriodInSeconds 1800s is the default expiration period.
	ExpirationPeriodInSeconds = 1800
)

// Region map of baidubce
var Region = map[string]string{
	"bj": "bj.bcebos.com",
	"gz": "gz.bcebos.com",
}

// Credentials struct for baidubce.
type Credentials struct {
	AccessKeyID     string
	SecretAccessKey string
}

// NewCredentials returns an instance of type `Credentials`.
func NewCredentials(AccessKeyID, secretAccessKey string) *Credentials {
	return &Credentials{AccessKeyID, secretAccessKey}
}

// DefaultCredentials provided a default `Credentials` instance.
var DefaultCredentials = Credentials{
	os.Getenv("BAIDU_BCE_AK"),
	os.Getenv("BAIDU_BCE_SK"),
}

// Config contains options for baidubce api.
type Config struct {
	Credentials
	Endpoint   string
	APIVersion string
}

// DefaultConfig provided a default `Config` instance.
var DefaultConfig = Config{DefaultCredentials, "", "v1"}

// SignOption contains options for signature of baidubce api.
type SignOption struct {
	Timestamp                 string
	ExpirationPeriodInSeconds int
	Headers                   map[string]string
	HeadersToSign             []string
	headersToSignSpecified    bool
}

// NewSignOption is the instance factory for `SignOption`.
func NewSignOption(timestamp string, expirationPeriodInSeconds int,
	headers map[string]string, headersToSign []string) *SignOption {

	return &SignOption{timestamp, expirationPeriodInSeconds,
		headers, headersToSign, len(headersToSign) > 0}
}

func CheckSignOption(option *SignOption) *SignOption {
	if option == nil {
		return &SignOption{}
	}

	return option
}

func (option *SignOption) AddHeadersToSign(headers ...string) {
	if option.HeadersToSign == nil {
		option.HeadersToSign = []string{}
		option.HeadersToSign = append(option.HeadersToSign, headers...)
	} else {
		for _, header := range headers {
			if !util.Contains(option.HeadersToSign, header, true) {
				option.HeadersToSign = append(option.HeadersToSign, header)
			}
		}
	}
}

func (option *SignOption) AddHeader(key, value string) {
	if option.Headers == nil {
		option.Headers = make(map[string]string)
		option.Headers[key] = value
	}

	if !util.MapContains(option.Headers, generateHeaderValidCompareFunc(key)) {
		option.Headers[key] = value
	}
}

func (option *SignOption) AddHeaders(headers map[string]string) {
	if headers == nil {
		return
	}

	if option.Headers == nil {
		option.Headers = make(map[string]string)
	}

	for key, value := range headers {
		option.AddHeader(key, value)
	}
}

func (option *SignOption) init() {
	if option.Timestamp == "" {
		option.Timestamp = util.TimeToUTCString(time.Now())
	}

	if option.ExpirationPeriodInSeconds <= 0 {
		option.ExpirationPeriodInSeconds = ExpirationPeriodInSeconds
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

	if !option.headersToSignSpecified {
		option.HeadersToSign = append(option.HeadersToSign, "x-bce-date")
		option.Headers["x-bce-date"] = option.Timestamp
	} else if util.Contains(option.HeadersToSign, "date", true) {
		if !util.MapContains(option.Headers, generateHeaderValidCompareFunc("date")) {
			option.Headers["date"] = time.Now().Format(time.RFC1123)
		} else {
			option.Headers["date"] = util.TimeStringToRFC1123(util.GetMapValue(option.Headers, "date", true))
		}
	} else if util.Contains(option.HeadersToSign, "x-bce-date", true) {
		if !util.MapContains(option.Headers, generateHeaderValidCompareFunc("x-bce-date")) {
			option.Headers["x-bce-date"] = option.Timestamp
		}
	}
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

// GenerateAuthorization returns authorization code of baidubce api.
func GenerateAuthorization(credentials Credentials, req Request, option *SignOption) string {
	if option == nil {
		option = &SignOption{}
	}
	option.init()

	authorization := "bce-auth-v1/" + credentials.AccessKeyID
	authorization += "/" + option.Timestamp
	authorization += "/" + strconv.Itoa(option.ExpirationPeriodInSeconds)
	signature := sign(credentials, req, option)
	authorization += "/" + option.signedHeadersToString() + "/" + signature

	req.addHeader("Authorization", authorization)

	return authorization
}

// Client is the base client struct for all products of baidubce.
type Client struct {
	Config
}

func (c *Client) GetUriPath(uriPath string) string {
	if strings.Index(uriPath, "/") == 0 {
		uriPath = uriPath[1:]
	}

	if c.APIVersion != "" {
		return fmt.Sprintf("/%s/%s", c.APIVersion, uriPath)
	}

	return fmt.Sprintf("/%s", uriPath)
}

// SendRequest sends a http request to the endpoint of baidubce api.
func (c *Client) SendRequest(req *Request, option *SignOption, autoReadAllBytesFromResponseBody bool) (*Response, *Error) {
	GenerateAuthorization(c.Credentials, *req, option)
	httpClient := http.Client{}
	res, err := httpClient.Do(req.raw())

	if err != nil {
		return nil, NewError(err)
	}

	bceResponse, err := NewResponse(res, autoReadAllBytesFromResponseBody)

	if err != nil {
		return nil, NewError(err)
	}

	if res.StatusCode >= 400 {
		if bceResponse.Body == nil {
			body, err := ioutil.ReadAll(bceResponse.Response.Body)

			if err != nil {
				return nil, NewError(err)
			}

			bceResponse.Body = body
		}

		return bceResponse, NewErrorFromJSON(bceResponse.Body)
	}

	return bceResponse, nil
}

func generateHeaderValidCompareFunc(headerKey string) func(string, string) bool {
	return func(key, value string) bool {
		return strings.ToLower(key) == strings.ToLower(headerKey) && value != ""
	}
}

// sign returns signed signature.
func sign(credentials Credentials, req Request, option *SignOption) string {
	signingKey := getSigningKey(credentials, option)
	req.prepareHeaders(option)
	canonicalRequest := req.canonical(option)
	signature := util.HmacSha256Hex(signingKey, canonicalRequest)

	return signature
}

func getSigningKey(credentials Credentials, option *SignOption) string {
	var authStringPrefix = fmt.Sprintf("bce-auth-v1/%s", credentials.AccessKeyID)
	authStringPrefix += "/" + option.Timestamp
	authStringPrefix += "/" + strconv.Itoa(option.ExpirationPeriodInSeconds)

	return util.HmacSha256Hex(credentials.SecretAccessKey, authStringPrefix)
}
