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

package unittesting

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"gopkg.in/h2non/gock.v1"
)

const testAppPrefix = "go-test-"

var licenseKeySet = false

// coreURL returns the base URL for the running core instance.
func coreURL() string {
	host := os.Getenv("SUPERTOKENS_CORE_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("SUPERTOKENS_CORE_PORT")
	if port == "" {
		port = "3567"
	}
	return fmt.Sprintf("http://%s:%s", host, port)
}

func ensureLicenseKey(connectionURI string) {
	licenseKey := os.Getenv("SUPERTOKENS_LICENSE_KEY")
	if licenseKey == "" {
		// Use the hardcoded test key for multitenancy
		licenseKey = "ijaleljUd2kU9XXWLiqFYv5br8nutTxbyBqWypQdv2N-BocoNriPrnYQd0NXPm8rVkeEocN9ayq0B7c3Pv-BTBIhAZSclXMlgyfXtlwAOJk=9BfESEleW6LyTov47dXu"
	}

	jsonData, _ := json.Marshal(map[string]interface{}{
		"licenseKey": licenseKey,
	})
	req, _ := http.NewRequest("PUT", connectionURI+"/ee/license", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(fmt.Sprintf("ensureLicenseKey failed: %s", err))
	}
	defer resp.Body.Close()
}

// SetUpST is now a no-op — app-based testing doesn't need per-test setup.
// Kept for backward compatibility with BeforeEach/AfterEach patterns.
func SetUpST() {}

// CleanST is now a no-op — app cleanup is handled by RemoveCoreApp.
func CleanST() {}

// StartUpST creates a core application and returns the connection URI.
// The host and port parameters are ignored — the shared core is used.
func StartUpST(host string, port string) string {
	return CreateCoreApp(nil)
}

// StartUpSTWithMultitenancy creates a core application with multitenancy support.
// The host and port parameters are ignored.
func StartUpSTWithMultitenancy(host string, port string) string {
	return CreateCoreApp(nil)
}

// CreateCoreApp creates a new isolated application in the running core.
// coreConfig is an optional map of core configuration overrides.
// Returns the connection URI for the created app.
func CreateCoreApp(coreConfig map[string]interface{}) string {
	base := coreURL()

	if !licenseKeySet {
		ensureLicenseKey(base)
		licenseKeySet = true
	}

	appId := testAppPrefix + uuid.New().String()

	if coreConfig == nil {
		coreConfig = map[string]interface{}{}
	}

	body, _ := json.Marshal(map[string]interface{}{
		"appId":      appId,
		"coreConfig": coreConfig,
	})

	req, _ := http.NewRequest("PUT", base+"/recipe/multitenancy/app/v2", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(fmt.Sprintf("CreateCoreApp failed: %s", err))
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if result["status"] != "OK" {
		respBytes, _ := json.Marshal(result)
		panic(fmt.Sprintf("CreateCoreApp failed: %s", string(respBytes)))
	}

	connectionURI := fmt.Sprintf("%s/appid-%s", base, appId)
	ensureLicenseKey(connectionURI)

	return connectionURI
}

// RemoveCoreApp deletes an application from the running core.
func RemoveCoreApp(connectionURI string) {
	if connectionURI == "" {
		return
	}

	base, appId := parseConnectionURI(connectionURI)
	if appId == "" {
		return
	}

	body, _ := json.Marshal(map[string]interface{}{
		"appId": appId,
	})

	req, _ := http.NewRequest("POST", base+"/recipe/multitenancy/app/remove", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return // best-effort cleanup
	}
	defer resp.Body.Close()
}

// KillAllST cleans up all test applications created by this process.
func KillAllST() {
	CleanupAllCoreApps()
}

// CleanupAllCoreApps removes all test applications matching our prefix.
func CleanupAllCoreApps() {
	base := coreURL()

	resp, err := http.Get(base + "/recipe/multitenancy/app/list/v2")
	if err != nil {
		return // core may not be running
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	apps, ok := result["apps"].([]interface{})
	if !ok {
		return
	}

	for _, app := range apps {
		appMap, ok := app.(map[string]interface{})
		if !ok {
			continue
		}
		appId, ok := appMap["appId"].(string)
		if !ok {
			continue
		}
		if strings.HasPrefix(appId, testAppPrefix) {
			RemoveCoreApp(fmt.Sprintf("%s/appid-%s", base, appId))
		}
	}
}

// SetKeyValueInConfig is kept for backward compatibility but should not be used
// with app-based testing. Use the coreConfig parameter of CreateCoreApp instead.
// This function panics to catch any remaining callers that need to be migrated.
func SetKeyValueInConfig(key string, value string) {
	panic("SetKeyValueInConfig is not supported with app-based testing. Pass coreConfig to CreateCoreApp instead.")
}

func parseConnectionURI(connectionURI string) (base string, appId string) {
	if strings.Contains(connectionURI, "appid-") {
		parts := strings.SplitN(connectionURI, "/appid-", 2)
		return parts[0], parts[1]
	}
	return connectionURI, ""
}

// MaxVersion returns max of v1 and v2
func MaxVersion(version1 string, version2 string) string {
	var splittedv1 = strings.Split(version1, ".")
	var splittedv2 = strings.Split(version2, ".")
	var minLength = len(splittedv1)
	if minLength > len(splittedv2) {
		minLength = len(splittedv2)
	}
	for i := 0; i < minLength; i++ {
		var v1, _ = strconv.Atoi(splittedv1[i])
		var v2, _ = strconv.Atoi(splittedv2[i])
		if v1 > v2 {
			return version1
		} else if v2 > v1 {
			return version2
		}
	}
	if len(splittedv1) >= len(splittedv2) {
		return version1
	}
	return version2
}

func ExtractInfoFromResponse(res *http.Response) map[string]string {
	antiCsrf := res.Header["Anti-Csrf"]
	cookies := res.Header["Set-Cookie"]
	var refreshToken string
	var refreshTokenExpiry string
	var refreshTokenDomain string
	var refreshTokenHttpOnly = "false"
	var accessToken string
	var accessTokenExpiry string
	var accessTokenDomain string
	var accessTokenHttpOnly = "false"

	for _, cookie := range cookies {
		if strings.Split(strings.Split(cookie, ";")[0], "=")[0] == "sRefreshToken" {
			refreshToken = strings.Split(strings.Split(cookie, ";")[0], "=")[1]
			if strings.Split(strings.Split(cookie, ";")[2], "=")[0] == " Expires" {
				refreshTokenExpiry = strings.Split(strings.Split(cookie, ";")[2], "=")[1]
			} else if strings.Split(strings.Split(cookie, ";")[2], "=")[0] == " expires" {
				refreshTokenExpiry = strings.Split(strings.Split(cookie, ";")[2], "=")[1]
			} else {
				refreshTokenExpiry = strings.Split(strings.Split(cookie, ";")[3], "=")[1]
			}
			for _, property := range strings.Split(cookie, ";") {
				if strings.HasPrefix(property, " Domain=") {
					refreshTokenDomain = strings.TrimPrefix(property, " Domain=")
				}
				if strings.Index(property, "HttpOnly") == 1 {
					refreshTokenHttpOnly = "true"
				}
			}
		} else if strings.Split(strings.Split(cookie, ";")[0], "=")[0] == "sAccessToken" {
			accessToken = strings.Split(strings.Split(cookie, ";")[0], "=")[1]
			if strings.Split(strings.Split(cookie, ";")[2], "=")[0] == " Expires" {
				accessTokenExpiry = strings.Split(strings.Split(cookie, ";")[2], "=")[1]
			} else if strings.Split(strings.Split(cookie, ";")[2], "=")[0] == " expires" {
				accessTokenExpiry = strings.Split(strings.Split(cookie, ";")[2], "=")[1]
			} else {
				accessTokenExpiry = strings.Split(strings.Split(cookie, ";")[3], "=")[1]
			}
			for _, property := range strings.Split(cookie, ";") {
				if strings.HasPrefix(property, " Domain=") {
					accessTokenDomain = strings.TrimPrefix(property, " Domain=")
				}
				if strings.Index(property, "HttpOnly") == 1 {
					accessTokenHttpOnly = "true"
				}
			}
		}
	}
	antiCsrfVal := ""
	if len(antiCsrf) > 0 {
		antiCsrfVal = antiCsrf[0]
	}
	frontToken := res.Header.Get("front-token")

	var refreshTokenFromHeader = res.Header.Get("st-refresh-token")
	var accessTokenFromHeader = res.Header.Get("st-access-token")

	refreshTokenFromAny := refreshToken
	if refreshTokenFromAny == "" {
		refreshTokenFromAny = refreshTokenFromHeader
	}
	accessTokenFromAny := accessToken
	if accessTokenFromAny == "" {
		accessTokenFromAny = accessTokenFromHeader
	}

	return map[string]string{
		"antiCsrf":               antiCsrfVal,
		"sAccessToken":           accessToken,
		"sRefreshToken":          refreshToken,
		"refreshTokenExpiry":     refreshTokenExpiry,
		"refreshTokenDomain":     refreshTokenDomain,
		"refreshTokenHttpOnly":   refreshTokenHttpOnly,
		"accessTokenExpiry":      accessTokenExpiry,
		"accessTokenDomain":      accessTokenDomain,
		"accessTokenHttpOnly":    accessTokenHttpOnly,
		"frontToken":             frontToken,
		"refreshTokenFromHeader": refreshTokenFromHeader,
		"accessTokenFromHeader":  accessTokenFromHeader,
		"refreshTokenFromAny":    refreshTokenFromAny,
		"accessTokenFromAny":     accessTokenFromAny,
	}
}

func ExtractInfoFromResponseForAuthModeTests(res *http.Response) map[string]string {
	antiCsrf := res.Header["Anti-Csrf"]
	cookies := res.Header["Set-Cookie"]
	var refreshToken = "-not-present-"
	var refreshTokenExpiry = "-not-present-"
	var refreshTokenDomain = "-not-present-"
	var refreshTokenHttpOnly = "false"
	var accessToken = "-not-present-"
	var accessTokenExpiry = "-not-present-"
	var accessTokenDomain = "-not-present-"
	var accessTokenHttpOnly = "false"

	for _, cookie := range cookies {
		if strings.Split(strings.Split(cookie, ";")[0], "=")[0] == "sRefreshToken" {
			refreshToken = strings.Split(strings.Split(cookie, ";")[0], "=")[1]
			if strings.Split(strings.Split(cookie, ";")[2], "=")[0] == " Expires" {
				refreshTokenExpiry = strings.Split(strings.Split(cookie, ";")[2], "=")[1]
			} else if strings.Split(strings.Split(cookie, ";")[2], "=")[0] == " expires" {
				refreshTokenExpiry = strings.Split(strings.Split(cookie, ";")[2], "=")[1]
			} else {
				refreshTokenExpiry = strings.Split(strings.Split(cookie, ";")[3], "=")[1]
			}
			for _, property := range strings.Split(cookie, ";") {
				if strings.Index(property, "HttpOnly") == 1 {
					refreshTokenHttpOnly = "true"
					break
				}
			}
		} else if strings.Split(strings.Split(cookie, ";")[0], "=")[0] == "sAccessToken" {
			accessToken = strings.Split(strings.Split(cookie, ";")[0], "=")[1]
			if strings.Split(strings.Split(cookie, ";")[2], "=")[0] == " Expires" {
				accessTokenExpiry = strings.Split(strings.Split(cookie, ";")[2], "=")[1]
			} else if strings.Split(strings.Split(cookie, ";")[2], "=")[0] == " expires" {
				accessTokenExpiry = strings.Split(strings.Split(cookie, ";")[2], "=")[1]
			} else {
				accessTokenExpiry = strings.Split(strings.Split(cookie, ";")[3], "=")[1]
			}
			for _, property := range strings.Split(cookie, ";") {
				if strings.Index(property, "HttpOnly") == 1 {
					accessTokenHttpOnly = "true"
					break
				}
			}
		}
	}
	antiCsrfVal := "-not-present-"
	if len(antiCsrf) > 0 {
		antiCsrfVal = antiCsrf[0]
	}
	frontToken := res.Header.Get("front-token")

	var refreshTokenFromHeader = "-not-present-"
	if len(res.Header.Values("st-refresh-token")) > 0 {
		refreshTokenFromHeader = res.Header.Get("st-refresh-token")
	}
	var accessTokenFromHeader = "-not-present-"
	if len(res.Header.Values("st-access-token")) > 0 {
		accessTokenFromHeader = res.Header.Get("st-access-token")
	}

	return map[string]string{
		"antiCsrf":               antiCsrfVal,
		"sAccessToken":           accessToken,
		"sRefreshToken":          refreshToken,
		"refreshTokenExpiry":     refreshTokenExpiry,
		"refreshTokenDomain":     refreshTokenDomain,
		"refreshTokenHttpOnly":   refreshTokenHttpOnly,
		"accessTokenExpiry":      accessTokenExpiry,
		"accessTokenDomain":      accessTokenDomain,
		"accessTokenHttpOnly":    accessTokenHttpOnly,
		"frontToken":             frontToken,
		"refreshTokenFromHeader": refreshTokenFromHeader,
		"accessTokenFromHeader":  accessTokenFromHeader,
	}
}

func ExtractInfoFromResponseWhenAntiCSRFisNone(res *http.Response) map[string]string {
	cookies := res.Header["Set-Cookie"]
	var refreshToken string
	var refreshTokenExpiry string
	var refreshTokenDomain string
	var refreshTokenHttpOnly = "false"
	var accessToken string
	var accessTokenExpiry string
	var accessTokenDomain string
	var accessTokenHttpOnly = "false"
	for _, cookie := range cookies {
		if strings.Split(strings.Split(cookie, ";")[0], "=")[0] == "sRefreshToken" {
			refreshToken = strings.Split(strings.Split(cookie, ";")[0], "=")[1]
			if strings.Split(strings.Split(cookie, ";")[2], "=")[0] == " Expires" {
				refreshTokenExpiry = strings.Split(strings.Split(cookie, ";")[2], "=")[1]
			} else if strings.Split(strings.Split(cookie, ";")[2], "=")[0] == " expires" {
				refreshTokenExpiry = strings.Split(strings.Split(cookie, ";")[2], "=")[1]
			} else {
				refreshTokenExpiry = strings.Split(strings.Split(cookie, ";")[3], "=")[1]
			}
			for _, property := range strings.Split(cookie, ";") {
				if strings.Index(property, "HttpOnly") == 1 {
					refreshTokenHttpOnly = "true"
					break
				}
			}
		} else if strings.Split(strings.Split(cookie, ";")[0], "=")[0] == "sAccessToken" {
			accessToken = strings.Split(strings.Split(cookie, ";")[0], "=")[1]
			if strings.Split(strings.Split(cookie, ";")[2], "=")[0] == " Expires" {
				accessTokenExpiry = strings.Split(strings.Split(cookie, ";")[2], "=")[1]
			} else if strings.Split(strings.Split(cookie, ";")[2], "=")[0] == " expires" {
				accessTokenExpiry = strings.Split(strings.Split(cookie, ";")[2], "=")[1]
			} else {
				accessTokenExpiry = strings.Split(strings.Split(cookie, ";")[3], "=")[1]
			}
			for _, property := range strings.Split(cookie, ";") {
				if strings.Index(property, "HttpOnly") == 1 {
					accessTokenHttpOnly = "true"
					break
				}
			}
		}
	}
	return map[string]string{
		"sAccessToken":         accessToken,
		"sRefreshToken":        refreshToken,
		"refreshTokenExpiry":   refreshTokenExpiry,
		"refreshTokenDomain":   refreshTokenDomain,
		"refreshTokenHttpOnly": refreshTokenHttpOnly,
		"accessTokenExpiry":    accessTokenExpiry,
		"accessTokenDomain":    accessTokenDomain,
		"accessTokenHttpOnly":  accessTokenHttpOnly,
	}
}

func SignupRequest(email string, password string, testUrl string) (*http.Response, error) {
	formFields := map[string][]map[string]string{
		"formFields": {
			{"id": "email", "value": email},
			{"id": "password", "value": password},
		},
	}
	postBody, err := json.Marshal(formFields)
	if err != nil {
		return nil, err
	}
	return http.Post(testUrl+"/auth/signup", "application/json", bytes.NewBuffer(postBody))
}

func SignupRequestWithTenantId(tenantId string, email string, password string, testUrl string) (*http.Response, error) {
	formFields := map[string][]map[string]string{
		"formFields": {
			{"id": "email", "value": email},
			{"id": "password", "value": password},
		},
	}
	postBody, err := json.Marshal(formFields)
	if err != nil {
		return nil, err
	}
	return http.Post(testUrl+fmt.Sprintf("/auth/%s/signup", tenantId), "application/json", bytes.NewBuffer(postBody))
}

func SignInRequest(email string, password string, testUrl string) (*http.Response, error) {
	formFields := map[string][]map[string]string{
		"formFields": {
			{"id": "email", "value": email},
			{"id": "password", "value": password},
		},
	}
	postBody, err := json.Marshal(formFields)
	if err != nil {
		return nil, err
	}
	return http.Post(testUrl+"/auth/signin", "application/json", bytes.NewBuffer(postBody))
}

func SignInRequestWithThirdpartyemailpasswordRid(email string, password string, testUrl string) (*http.Response, error) {
	formFields := map[string][]map[string]string{
		"formFields": {
			{"id": "email", "value": email},
			{"id": "password", "value": password},
		},
	}
	postBody, err := json.Marshal(formFields)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req, _ := http.NewRequest("POST", testUrl+"/auth/signin", bytes.NewBuffer(postBody))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("rid", "thirdpartyemailpassword")
	return client.Do(req)
}

func SignInRequestWithTenantId(tenantId string, email string, password string, testUrl string) (*http.Response, error) {
	formFields := map[string][]map[string]string{
		"formFields": {
			{"id": "email", "value": email},
			{"id": "password", "value": password},
		},
	}
	postBody, err := json.Marshal(formFields)
	if err != nil {
		return nil, err
	}
	return http.Post(testUrl+fmt.Sprintf("/auth/%s/signin", tenantId), "application/json", bytes.NewBuffer(postBody))
}

func EmailVerifyTokenRequest(testUrl string, userId string, accessToken string, antiCsrf string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, testUrl+"/auth/user/email/verify/token", bytes.NewBuffer([]byte(userId)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Cookie", "sAccessToken="+accessToken)
	req.Header.Add("anti-csrf", antiCsrf)
	return http.DefaultClient.Do(req)
}

func SignoutRequest(testUrl string, accessToken string, antiCsrf string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, testUrl+"/auth/signout", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Cookie", "sAccessToken="+accessToken)
	req.Header.Add("anti-csrf", antiCsrf)
	return http.DefaultClient.Do(req)
}

func SessionRefresh(testUrl string, refreshToken string, antiCsrf string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, testUrl+"/auth/session/refresh", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Cookie", "sRefreshToken="+refreshToken)
	req.Header.Add("anti-csrf", antiCsrf)
	return http.DefaultClient.Do(req)
}

func ReturnCustomProviderWithAuthRedirectParams() tpmodels.ProviderInput {
	return tpmodels.ProviderInput{
		Config: tpmodels.ProviderConfig{
			ThirdPartyId:          "custom",
			AuthorizationEndpoint: "https://test.com/oauth/auth",
			AuthorizationEndpointQueryParams: map[string]interface{}{
				"scope":     "test",
				"client_id": "supertokens",
			},
			TokenEndpoint: "https://test.com/oauth/token",
			Clients: []tpmodels.ProviderClientConfig{
				{ClientID: "supertokens", Scope: []string{"test"}},
			},
		},
		Override: func(originalImplementation *tpmodels.TypeProvider) *tpmodels.TypeProvider {
			oGetConfig := originalImplementation.GetConfigForClientType
			originalImplementation.GetConfigForClientType = func(clientType *string, userContext supertokens.UserContext) (tpmodels.ProviderConfigForClientType, error) {
				config, err := oGetConfig(clientType, userContext)
				if err != nil {
					return config, err
				}
				if _default, ok := (*userContext)["_default"].(map[string]interface{}); ok {
					if req, ok := _default["request"].(*http.Request); ok {
						config.AuthorizationEndpointQueryParams["dynamic"] = req.URL.Query().Get("dynamic")
					}
				}
				return config, nil
			}
			originalImplementation.GetUserInfo = func(oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
				return tpmodels.TypeUserInfo{
					ThirdPartyUserId: "user",
					Email:            &tpmodels.EmailStruct{ID: "email@test.com", IsVerified: true},
				}, nil
			}
			return originalImplementation
		},
	}
}

func ReturnCustomProviderWithoutAuthRedirectParams() tpmodels.ProviderInput {
	return tpmodels.ProviderInput{
		Config: tpmodels.ProviderConfig{
			ThirdPartyId:          "custom",
			AuthorizationEndpoint: "https://test.com/oauth/auth",
			TokenEndpoint:         "https://test.com/oauth/token",
			Clients: []tpmodels.ProviderClientConfig{
				{ClientID: "supertokens", Scope: []string{"test"}},
			},
		},
		Override: func(originalImplementation *tpmodels.TypeProvider) *tpmodels.TypeProvider {
			originalImplementation.GetUserInfo = func(oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
				return tpmodels.TypeUserInfo{
					ThirdPartyUserId: "user",
					Email:            &tpmodels.EmailStruct{ID: "email@test.com", IsVerified: true},
				}, nil
			}
			return originalImplementation
		},
	}
}

func SigninupCustomRequest(testServerUrl string, email string, id string) (*http.Response, error) {
	defer gock.OffAll()
	gock.New("https://test.com/").Post("oauth/token").Reply(200).JSON(map[string]interface{}{
		"email": email,
		"id":    id,
	})
	postData := map[string]interface{}{
		"thirdPartyId": "custom",
		"redirectURIInfo": map[string]interface{}{
			"redirectURIOnProviderDashboard": "http://localhost.org",
			"redirectURIQueryParams": map[string]interface{}{
				"code": "32432432",
			},
		},
	}
	postBody, err := json.Marshal(postData)
	if err != nil {
		return nil, err
	}
	gock.New(testServerUrl).EnableNetworking().Persist()
	gock.New("http://localhost:3567/").EnableNetworking().Persist()
	return http.Post(testServerUrl+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
}

func HttpResponseToConsumableInformation(body io.ReadCloser) *map[string]interface{} {
	dataInBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil
	}
	body.Close()
	var result map[string]interface{}
	err = json.Unmarshal(dataInBytes, &result)
	if err != nil {
		return nil
	}
	return &result
}

func GenerateRandomCode(size int) string {
	characters := "ABCDEFGHIJKLMNOPQRSTUVWXTZabcdefghiklmnopqrstuvwxyz"
	randomString := ""
	for i := 0; i < size; i++ {
		randomNumber := rand.Intn(len(characters))
		randomString += characters[randomNumber : randomNumber+1]
	}
	return randomString
}

func EmailVerificationTokenRequest(cookies []*http.Cookie, testUrl string) (*http.Response, error) {
	req, _ := http.NewRequest("POST", testUrl+"/auth/user/email/verify/token", nil)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	return http.DefaultClient.Do(req)
}

func PasswordResetTokenRequest(email string, testUrl string) (*http.Response, error) {
	formFields := map[string][]map[string]string{
		"formFields": {{"id": "email", "value": email}},
	}
	postBody, err := json.Marshal(formFields)
	if err != nil {
		return nil, err
	}
	return http.Post(testUrl+"/auth/user/password/reset/token", "application/json", bytes.NewBuffer(postBody))
}

func PasswordlessEmailLoginRequest(email string, testUrl string) (*http.Response, error) {
	body := map[string]string{"email": email}
	postBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return http.Post(testUrl+"/auth/signinup/code", "application/json", bytes.NewBuffer(postBody))
}

func PasswordlessPhoneLoginRequest(phone string, testUrl string) (*http.Response, error) {
	body := map[string]string{"phoneNumber": phone}
	postBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return http.Post(testUrl+"/auth/signinup/code", "application/json", bytes.NewBuffer(postBody))
}

func PasswordlessLoginResendRequest(deviceId string, preAuthSessionId string, testUrl string) (*http.Response, error) {
	body := map[string]interface{}{
		"deviceId":         deviceId,
		"preAuthSessionId": preAuthSessionId,
	}
	postBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return http.Post(testUrl+"/auth/signinup/code/resend", "application/json", bytes.NewBuffer(postBody))
}

func PasswordlessLoginWithCodeRequest(deviceId string, preAuthSessionId string, code string, testUrl string) (*http.Response, error) {
	body := map[string]string{
		"deviceId":         deviceId,
		"preAuthSessionId": preAuthSessionId,
		"userInputCode":    code,
	}
	postBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return http.Post(testUrl+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(postBody))
}

func GetRequestWithJSONResult(url string, cookies []*http.Cookie) (int, map[string]interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, nil, err
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	respObj := map[string]interface{}{}
	err = json.NewDecoder(res.Body).Decode(&respObj)
	if err != nil {
		return 0, nil, err
	}
	res.Body.Close()
	return res.StatusCode, respObj, nil
}

type InfoLogData struct {
	LastLine string
	Output   []string
}

func GetInfoLogData(t *testing.T, startWith string) InfoLogData {
	// With app-based testing on a shared core, we can't isolate per-test logs.
	// Tests that depend on this need to be reworked.
	t.Log("GetInfoLogData: not supported with app-based testing")
	return InfoLogData{}
}
