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

package openid

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/openid/openidmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestCreateJWTUsesConfiguredJwtValiditySeconds(t *testing.T) {
	BeforeEach()
	connectionURI := unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	jwtValiditySeconds := uint64(60)
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
			Init(&openidmodels.TypeInput{
				JwtValiditySeconds: &jwtValiditySeconds,
			}),
		},
	}
	err := supertokens.Init(configValue)
	assert.NoError(t, err)

	now := time.Now().Unix()
	resp, err := CreateJWT(map[string]interface{}{"customClaim": "value"}, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp.OK)

	claims := jwt.MapClaims{}
	_, _, err = jwt.NewParser().ParseUnverified(resp.OK.Jwt, claims)
	assert.NoError(t, err)
	exp, ok := claims["exp"].(float64)
	assert.True(t, ok)
	assert.InDelta(t, float64(now+60), exp, 10)
}
