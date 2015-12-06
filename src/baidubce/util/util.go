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

// GetURIPath returns the path part of URI.
func GetURIPath(uri string) string {
	uri = strings.Replace(uri, "://", "", 1)
	index := strings.Index(uri, "/")
	return uri[index:]
}

// URIEncodeExceptSlash encodes all characters of a string except the slash character.
func URIEncodeExceptSlash(uri string) string {
	var result string

	for _, char := range uri {
		str := fmt.Sprintf("%c", char)
		if str == "/" {
			result += str
		} else {
			result += URLEncode(str)
		}
	}

	return result
}

// HmacSha256Hex returns a encrypted string.
func HmacSha256Hex(key, message string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

// Contains determines whether a string slice contains a certain value.
// Ignore case when comparing if case insensitive.
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

// MapContains determines whether the string map contains a uncertain value.
// The result is determined by compare function.
func MapContains(m map[string]string, compareFunc func(string, string) bool) bool {
	for key, value := range m {
		if compareFunc(key, value) {
			return true
		}
	}

	return false
}

// GetMapKey returns the key of the map for a certain value.
// Ignore case when comparing if case insensitive.
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

// GetMapValue returns the value of the map for a certain key.
// Ignore case when comparing if case insensitive.
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

// TimeToUTCString returns a utc string of a time instance.
func TimeToUTCString(t time.Time) string {
	format := time.RFC3339 // 2006-01-02T15:04:05Z07:00
	return t.UTC().Format(format)
}

// TimeStringToRFC1123 returns a formatted string of `time.RFC1123` format.
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

// HostToURL returns the whole URL string.
func HostToURL(host string) string {
	if matched, _ := regexp.MatchString("^[[:alpha:]]+:", host); matched {
		return host
	}

	return "http://" + host
}

// ToCanonicalQueryString returns the canonicalized query string.
func ToCanonicalQueryString(params map[string]string) string {
	if params == nil {
		return ""
	}

	encodedQueryStrings := make([]string, 0, 10)
	var query string

	for key, value := range params {
		if key != "" {
			query = URLEncode(key) + "="
			if value != "" {
				query += URLEncode(value)
			}
			encodedQueryStrings = append(encodedQueryStrings, query)
		}
	}

	sort.Strings(encodedQueryStrings)

	return strings.Join(encodedQueryStrings, "&")
}

// ToCanonicalHeaderString returns the canonicalized string.
func ToCanonicalHeaderString(headerMap map[string]string) string {
	headers := make([]string, 0, len(headerMap))
	for key, value := range headerMap {
		headers = append(headers,
			fmt.Sprintf("%s:%s", URLEncode(strings.ToLower(key)),
				URLEncode(strings.TrimSpace(value))))
	}

	sort.Strings(headers)

	return strings.Join(headers, "\n")
}

// URLEncode encodes a string like Javascript's encodeURIComponent()
func URLEncode(str string) string {
	// BUG(go): see https://github.com/golang/go/issues/4013
	// use %20 instead of the + character for encoding a space
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

// SliceToLower transforms each item of a slice to lowercase.
func SliceToLower(slice []string) {
	for index, value := range slice {
		slice[index] = strings.ToLower(value)
	}
}

// MapKeyToLower transforms each item of a map to lowercase.
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
