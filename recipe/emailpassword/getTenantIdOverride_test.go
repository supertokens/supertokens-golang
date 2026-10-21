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

package emailpassword

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/multitenancy"
	"github.com/supertokens/supertokens-golang/recipe/multitenancy/multitenancymodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestGetTenantIdOverrideIsUsedWithAndWithoutRidHeader(t *testing.T) {
	BeforeEach()
	connectionURI := unittesting.StartUpSTWithMultitenancy("localhost", "8080")
	defer AfterEach()

	var tenantIdsFromFrontend []string
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: connectionURI,
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
			session.Init(nil),
			multitenancy.Init(&multitenancymodels.TypeInput{
				Override: &multitenancymodels.OverrideStruct{
					Functions: func(originalImplementation multitenancymodels.RecipeInterface) multitenancymodels.RecipeInterface {
						oGetTenantId := *originalImplementation.GetTenantId
						nGetTenantId := func(tenantIdFromFrontend string, userContext supertokens.UserContext) (string, error) {
							tenantIdsFromFrontend = append(tenantIdsFromFrontend, tenantIdFromFrontend)
							if tenantIdFromFrontend == "blocked" {
								return "", supertokens.BadInputError{Msg: "tenant is blocked"}
							}
							return oGetTenantId(tenantIdFromFrontend, userContext)
						}
						*originalImplementation.GetTenantId = nGetTenantId
						return originalImplementation
					},
				},
			}),
		},
	}
	err := supertokens.Init(configValue)
	assert.NoError(t, err)

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	signUp := func(path string, rid string, email string) (int, map[string]interface{}) {
		postBody, err := json.Marshal(map[string]interface{}{
			"formFields": []map[string]interface{}{
				{"id": "email", "value": email},
				{"id": "password", "value": "validpass123"},
			},
		})
		assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, testServer.URL+path, bytes.NewBuffer(postBody))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		if rid != "" {
			req.Header.Set("rid", rid)
		}
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		dataInBytes, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		resp.Body.Close()
		var data map[string]interface{}
		_ = json.Unmarshal(dataInBytes, &data)
		return resp.StatusCode, data
	}

	for _, rid := range []string{"emailpassword", ""} {
		status, data := signUp("/auth/blocked/signup", rid, "blocked-"+rid+"x@example.com")
		assert.Equal(t, 400, status, "rid: %q", rid)
		assert.Equal(t, "tenant is blocked", data["message"], "rid: %q", rid)
	}

	for _, rid := range []string{"emailpassword", ""} {
		tenantIdsFromFrontend = nil
		status, data := signUp("/auth/signup", rid, "public-"+rid+"x@example.com")
		assert.Equal(t, 200, status, "rid: %q", rid)
		assert.Equal(t, "OK", data["status"], "rid: %q", rid)
		assert.Equal(t, []string{"public"}, tenantIdsFromFrontend, "rid: %q", rid)
	}
}
