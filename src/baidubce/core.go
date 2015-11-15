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

// Package core define a set of core data structure, and implements a set of core functions
package baidubce

import (
	"baidubce/util"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	EXPIRATION_PERIOD_IN_SECONDS = 1800
)

var Region map[string]string = map[string]string{
	"bj": "http://bj.bcebos.com",
}

type Credentials struct {
	AccessKeyId     string
	SecretAccessKey string
}

type SignOption struct {
	Timestamp                 string
	ExpirationPeriodInSeconds int
}

type Request struct {
	HttpMethod string
	URI        string
	Params     map[string]string
	Header     http.Header
}

var canonicalHeaders []string = []string{
	"host",
	"content-length",
	"content-type",
	"content-md5",
}

func NewCredentials(accessKeyId, secretAccessKey string) *Credentials {
	return &Credentials{accessKeyId, secretAccessKey}
}

func NewSignOption(timestamp string, expirationPeriodInSeconds int) *SignOption {
	option := &SignOption{timestamp, expirationPeriodInSeconds}
	option.init()

	return option
}

func NewRequest(httpMethod, URI string, params map[string]string, header http.Header) *Request {
	return &Request{httpMethod, URI, params, header}
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
	authorization += "//" + signature

	return authorization
}

func (option *SignOption) init() {
	if option.Timestamp == "" {
		option.Timestamp = util.TimeToUTCString(time.Now())
	}
	if option.ExpirationPeriodInSeconds <= 0 && option.ExpirationPeriodInSeconds != -1 {
		option.ExpirationPeriodInSeconds = EXPIRATION_PERIOD_IN_SECONDS
	}
}

func (req *Request) canonical() string {
	canonicalStrings := make([]string, 0, 4)

	canonicalHttpMethod := strings.ToUpper(req.HttpMethod)
	canonicalStrings = append(canonicalStrings, canonicalHttpMethod)

	canonicalURI := util.UriEncodeExceptSlash(util.GetUriPath(req.URI))
	canonicalStrings = append(canonicalStrings, canonicalURI)

	canonicalQueryString := getCanonicalQueryString(req.Params)
	canonicalStrings = append(canonicalStrings, canonicalQueryString)

	canonicalHeader := getCanonicalHeader(req.Header)
	canonicalStrings = append(canonicalStrings, canonicalHeader)

	return strings.Join(canonicalStrings, "\n")
}

func (req *Request) ParamsToCanonicalQueryString() string {
	return getCanonicalQueryString(req.Params)
}

// generate signature
func sign(credentials Credentials, req Request, option *SignOption) string {
	signingKey := getSigningKey(credentials, option)
	req.Header.Add("x-bce-date", option.Timestamp)
	canonicalRequest := req.canonical()
	signature := util.HmacSha256Hex(signingKey, canonicalRequest)

	return signature
}

func getSigningKey(credentials Credentials, option *SignOption) string {
	var authStringPrefix = fmt.Sprintf("bce-auth-v1/%s", credentials.AccessKeyId)
	authStringPrefix += "/" + option.Timestamp
	authStringPrefix += "/" + strconv.Itoa(option.ExpirationPeriodInSeconds)

	return util.HmacSha256Hex(credentials.SecretAccessKey, authStringPrefix)
}

func getCanonicalQueryString(params map[string]string) string {
	if params == nil {
		return ""
	}

	encodedQueryStrings := make([]string, 0, 10)
	var query string

	for key, value := range params {
		if key != "" {
			query = url.QueryEscape(key) + "="
			if value != "" {
				query += url.QueryEscape(value)
			}
			encodedQueryStrings = append(encodedQueryStrings, query)
		}
	}

	sort.Strings(encodedQueryStrings)

	return strings.Join(encodedQueryStrings, "&")
}

func getCanonicalHeader(header http.Header) string {
	headers := make([]string, 0, len(header))
	for key, value := range header {
		if isCanonicalHeader(key) {
			headers = append(headers, fmt.Sprintf("%s:%s",
				url.QueryEscape(strings.ToLower(key)),
				url.QueryEscape(strings.TrimSpace(value[0]))))
		}
	}

	sort.Strings(headers)

	return strings.Join(headers, "\n")
}

func isCanonicalHeader(key string) bool {
	key = strings.ToLower(key)
	return util.Contains(canonicalHeaders, key) || strings.Index(key, "x-bce-") == 0
}
