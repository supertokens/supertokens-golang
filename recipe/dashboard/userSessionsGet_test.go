/* Copyright (c) 2026, VRAI Labs and/or its affiliates. All rights reserved.
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

package dashboard

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

// fakeCoreForUserSessions answers the core calls made by the user sessions API,
// with the response fields that the session recipe reads. A handle listed in
// failingHandles gets a 500, as a core would answer on an internal error.
func fakeCoreForUserSessions(handles []string, failingHandles map[string]bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/apiversion":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"versions": []string{"3.1"}})
		case strings.HasSuffix(r.URL.Path, "/recipe/session/user"):
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"status": "OK", "sessionHandles": handles})
		case strings.HasSuffix(r.URL.Path, "/recipe/session"):
			handle := r.URL.Query().Get("sessionHandle")
			if failingHandles[handle] {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("internal error"))
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"status":             "OK",
				"sessionHandle":      handle,
				"userId":             "user-1",
				"tenantId":           "public",
				"userDataInDatabase": map[string]interface{}{},
				"userDataInJWT":      map[string]interface{}{},
				"expiry":             1,
				"timeCreated":        1,
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func setUpUserSessionsAPI(t *testing.T, core *httptest.Server) *httptest.Server {
	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{ConnectionURI: core.URL},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			session.Init(nil),
			Init(&dashboardmodels.TypeInput{ApiKey: "testapikey"}),
		},
	})
	assert.NoError(t, err)
	return httptest.NewServer(supertokens.Middleware(http.NewServeMux()))
}

func getUserSessions(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url+"/auth/dashboard/api/user/sessions?userId=user-1", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer testapikey")
	return (&http.Client{Timeout: 10 * time.Second}).Do(req)
}

// Run with -race. The handler reads every session in its own goroutine, so a
// user context or a slice shared between them is reported as a data race here,
// and shows up in production as "fatal error: concurrent map writes".
func TestThatUserSessionsGetReturnsAllSessionsWhenReadConcurrently(t *testing.T) {
	resetAll()
	defer resetAll()

	handles := []string{}
	for i := 0; i < 64; i++ {
		handles = append(handles, fmt.Sprintf("handle-%d", i))
	}
	core := fakeCoreForUserSessions(handles, nil)
	defer core.Close()
	server := setUpUserSessionsAPI(t, core)
	defer server.Close()

	res, err := getUserSessions(server.URL)
	if !assert.NoError(t, err) {
		return
	}
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	var response struct {
		Status   string `json:"status"`
		Sessions []struct {
			SessionHandle string `json:"sessionHandle"`
		} `json:"sessions"`
	}
	assert.NoError(t, json.Unmarshal(body, &response))
	assert.Equal(t, "OK", response.Status)

	returned := []string{}
	for _, s := range response.Sessions {
		returned = append(returned, s.SessionHandle)
	}
	assert.ElementsMatch(t, handles, returned)
}

func TestThatUserSessionsGetAnswersWhenReadingOneSessionFails(t *testing.T) {
	resetAll()
	defer resetAll()

	handles := []string{"handle-ok", "handle-bad"}
	core := fakeCoreForUserSessions(handles, map[string]bool{"handle-bad": true})
	defer core.Close()
	server := setUpUserSessionsAPI(t, core)

	res, err := getUserSessions(server.URL)
	// Closed only on success: httptest's Close waits for in-flight requests, so
	// it would block forever on a handler that never answers.
	if !assert.NoError(t, err, "the API did not answer") {
		return
	}
	defer server.Close()
	defer res.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}
