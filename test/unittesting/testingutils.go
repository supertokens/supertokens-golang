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
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"gopkg.in/h2non/gock.v1"
)

// containerCounter tracks running containers for unique naming
var containerCounter int

// configDir is the temporary directory holding config.yaml for the current test
var configDir string

func getCoreImage() string {
	img := os.Getenv("SUPERTOKENS_CORE_IMAGE")
	if img != "" {
		return img
	}
	version := os.Getenv("SUPERTOKENS_CORE_VERSION")
	if version == "" {
		version = "master"
	}
	return "supertokens/supertokens-dev-postgresql:" + version
}

func SetUpST() {
	dir, err := os.MkdirTemp("", "st-config-*")
	if err != nil {
		panic(fmt.Sprintf("failed to create temp config dir: %s", err))
	}
	configDir = dir

	// Write a minimal default config
	// info_log_path: null → log to stdout (captured by docker logs)
	// error_log_path: null → log to stderr (captured by docker logs)
	defaultConfig := "core_config_version: 0\ninfo_log_path: null\nerror_log_path: null\n"
	err = os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte(defaultConfig), 0644)
	if err != nil {
		panic(fmt.Sprintf("failed to write config.yaml: %s", err))
	}
}

func StartUpST(host string, port string) string {
	containerCounter++
	containerName := fmt.Sprintf("supertokens-test-%d-%d", os.Getpid(), containerCounter)

	args := []string{
		"run", "-d",
		"--name", containerName,
		"--platform", "linux/amd64",
		"-p", fmt.Sprintf("%s:3567", port),
	}

	// Mount config if we have one
	if configDir != "" {
		configPath := filepath.Join(configDir, "config.yaml")
		args = append(args, "-v", fmt.Sprintf("%s:/usr/lib/supertokens/config.yaml", configPath))
	}

	args = append(args, getCoreImage(),
		"/usr/lib/supertokens/jre/bin/java",
		"-classpath", "/usr/lib/supertokens/core/*:/usr/lib/supertokens/plugin-interface/*:/usr/lib/supertokens/ee/*",
		"io.supertokens.Main", "/usr/lib/supertokens/", "DEV",
		fmt.Sprintf("host=%s", "0.0.0.0"),
		fmt.Sprintf("port=%s", "3567"),
		"test_mode",
	)

	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("could not start ST container: %s\nOutput: %s", err, string(output)))
	}

	// Wait for core to be ready
	startTime := time.Now()
	for time.Since(startTime) < 30*time.Second {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/hello", port))
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == 200 {
				return containerName
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	panic("could not start ST process")
}

func StartUpSTWithMultitenancy(host string, port string) string {
	containerName := StartUpST(host, port)

	const OPAQUE_KEY_WITH_MULTITENANCY_FEATURE = "ijaleljUd2kU9XXWLiqFYv5br8nutTxbyBqWypQdv2N-BocoNriPrnYQd0NXPm8rVkeEocN9ayq0B7c3Pv-BTBIhAZSclXMlgyfXtlwAOJk=9BfESEleW6LyTov47dXu"

	jsonData, err := json.Marshal(map[string]interface{}{
		"licenseKey": OPAQUE_KEY_WITH_MULTITENANCY_FEATURE,
	})
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("http://%s:%s/ee/license", host, port), bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	return containerName
}

func stopContainer(name string) {
	// Use rm -f to skip the graceful shutdown wait (docker stop has a 10s default timeout)
	exec.Command("docker", "rm", "-f", name).Run()
}

func CleanST() {
	if configDir != "" {
		os.RemoveAll(configDir)
		configDir = ""
	}
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

func KillAllST() {
	// Stop all containers matching our naming pattern
	cmd := exec.Command("docker", "ps", "-a", "--filter", fmt.Sprintf("name=supertokens-test-%d-", os.Getpid()), "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return
	}
	names := strings.TrimSpace(string(output))
	if names == "" {
		return
	}
	for _, name := range strings.Split(names, "\n") {
		name = strings.TrimSpace(name)
		if name != "" {
			stopContainer(name)
		}
	}
}

func SetKeyValueInConfig(key string, value string) {
	if configDir == "" {
		panic("SetKeyValueInConfig called before SetUpST")
	}
	pathToConfigYamlFile := filepath.Join(configDir, "config.yaml")
	f, err := os.OpenFile(pathToConfigYamlFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(key + ": " + value + "\n"); err != nil {
		panic(err)
	}
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

	// Cookie stuff
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

	// Header stuff
	var refreshTokenFromHeader string = res.Header.Get("st-refresh-token")
	var accessTokenFromHeader string = res.Header.Get("st-access-token")

	refreshTokenFromAny := refreshToken

	if refreshTokenFromAny == "" {
		refreshTokenFromAny = refreshTokenFromHeader
	}

	accessTokenFromAny := accessToken

	if accessTokenFromAny == "" {
		accessTokenFromAny = accessTokenFromHeader
	}

	return map[string]string{
		"antiCsrf":             antiCsrfVal,
		"sAccessToken":         accessToken,
		"sRefreshToken":        refreshToken,
		"refreshTokenExpiry":   refreshTokenExpiry,
		"refreshTokenDomain":   refreshTokenDomain,
		"refreshTokenHttpOnly": refreshTokenHttpOnly,
		"accessTokenExpiry":    accessTokenExpiry,
		"accessTokenDomain":    accessTokenDomain,
		"accessTokenHttpOnly":  accessTokenHttpOnly,
		"frontToken":           frontToken,

		"refreshTokenFromHeader": refreshTokenFromHeader,
		"accessTokenFromHeader":  accessTokenFromHeader,
		"refreshTokenFromAny":    refreshTokenFromAny,
		"accessTokenFromAny":     accessTokenFromAny,
	}
}

func ExtractInfoFromResponseForAuthModeTests(res *http.Response) map[string]string {
	antiCsrf := res.Header["Anti-Csrf"]
	cookies := res.Header["Set-Cookie"]
	var refreshToken string = "-not-present-"
	var refreshTokenExpiry string = "-not-present-"
	var refreshTokenDomain string = "-not-present-"
	var refreshTokenHttpOnly = "false"
	var accessToken string = "-not-present-"
	var accessTokenExpiry string = "-not-present-"
	var accessTokenDomain string = "-not-present-"
	var accessTokenHttpOnly = "false"

	// Cookie stuff
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

	// Header stuff
	var refreshTokenFromHeader string = "-not-present-"
	if len(res.Header.Values("st-refresh-token")) > 0 {
		refreshTokenFromHeader = res.Header.Get("st-refresh-token")
	}
	var accessTokenFromHeader string = "-not-present-"
	if len(res.Header.Values("st-access-token")) > 0 {
		accessTokenFromHeader = res.Header.Get("st-access-token")
	}

	return map[string]string{
		"antiCsrf":             antiCsrfVal,
		"sAccessToken":         accessToken,
		"sRefreshToken":        refreshToken,
		"refreshTokenExpiry":   refreshTokenExpiry,
		"refreshTokenDomain":   refreshTokenDomain,
		"refreshTokenHttpOnly": refreshTokenHttpOnly,
		"accessTokenExpiry":    accessTokenExpiry,
		"accessTokenDomain":    accessTokenDomain,
		"accessTokenHttpOnly":  accessTokenHttpOnly,
		"frontToken":           frontToken,

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
			{
				"id":    "email",
				"value": email,
			},
			{
				"id":    "password",
				"value": password,
			},
		},
	}

	postBody, err := json.Marshal(formFields)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	resp, err := http.Post(testUrl+"/auth/signup", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return resp, nil
}

func SignupRequestWithTenantId(tenantId string, email string, password string, testUrl string) (*http.Response, error) {
	formFields := map[string][]map[string]string{
		"formFields": {
			{
				"id":    "email",
				"value": email,
			},
			{
				"id":    "password",
				"value": password,
			},
		},
	}

	postBody, err := json.Marshal(formFields)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	resp, err := http.Post(testUrl+fmt.Sprintf("/auth/%s/signup", tenantId), "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return resp, nil
}

func SignInRequest(email string, password string, testUrl string) (*http.Response, error) {
	formFields := map[string][]map[string]string{
		"formFields": {
			{
				"id":    "email",
				"value": email,
			},
			{
				"id":    "password",
				"value": password,
			},
		},
	}

	postBody, err := json.Marshal(formFields)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	resp, err := http.Post(testUrl+"/auth/signin", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return resp, nil
}

func SignInRequestWithThirdpartyemailpasswordRid(email string, password string, testUrl string) (*http.Response, error) {
	formFields := map[string][]map[string]string{
		"formFields": {
			{
				"id":    "email",
				"value": email,
			},
			{
				"id":    "password",
				"value": password,
			},
		},
	}

	postBody, err := json.Marshal(formFields)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	client := &http.Client{}
	req, _ := http.NewRequest("POST", testUrl+"/auth/signin", bytes.NewBuffer(postBody))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("rid", "thirdpartyemailpassword")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return resp, nil
}

func SignInRequestWithTenantId(tenantId string, email string, password string, testUrl string) (*http.Response, error) {
	formFields := map[string][]map[string]string{
		"formFields": {
			{
				"id":    "email",
				"value": email,
			},
			{
				"id":    "password",
				"value": password,
			},
		},
	}

	postBody, err := json.Marshal(formFields)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	resp, err := http.Post(testUrl+fmt.Sprintf("/auth/%s/signin", tenantId), "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return resp, nil
}

func EmailVerifyTokenRequest(testUrl string, userId string, accessToken string, antiCsrf string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, testUrl+"/auth/user/email/verify/token", bytes.NewBuffer([]byte(userId)))
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Cookie", "sAccessToken="+accessToken)
	req.Header.Add("anti-csrf", antiCsrf)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return resp, nil
}

func SignoutRequest(testUrl string, accessToken string, antiCsrf string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, testUrl+"/auth/signout", nil)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Cookie", "sAccessToken="+accessToken)
	req.Header.Add("anti-csrf", antiCsrf)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return resp, nil
}

func SessionRefresh(testUrl string, refreshToken string, antiCsrf string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, testUrl+"/auth/session/refresh", nil)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Cookie", "sRefreshToken="+refreshToken)
	req.Header.Add("anti-csrf", antiCsrf)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return resp, nil
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
				{
					ClientID: "supertokens",
					Scope:    []string{"test"},
				},
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
					Email: &tpmodels.EmailStruct{
						ID:         "email@test.com",
						IsVerified: true,
					},
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
				{
					ClientID: "supertokens",
					Scope:    []string{"test"},
				},
			},
		},
		Override: func(originalImplementation *tpmodels.TypeProvider) *tpmodels.TypeProvider {
			originalImplementation.GetUserInfo = func(oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
				return tpmodels.TypeUserInfo{
					ThirdPartyUserId: "user",
					Email: &tpmodels.EmailStruct{
						ID:         "email@test.com",
						IsVerified: true,
					},
				}, nil
			}
			return originalImplementation
		},
	}
}

