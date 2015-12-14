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
 * @file request.go
 * @author guoyao
 */

// Package baidubce define a set of core data structure, and implements a set of core functions.
package baidubce

import (
	"io/ioutil"
	"net/http"
)

// Response holds an instance of type `http response`, and has some custom data and functions.
type Response struct {
	Body []byte
	*http.Response
}

// NewResponse returns an instance of type `Response`
func NewResponse(res *http.Response) *Response {
	response := &Response{Response: res}
	body, _ := ioutil.ReadAll(res.Body)
	response.Body = body
	return response
}
