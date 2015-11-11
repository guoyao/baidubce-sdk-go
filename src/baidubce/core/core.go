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
 * @file structure.go
 * @author guoyao
 */

// Package core define a set of core data structure, and implements a set of core functions
package core

import (
	"baidubce/util"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type Credentials struct {
	AccessKeyId     string
	SecretAccessKey string
}

type SignOption struct {
	Timestamp                 string
	ExpirationPeriodInSeconds uint
}

type Request struct {
	HttpMethod  string
	URI         string
	QueryString string
	Header      http.Header
}

var canonicalHeaders []string = []string{
	"host",
	"content-length",
	"content-type",
	"content-md5",
}

func NewCredentials(accessKeyId, secretAccessKey string) Credentials {
	return Credentials{accessKeyId, secretAccessKey}
}

func NewSignOption(timestamp string, expirationPeriodInSeconds uint) SignOption {
	return SignOption{timestamp, expirationPeriodInSeconds}
}

func NewRequest(httpMethod, URI, queryString string, header http.Header) Request {
	return Request{httpMethod, URI, queryString, header}
}

func Sign(credentials Credentials, req Request, option SignOption) string {
	signingKey := getSigningKey(credentials, option)
	canonicalRequest := req.canonical()
	signature := util.HmacSha256Hex(signingKey, canonicalRequest)

	return signature
}

func (req *Request) canonical() string {
	canonicalStrings := make([]string, 0, 4)

	canonicalHttpMethod := strings.ToUpper(req.HttpMethod)
	canonicalStrings = append(canonicalStrings, canonicalHttpMethod)

	canonicalURI := util.UriEncodeExceptSlash(util.GetUriPath(req.URI))
	canonicalStrings = append(canonicalStrings, canonicalURI)

	canonicalQueryString := getCanonicalQueryString(req.QueryString)
	canonicalStrings = append(canonicalStrings, canonicalQueryString)

	canonicalHeader := getCanonicalHeader(req.Header)
	canonicalStrings = append(canonicalStrings, canonicalHeader)

	return strings.Join(canonicalStrings, "\n")
}

func getSigningKey(credentials Credentials, option SignOption) string {
	var authStringPrefix = fmt.Sprintf("bce-auth-v1/%s", credentials.AccessKeyId)

	if option.Timestamp != "" {
		authStringPrefix += "/" + option.Timestamp
	}

	if option.ExpirationPeriodInSeconds > 0 {
		authStringPrefix += "/" + strconv.Itoa(int(option.ExpirationPeriodInSeconds))
	}

	return util.HmacSha256Hex(credentials.SecretAccessKey, authStringPrefix)
}

func getCanonicalQueryString(queryString string) string {
	arr := strings.Split(queryString, "&")
	encodedQueryStrings := make([]string, 0, 10)

	for _, value := range arr {
		if value != "" {
			keyValueArr := strings.Split(value, "=")
			query := url.QueryEscape(keyValueArr[0]) + "="
			if len(keyValueArr) > 1 {
				query += url.QueryEscape(keyValueArr[1])
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
			headerValue := url.QueryEscape(strings.TrimSpace(value[0]))
			headers = append(headers, fmt.Sprintf("%s:%s", strings.ToLower(key), headerValue))
		}
	}

	sort.Strings(headers)

	return strings.Join(headers, "\n")
}

func isCanonicalHeader(key string) bool {
	key = strings.ToLower(key)
	return util.Contains(canonicalHeaders, key) || strings.Index(key, "x-bce-") == 0
}