func SigninupCustomRequest(testServerUrl string, email string, id string) (*http.Response, error) {
	defer gock.OffAll()
	gock.New("https://test.com/").
		Post("oauth/token").
		Reply(200).
		JSON(map[string]interface{}{
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
	gock.New("http://localhost:8080/").EnableNetworking().Persist()

	resp, err := http.Post(testServerUrl+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		return nil, err
	}
	return resp, nil
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

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return resp, nil
}

func PasswordResetTokenRequest(email string, testUrl string) (*http.Response, error) {
	formFields := map[string][]map[string]string{
		"formFields": {
			{
				"id":    "email",
				"value": email,
			},
		},
	}

	postBody, err := json.Marshal(formFields)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	resp, err := http.Post(testUrl+"/auth/user/password/reset/token", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return resp, nil
}

func PasswordlessEmailLoginRequest(email string, testUrl string) (*http.Response, error) {
	body := map[string]string{
		"email": email,
	}

	postBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	resp, err := http.Post(testUrl+"/auth/signinup/code", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return resp, nil
}

func PasswordlessPhoneLoginRequest(phone string, testUrl string) (*http.Response, error) {
	body := map[string]string{
		"phoneNumber": phone,
	}

	postBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	resp, err := http.Post(testUrl+"/auth/signinup/code", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return resp, nil
}

func PasswordlessLoginResendRequest(deviceId string, preAuthSessionId string, testUrl string) (*http.Response, error) {
	body := map[string]interface{}{
		"deviceId":         deviceId,
		"preAuthSessionId": preAuthSessionId,
	}

	postBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	resp, err := http.Post(testUrl+"/auth/signinup/code/resend", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return resp, nil
}

func PasswordlessLoginWithCodeRequest(deviceId string, preAuthSessionId string, code string, testUrl string) (*http.Response, error) {
	body := map[string]string{
		"deviceId":         deviceId,
		"preAuthSessionId": preAuthSessionId,
		"userInputCode":    code,
	}

	postBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	resp, err := http.Post(testUrl+"/auth/signinup/code/consume", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return resp, nil
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

func getRunningContainerName() string {
	cmd := exec.Command("docker", "ps", "--filter", fmt.Sprintf("name=supertokens-test-%d-", os.Getpid()), "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	names := strings.TrimSpace(string(output))
	if names == "" {
		return ""
	}
	// Return the first running container
	return strings.Split(names, "\n")[0]
}

func GetInfoLogData(t *testing.T, startWith string) InfoLogData {
	containerName := getRunningContainerName()
	if containerName == "" {
		t.Log("GetInfoLogData: no running container found")
		return InfoLogData{}
	}

	cmd := exec.Command("docker", "logs", containerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("GetInfoLogData: docker logs failed: %s", err)
		return InfoLogData{}
	}

	lines := strings.Split(string(output), "\n")
	var lastLine string
	var resultLines []string

	shouldRecord := startWith == ""

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if !shouldRecord && startWith != "" && strings.Contains(line, startWith) {
			shouldRecord = true
			continue
		}

		if shouldRecord {
			resultLines = append(resultLines, line)
		}

		lastLine = line
	}

	return InfoLogData{
		LastLine: lastLine,
		Output:   resultLines,
	}
}
