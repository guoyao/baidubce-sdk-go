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
 * @file util.go
 * @author guoyao
 */

// Package util implements a set of util functions.
package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
)

func GetUriPath(uri string) string {
	uri = strings.Replace(uri, "://", "", 1)
	index := strings.Index(uri, "/")
	return uri[index:]
}

func UriEncodeExceptSlash(uri string) string {
	var result string

	for _, char := range uri {
		str := fmt.Sprintf("%c", char)
		if str == "/" {
			result += str
		} else {
			result += UrlEncode(str)
		}
	}

	return result
}

func HmacSha256Hex(key, message string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

// Whether the string slice contains a certain value
// Ignore case when comparing if case insensitive
func Contains(slice []string, value string, caseInsensitive bool) bool {
	if caseInsensitive {
		value = strings.ToLower(value)
	}

	for _, v := range slice {
		if caseInsensitive {
			v = strings.ToLower(v)
		}

		if value == v {
			return true
		}
	}

	return false
}

// Whether the string map contains a uncertain value
// The result is determined by compare function
func MapContains(m map[string]string, compareFunc func(string, string) bool) bool {
	for key, value := range m {
		if compareFunc(key, value) {
			return true
		}
	}

	return false
}

func GetMapKey(m map[string]string, key string, caseInsensitive bool) string {
	if caseInsensitive {
		key = strings.ToLower(key)
	}

	var tempKey string

	for k := range m {
		tempKey = k

		if caseInsensitive {
			tempKey = strings.ToLower(k)
		}

		if tempKey == key {
			return k
		}
	}

	return ""
}

func GetMapValue(m map[string]string, key string, caseInsensitive bool) string {
	if caseInsensitive {
		for k, v := range m {
			if strings.ToLower(k) == strings.ToLower(key) {
				return v
			}
		}
	}

	return m[key]
}

func TimeToUTCString(t time.Time) string {
	format := time.RFC3339 // 2006-01-02T15:04:05Z07:00
	return t.UTC().Format(format)
}

// format string to time.RFC1123
func TimeStringToRFC1123(str string) string {
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		t, err = time.Parse(time.RFC1123, str)
		if err != nil {
			panic("Time format invalid. The time format must be time.RFC3339 or time.RFC1123")
		}
	}

	return t.Format(time.RFC1123)
}

func HostToUrl(host string) string {
	if matched, _ := regexp.MatchString("^[[:alpha:]]+:", host); matched {
		return host
	}

	return "http://" + host
}

func ToCanonicalQueryString(params map[string]string) string {
	if params == nil {
		return ""
	}

	encodedQueryStrings := make([]string, 0, 10)
	var query string

	for key, value := range params {
		if key != "" {
			query = UrlEncode(key) + "="
			if value != "" {
				query += UrlEncode(value)
			}
			encodedQueryStrings = append(encodedQueryStrings, query)
		}
	}

	sort.Strings(encodedQueryStrings)

	return strings.Join(encodedQueryStrings, "&")
}

func ToCanonicalHeaderString(headerMap map[string]string) string {
	headers := make([]string, 0, len(headerMap))
	for key, value := range headerMap {
		headers = append(headers,
			fmt.Sprintf("%s:%s", UrlEncode(strings.ToLower(key)),
				UrlEncode(strings.TrimSpace(value))))
	}

	sort.Strings(headers)

	return strings.Join(headers, "\n")
}

// UrlEncoded encodes a string like Javascript's encodeURIComponent()
func UrlEncode(str string) string {
	// BUG(go): see https://github.com/golang/go/issues/4013
	// use %20 instead of the + character for encoding a space
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

func SliceToLower(slice []string) {
	for index, value := range slice {
		slice[index] = strings.ToLower(value)
	}
}

func MapKeyToLower(m map[string]string) {
	temp := make(map[string]string, len(m))
	for key, value := range m {
		temp[strings.ToLower(key)] = value
		delete(m, key)
	}
	for key, value := range temp {
		m[key] = value
	}
}
