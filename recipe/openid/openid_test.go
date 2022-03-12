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

package openid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/openid/openidmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestWithDefaultConfigDisoveryConfigurationIsAsExpected(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.8") == "2.8" {
		return
	}

	resp, err := GetOpenIdDiscoveryConfiguration()
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "https://api.supertokens.io/auth", resp.OK.Issuer)
	assert.Equal(t, "https://api.supertokens.io/auth/jwt/jwks.json", resp.OK.Jwks_uri)
}

func TestWithAPIBasePathDisoveryConfigurationIsAsExpected(t *testing.T) {
	customAPIBasePath := "/"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIBasePath:   &customAPIBasePath,
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.8") == "2.8" {
		return
	}

	resp, err := GetOpenIdDiscoveryConfiguration()
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "https://api.supertokens.io", resp.OK.Issuer)
	assert.Equal(t, "https://api.supertokens.io/jwt/jwks.json", resp.OK.Jwks_uri)
}

func TestWithIssuerDisoveryConfigurationIsAsExpected(t *testing.T) {
	customIssuer := "https://cusomissuer/auth"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&openidmodels.TypeInput{
				Issuer: &customIssuer,
			}),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.8") == "2.8" {
		return
	}

	resp, err := GetOpenIdDiscoveryConfiguration()
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "https://cusomissuer/auth", resp.OK.Issuer)
	assert.Equal(t, "https://cusomissuer/auth/jwt/jwks.json", resp.OK.Jwks_uri)
}
