/* Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
 *
 * This software is licensed under the Apache License, Version 2.0 (the
 * "License") as published by the Apache Software Foundation.
 *
 * You may not use this file except in compliance with the License. You may
 * obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

package supertokens

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type NormalisedURLPathTest struct {
	Input  string
	Output string
}

func TestNormaliseURLPathOrThrowError(t *testing.T) {
	input := []NormalisedURLPathTest{{
		Input:  "exists?email=john.doe%40gmail.com",
		Output: "/exists",
	}, {
		Input:  "/auth/email/exists?email=john.doe%40gmail.com",
		Output: "/auth/email/exists",
	}, {
		Input:  "exists",
		Output: "/exists",
	}, {
		Input:  "/exists",
		Output: "/exists",
	}, {
		Input:  "/exists?email=john.doe%40gmail.com",
		Output: "/exists",
	}, {
		Input:  "http://api.example.com",
		Output: "",
	}, {
		Input:  "https://api.example.com",
		Output: "",
	}, {
		Input:  "http://api.example.com?hello=1",
		Output: "",
	}, {
		Input:  "http://api.example.com/hello",
		Output: "/hello",
	}, {
		Input:  "http://api.example.com/",
		Output: "",
	}, {
		Input:  "http://api.example.com:8080",
		Output: "",
	}, {
		Input:  "http://api.example.com#random2",
		Output: "",
	}, {
		Input:  "api.example.com/",
		Output: "",
	}, {
		Input:  "api.example.com#random",
		Output: "",
	}, {
		Input:  ".example.com",
		Output: "",
	}, {
		Input:  "api.example.com/?hello=1&bye=2",
		Output: "",
	}, {
		Input:  "http://api.example.com/one/two",
		Output: "/one/two",
	}, {
		Input:  "http://1.2.3.4/one/two",
		Output: "/one/two",
	}, {
		Input:  "1.2.3.4/one/two",
		Output: "/one/two",
	}, {
		Input:  "https://api.example.com/one/two/",
		Output: "/one/two",
	}, {
		Input:  "http://api.example.com/one/two?hello=1",
		Output: "/one/two",
	}, {
		Input:  "http://api.example.com/hello",
		Output: "/hello",
	}, {
		Input:  "http://api.example.com/one/two/",
		Output: "/one/two",
	}, {
		Input:  "http://api.example.com:8080/one/two",
		Output: "/one/two",
	}, {
		Input:  "http://api.example.com/one/two#random2",
		Output: "/one/two",
	}, {
		Input:  "api.example.com/one/two",
		Output: "/one/two",
	}, {
		Input:  "api.example.com/one/two/#random",
		Output: "/one/two",
	}, {
		Input:  ".example.com/one/two",
		Output: "/one/two",
	}, {
		Input:  "api.example.com/one/two?hello=1&bye=2",
		Output: "/one/two",
	}, {
		Input:  "/one/two",
		Output: "/one/two",
	}, {
		Input:  "one/two",
		Output: "/one/two",
	}, {
		Input:  "one/two/",
		Output: "/one/two",
	}, {
		Input:  "/one",
		Output: "/one",
	}, {
		Input:  "one",
		Output: "/one",
	}, {
		Input:  "one/",
		Output: "/one",
	}, {
		Input:  "/one/two/",
		Output: "/one/two",
	}, {
		Input:  "/one/two?hello=1",
		Output: "/one/two",
	}, {
		Input:  "one/two?hello=1",
		Output: "/one/two",
	}, {
		Input:  "/one/two/#random",
		Output: "/one/two",
	}, {
		Input:  "one/two#random",
		Output: "/one/two",
	}, {
		Input:  "localhost:4000/one/two",
		Output: "/one/two",
	}, {
		Input:  "127.0.0.1:4000/one/two",
		Output: "/one/two",
	}, {
		Input:  "127.0.0.1/one/two",
		Output: "/one/two",
	}, {
		Input:  "https://127.0.0.1:80/one/two",
		Output: "/one/two",
	}, {
		Input:  "/",
		Output: "",
	}, {
		Input:  "/.netlify/functions/api",
		Output: "/.netlify/functions/api",
	}, {
		Input:  "/netlify/.functions/api",
		Output: "/netlify/.functions/api",
	}, {
		Input:  "app.example.com/.netlify/functions/api",
		Output: "/.netlify/functions/api",
	}, {
		Input:  "app.example.com/netlify/.functions/api",
		Output: "/netlify/.functions/api",
	}, {
		Input:  "/app.example.com",
		Output: "/app.example.com",
	}}
	for _, val := range input {
		path, _ := NewNormalisedURLPath(val.Input)
		assert.Equal(t, val.Output, path.value, val.Input)
	}
}
