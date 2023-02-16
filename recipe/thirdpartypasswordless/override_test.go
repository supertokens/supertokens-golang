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

package thirdpartypasswordless

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/tplmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
	"gopkg.in/h2non/gock.v1"
)

func TestOverridingFunctions(t *testing.T) {
	var userRef *tplmodels.User
	var newUser bool
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.ProviderInput{customProvider1},
				Override: &tplmodels.OverrideStruct{
					Functions: func(originalImplementation tplmodels.RecipeInterface) tplmodels.RecipeInterface {
						originalThirdPartySignInUp := *originalImplementation.ThirdPartySignInUp
						*originalImplementation.ThirdPartySignInUp = func(thirdPartyID, thirdPartyUserID, email string, oAuthTokens tpmodels.TypeOAuthTokens, rawUserInfoFromProvider tpmodels.TypeRawUserInfoFromProvider, userContext supertokens.UserContext) (tplmodels.ThirdPartySignInUp, error) {
							resp, err := originalThirdPartySignInUp(thirdPartyID, thirdPartyUserID, email, oAuthTokens, rawUserInfoFromProvider, userContext)
							userRef = &resp.OK.User
							newUser = resp.OK.CreatedNewUser
							return resp, err
						}
						originalGetUserById := *originalImplementation.GetUserByID
						*originalImplementation.GetUserByID = func(userID string, tenantId *string, userContext supertokens.UserContext) (*tplmodels.User, error) {
							resp, err := originalGetUserById(userID, tenantId, userContext)
							userRef = resp
							return resp, err
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
	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/user", func(rw http.ResponseWriter, r *http.Request) {
		userId := r.URL.Query().Get("userId")
		fetchedUser, err := GetUserByID(userId, nil)
		if err != nil {
			t.Error(err.Error())
		}
		jsonResp, err := json.Marshal(fetchedUser)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)
	})

	defer gock.OffAll()
	gock.New("https://test.com").
		Post("/oauth/token").
		Persist().
		Reply(200).
		JSON(map[string]string{})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	formFields := map[string]interface{}{
		"thirdPartyId": "custom",
		"redirectURIInfo": map[string]interface{}{
			"redirectURIOnProviderDashboard": testServer.URL + "/callback",
			"redirectURIQueryParams": map[string]interface{}{
				"code": "abcdefghj",
			},
		},
	}

	postBody, err := json.Marshal(formFields)
	if err != nil {
		t.Error(err.Error())
	}

	gock.New(testServer.URL).EnableNetworking().Persist()
	gock.New("http://localhost:8080/").EnableNetworking().Persist()

	resp, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	signUpResponse := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	fetchedUser := signUpResponse["user"].(map[string]interface{})

	assert.NotNil(t, userRef)
	assert.True(t, newUser)
	assert.Equal(t, fetchedUser["email"], *userRef.Email)
	assert.Equal(t, fetchedUser["id"], userRef.ID)
	assert.Equal(t, fetchedUser["thirdParty"].(map[string]interface{})["id"], userRef.ThirdParty.ID)
	assert.Equal(t, fetchedUser["thirdParty"].(map[string]interface{})["userId"], userRef.ThirdParty.UserID)

	userRef = nil
	assert.Nil(t, userRef)

	formFields = map[string]interface{}{
		"thirdPartyId": "custom",
		"redirectURIInfo": map[string]interface{}{
			"redirectURIOnProviderDashboard": testServer.URL + "/callback",
			"redirectURIQueryParams": map[string]interface{}{
				"code": "abcdefghj",
			},
		},
	}

	postBody, err = json.Marshal(formFields)
	if err != nil {
		t.Error(err.Error())
	}

	resp, err = http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	signInResponse := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	fetchedUserFromSignIn := signInResponse["user"].(map[string]interface{})

	assert.NotNil(t, userRef)
	assert.False(t, newUser)
	assert.Equal(t, fetchedUserFromSignIn["email"], *userRef.Email)
	assert.Equal(t, fetchedUserFromSignIn["id"], userRef.ID)
	assert.Equal(t, fetchedUserFromSignIn["thirdParty"].(map[string]interface{})["id"], userRef.ThirdParty.ID)
	assert.Equal(t, fetchedUserFromSignIn["thirdParty"].(map[string]interface{})["userId"], userRef.ThirdParty.UserID)

	userRef = nil
	assert.Nil(t, userRef)

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/user", nil)
	assert.NoError(t, err)

	query := req.URL.Query()
	query.Add("userId", fetchedUserFromSignIn["id"].(string))
	req.URL.RawQuery = query.Encode()

	res, err := http.DefaultClient.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	userByIdResponse := *unittesting.HttpResponseToConsumableInformation(res.Body)

	assert.NotNil(t, userRef)
	assert.Equal(t, userByIdResponse["email"], *userRef.Email)
	assert.Nil(t, userByIdResponse["phoneNumber"])
	assert.Equal(t, userByIdResponse["id"], userRef.ID)
	assert.Equal(t, userByIdResponse["thirdParty"].(map[string]interface{})["id"], userRef.ThirdParty.ID)
	assert.Equal(t, userByIdResponse["thirdParty"].(map[string]interface{})["userId"], userRef.ThirdParty.UserID)
}

