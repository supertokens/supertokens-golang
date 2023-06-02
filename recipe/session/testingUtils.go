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
	"net/http"

	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

// Testing constants
var didGetSessionCallCore = false
var returnedFromCache = false
var urlsAttemptedForJWKSFetch []string

func resetAll() {
	supertokens.ResetForTest()
	ResetForTest()
	didGetSessionCallCore = false
	returnedFromCache = false
	urlsAttemptedForJWKSFetch = []string{}
	jwksCache = nil
}

func BeforeEach() {
	unittesting.KillAllST()
	resetAll()
	unittesting.SetUpST()
}

func AfterEach() {
	unittesting.KillAllST()
	resetAll()
	unittesting.CleanST()
}

type fakeRes struct{}

func (f fakeRes) Header() http.Header {
	return http.Header{}
}

func (f fakeRes) Write(body []byte) (int, error) {
	return len(body), nil
}

func (f fakeRes) WriteHeader(statusCode int) {}
