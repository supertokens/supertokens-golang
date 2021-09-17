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

package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type GetTopLevelDomainForSameSiteResolutionTest struct {
	Input  string
	Output string
}

func TestGetTopLevelDomainForSameSiteResolution(t *testing.T) {
	input := []GetTopLevelDomainForSameSiteResolutionTest{{
		Input:  "http://a.b.test.com",
		Output: "test.com",
	}, {
		Input:  "https://a.b.test.com",
		Output: "test.com",
	}, {
		Input:  "http://a.b.test.co.uk",
		Output: "test.co.uk",
	}, {
		Input:  "http://test.com",
		Output: "test.com",
	}, {
		Input:  "https://test.com",
		Output: "test.com",
	}, {
		Input:  "http://localhost",
		Output: "localhost",
	}, {
		Input:  "http://localhost.org",
		Output: "localhost",
	}, {
		Input:  "http://8.8.8.8",
		Output: "localhost",
	}, {
		Input:  "http://8.8.8.8:8080",
		Output: "localhost",
	}, {
		Input:  "http://localhost:3000",
		Output: "localhost",
	}, {
		Input:  "http://test.com:3567",
		Output: "test.com",
	}, {
		Input:  "https://test.com:3567",
		Output: "test.com",
	}}
	for _, val := range input {
		domain, _ := GetTopLevelDomainForSameSiteResolution(val.Input)
		assert.Equal(t, val.Output, domain, val.Input)
	}
}