func TestOverridingAPIs(t *testing.T) {
	var userRef *tplmodels.User
	var newUser bool
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
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			Init(tplmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
					CreateAndSendCustomEmail: func(email string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						return nil
					},
				},
				Providers: []tpmodels.ProviderInput{customProvider1},
				Override: &tplmodels.OverrideStruct{
					APIs: func(originalImplementation tplmodels.APIInterface) tplmodels.APIInterface {
						originalThirdPartySignInUpPost := *originalImplementation.ThirdPartySignInUpPOST
						*originalImplementation.ThirdPartySignInUpPOST = func(provider *tpmodels.TypeProvider, input tpmodels.TypeSignInUpInput, options tpmodels.APIOptions, userContext supertokens.UserContext) (tplmodels.ThirdPartySignInUpPOSTResponse, error) {
							resp, err := originalThirdPartySignInUpPost(provider, input, options, userContext)
							userRef = &resp.OK.User
							newUser = resp.OK.CreatedNewUser
							return resp, err
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
	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	apiV, err := q.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	if unittesting.MaxVersion(apiV, "2.11") == "2.11" {
		return
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/user", func(rw http.ResponseWriter, r *http.Request) {
		userId := r.URL.Query().Get("userId")
		fetchedUser, err := GetUserByID(userId, nil)
		if err != nil {
			t.Error(err.Error())
		}
		jsonResp, err := json.Marshal(fetchedUser)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)
	})

	defer gock.OffAll()
	gock.New("https://test.com").
		Post("/oauth/token").
		Persist().
		Reply(200).
		JSON(map[string]string{})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	formFields := map[string]interface{}{
		"thirdPartyId": "custom",
		"redirectURIInfo": map[string]interface{}{
			"redirectURIOnProviderDashboard": testServer.URL + "/callback",
			"redirectURIQueryParams": map[string]interface{}{
				"code": "abcdefghj",
			},
		},
	}

	postBody, err := json.Marshal(formFields)
	if err != nil {
		t.Error(err.Error())
	}

	gock.New(testServer.URL).EnableNetworking().Persist()
	gock.New("http://localhost:8080/").EnableNetworking().Persist()

	resp, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	signUpResponse := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	fetchedUser := signUpResponse["user"].(map[string]interface{})

	assert.NotNil(t, userRef)
	assert.True(t, newUser)
	assert.Equal(t, fetchedUser["email"], *userRef.Email)
	assert.Equal(t, fetchedUser["id"], userRef.ID)
	assert.Equal(t, fetchedUser["thirdParty"].(map[string]interface{})["id"], userRef.ThirdParty.ID)
	assert.Equal(t, fetchedUser["thirdParty"].(map[string]interface{})["userId"], userRef.ThirdParty.UserID)

	userRef = nil
	assert.Nil(t, userRef)

	formFields = map[string]interface{}{
		"thirdPartyId": "custom",
		"redirectURIInfo": map[string]interface{}{
			"redirectURIOnProviderDashboard": testServer.URL + "/callback",
			"redirectURIQueryParams": map[string]interface{}{
				"code": "abcdefghj",
			},
		},
	}

	postBody, err = json.Marshal(formFields)
	if err != nil {
		t.Error(err.Error())
	}

	resp, err = http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	signInResponse := *unittesting.HttpResponseToConsumableInformation(resp.Body)
	fetchedUserFromSignIn := signInResponse["user"].(map[string]interface{})

	assert.NotNil(t, userRef)
	assert.False(t, newUser)
	assert.Equal(t, fetchedUserFromSignIn["email"], *userRef.Email)
	assert.Equal(t, fetchedUserFromSignIn["id"], userRef.ID)
	assert.Equal(t, fetchedUserFromSignIn["thirdParty"].(map[string]interface{})["id"], userRef.ThirdParty.ID)
	assert.Equal(t, fetchedUserFromSignIn["thirdParty"].(map[string]interface{})["userId"], userRef.ThirdParty.UserID)
}
