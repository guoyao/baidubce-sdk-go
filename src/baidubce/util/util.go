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
			result += url.QueryEscape(str)
		}
	}

	return result
}

func HmacSha256Hex(key, message string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

func Contains(slice []string, value string) bool {
	for _, v := range slice {
		if value == v {
			return true
		}
	}

	return false
}

func TimeToUTCString(t time.Time) string {
	format := time.RFC3339 // 2006-01-02T15:04:05Z07:00
	return t.UTC().Format(format)
}

func ToTestError(funcName, got, expected string) string {
	return fmt.Sprintf("%s failed. Got %s, expected %s", funcName, got, expected)
}
