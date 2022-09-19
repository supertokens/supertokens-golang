/* Copyright (c) 2022, VRAI Labs and/or its affiliates. All rights reserved.
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

package userroles

import (
	"net/http"
	"testing"

	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func resetAll() {
	supertokens.ResetForTest()
	ResetForTest()
	session.ResetForTest()
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

func canRunTest(t *testing.T) bool {
	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}
	if unittesting.MaxVersion("2.13", cdiVersion) == "2.13" {
		return false
	}
	return true
}

type fakeRes struct{}

func (f fakeRes) Header() http.Header {
	return http.Header{}
}

func (f fakeRes) Write(body []byte) (int, error) {
	return len(body), nil
}

func (f fakeRes) WriteHeader(statusCode int) {}
