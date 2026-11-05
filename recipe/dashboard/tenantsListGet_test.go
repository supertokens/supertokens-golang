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
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/multitenancy"
	"github.com/supertokens/supertokens-golang/recipe/multitenancy/multitenancymodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestTenantsListReturnsOnlyTenantIdsAndEnabledFlags(t *testing.T) {
	BeforeEach()
	connectionURI := unittesting.StartUpSTWithMultitenancy("localhost", "8080")
	defer AfterEach()
	config := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: connectionURI,
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			emailpassword.Init(nil),
			thirdparty.Init(nil),
			multitenancy.Init(nil),
			Init(&dashboardmodels.TypeInput{
				ApiKey: "testapikey",
			}),
		},
	}
	err := supertokens.Init(config)
	assert.NoError(t, err)

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	trueValue := true
	falseValue := false
	_, err = multitenancy.CreateOrUpdateTenant("t1", multitenancymodels.TenantConfig{
		EmailPasswordEnabled: &falseValue,
		PasswordlessEnabled:  &falseValue,
		ThirdPartyEnabled:    &trueValue,
	})
	assert.NoError(t, err)
	_, err = multitenancy.CreateOrUpdateThirdPartyConfig("t1", tpmodels.ProviderConfig{
		ThirdPartyId: "google",
		Clients: []tpmodels.ProviderClientConfig{
			{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret-value",
			},
		},
	}, nil)
	assert.NoError(t, err)

	req, _ := http.NewRequest(http.MethodGet, testServer.URL+"/auth/dashboard/api/tenants/list", nil)
	req.Header.Set("Authorization", "Bearer testapikey")
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	res.Body.Close()

	bodyString := string(body)
	assert.False(t, strings.Contains(bodyString, "providers"))
	assert.False(t, strings.Contains(bodyString, "clients"))
	assert.False(t, strings.Contains(bodyString, "test-client-id"))
	assert.False(t, strings.Contains(bodyString, "test-client-secret-value"))

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, "OK", response["status"])

	tenants := response["tenants"].([]interface{})
	var t1 map[string]interface{}
	for _, tenant := range tenants {
		tenantMap := tenant.(map[string]interface{})
		if tenantMap["tenantId"] == "t1" {
			t1 = tenantMap
		}
	}
	assert.NotNil(t, t1)
	assert.Equal(t, map[string]interface{}{
		"tenantId":      "t1",
		"emailPassword": map[string]interface{}{"enabled": false},
		"passwordless":  map[string]interface{}{"enabled": false},
		"thirdParty":    map[string]interface{}{"enabled": true},
	}, t1)
}
