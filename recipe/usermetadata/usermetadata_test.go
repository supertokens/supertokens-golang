/*
 * Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
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

package usermetadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/usermetadata/usermetadatamodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestClearMetadata(t *testing.T) {
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
	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}
	if unittesting.MaxVersion("2.13", cdiVersion) == cdiVersion {
		err := ClearUserMetadata("userId")
		if err != nil {
			t.Error(err.Error())
		}
	}
}

func TestShouldClearStoredField(t *testing.T) {
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
	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}
	if unittesting.MaxVersion("2.13", cdiVersion) == cdiVersion {
		_, err := UpdateUserMetadata("userId", map[string]interface{}{
			"role": "admin",
		})
		if err != nil {
			t.Error(err.Error())
		}

		err = ClearUserMetadata("userId")
		if err != nil {
			t.Error(err.Error())
		}

		metadata, err := GetUserMetadata("userId")
		if err != nil {
			t.Error(err.Error())
		}
		assert.Equal(t, metadata, map[string]interface{}{})

		updatedContent, err := UpdateUserMetadata("userId", map[string]interface{}{
			"role": "admin2",
		})
		if err != nil {
			t.Error(err.Error())
		}

		assert.Equal(t, updatedContent, map[string]interface{}{
			"role": "admin2",
		})
	}
}

func TestGetUserMetadata(t *testing.T) {
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
	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}
	if unittesting.MaxVersion("2.13", cdiVersion) == cdiVersion {
		metadata, err := GetUserMetadata("userId")
		if err != nil {
			t.Error(err.Error())
		}
		assert.Equal(t, metadata, map[string]interface{}{})
	}
}

func TestGetUserMetadataIfCreated(t *testing.T) {
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
	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}
	if unittesting.MaxVersion("2.13", cdiVersion) == cdiVersion {
		updatedContent, err := UpdateUserMetadata("userId", map[string]interface{}{
			"role": "admin",
		})
		if err != nil {
			t.Error(err.Error())
		}
		assert.Equal(t, updatedContent, map[string]interface{}{
			"role": "admin",
		})

		metadata, err := GetUserMetadata("userId")
		if err != nil {
			t.Error(err.Error())
		}
		assert.Equal(t, metadata, map[string]interface{}{
			"role": "admin",
		})
	}
}

func TestOverride(t *testing.T) {
	calledOverride := false
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
			Init(&usermetadatamodels.TypeInput{
				Override: &usermetadatamodels.OverrideStruct{
					Functions: func(originalImplementation usermetadatamodels.RecipeInterface) usermetadatamodels.RecipeInterface {
						originalUpdateMetadata := *originalImplementation.UpdateUserMetadata

						(*originalImplementation.UpdateUserMetadata) = func(userID string, metadataUpdate map[string]interface{}, userContext supertokens.UserContext) (map[string]interface{}, error) {
							calledOverride = true
							return originalUpdateMetadata(userID, metadataUpdate, userContext)
						}

						return originalImplementation
					},
				},
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
	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}
	if unittesting.MaxVersion("2.13", cdiVersion) == cdiVersion {
		_, err := UpdateUserMetadata("userId", map[string]interface{}{
			"role": "admin",
		})
		if err != nil {
			t.Error(err.Error())
		}

		metadata, err := GetUserMetadata("userId")
		if err != nil {
			t.Error(err.Error())
		}
		assert.Equal(t, metadata, map[string]interface{}{
			"role": "admin",
		})

		assert.Equal(t, calledOverride, true)
	}
}

func TestShouldUpdateShallowMerge(t *testing.T) {
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
	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}
	if unittesting.MaxVersion("2.13", cdiVersion) == cdiVersion {
		testMetadata := map[string]interface{}{
			"updated": map[string]interface{}{
				"subObjectNull":    "this will become null",
				"subObjectCleared": "this will be removed",
				"subObjectUpdate":  "this will become a number",
			},
			"cleared":   "this should not be on the end result",
			"unchanged": float64(123),
		}

		testMetadataUpdate := map[string]interface{}{
			"updated": map[string]interface{}{
				"subObjectNull":    nil,
				"subObjectUpdate":  float64(123),
				"subObjectNewProp": "this will appear",
			},
			"cleared":     nil,
			"newRootProp": "this should appear on the end result",
		}

		expectedResult := map[string]interface{}{
			"updated": map[string]interface{}{
				"subObjectNull":    nil,
				"subObjectUpdate":  float64(123),
				"subObjectNewProp": "this will appear",
			},
			"newRootProp": "this should appear on the end result",
			"unchanged":   float64(123),
		}

		updatedContent, err := UpdateUserMetadata("userId", testMetadata)
		if err != nil {
			t.Error(err.Error())
		}
		assert.Equal(t, updatedContent, testMetadata)

		updatedContent, err = UpdateUserMetadata("userId", testMetadataUpdate)
		if err != nil {
			t.Error(err.Error())
		}

		assert.Equal(t, updatedContent, expectedResult)

		finalResult, err := GetUserMetadata("userId")
		if err != nil {
			t.Error(err.Error())
		}

		assert.Equal(t, finalResult, expectedResult)

	}
}
