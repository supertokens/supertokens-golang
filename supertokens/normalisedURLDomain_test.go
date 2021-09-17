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

type NormalisedURLDomainTest struct {
	Input  string
	Output string
}

func TestNormaliseURLDomainOrThrowError(t *testing.T) {
	input := []NormalisedURLDomainTest{{
		Input:  "http://api.example.com",
		Output: "http://api.example.com",
	}, {
		Input:  "https://api.example.com",
		Output: "https://api.example.com",
	}, {
		Input:  "http://api.example.com?hello=1",
		Output: "http://api.example.com",
	}, {
		Input:  "http://api.example.com/hello",
		Output: "http://api.example.com",
	}, {
		Input:  "http://api.example.com/",
		Output: "http://api.example.com",
	}, {
		Input:  "http://api.example.com:8080",
		Output: "http://api.example.com:8080",
	}, {
		Input:  "http://api.example.com#random2",
		Output: "http://api.example.com",
	}, {
		Input:  "api.example.com/",
		Output: "https://api.example.com",
	}, {
		Input:  "api.example.com#random",
		Output: "https://api.example.com",
	}, {
		Input:  ".example.com",
		Output: "https://example.com",
	}, {
		Input:  "api.example.com/?hello=1&bye=2",
		Output: "https://api.example.com",
	}, {
		Input:  "localhost",
		Output: "http://localhost",
	}, {
		Input:  "https://localhost",
		Output: "https://localhost",
	}, {
		Input:  "http://api.example.com/one/two",
		Output: "http://api.example.com",
	}, {
		Input:  "http://1.2.3.4/one/two",
		Output: "http://1.2.3.4",
	}, {
		Input:  "https://1.2.3.4/one/two",
		Output: "https://1.2.3.4",
	}, {
		Input:  "1.2.3.4/one/two",
		Output: "http://1.2.3.4",
	}, {
		Input:  "https://api.example.com/one/two/",
		Output: "https://api.example.com",
	}, {
		Input:  "http://api.example.com/one/two?hello=1",
		Output: "http://api.example.com",
	}, {
		Input:  "http://api.example.com/one/two#random2",
		Output: "http://api.example.com",
	}, {
		Input:  "api.example.com/one/two",
		Output: "https://api.example.com",
	}, {
		Input:  "api.example.com/one/two/#random",
		Output: "https://api.example.com",
	}, {
		Input:  ".example.com/one/two",
		Output: "https://example.com",
	}, {
		Input:  "localhost:4000",
		Output: "http://localhost:4000",
	}, {
		Input:  "127.0.0.1:4000",
		Output: "http://127.0.0.1:4000",
	}, {
		Input:  "127.0.0.1",
		Output: "http://127.0.0.1",
	}, {
		Input:  "https://127.0.0.1:80",
		Output: "https://127.0.0.1:80",
	}, {
		Input:  "https://127.0.0.1:80/",
		Output: "https://127.0.0.1:80",
	}}
	for _, val := range input {
		domain, _ := NewNormalisedURLDomain(val.Input)
		assert.Equal(t, val.Output, domain.value)
	}

	{
		_, err := NewNormalisedURLDomain("/one/two")
		assert.Equal(t, err.Error(), "please provide a valid domain name")
	}

	{
		_, err := NewNormalisedURLDomain("/.netlify/functions/api")
		assert.Equal(t, err.Error(), "please provide a valid domain name")
	}
}
