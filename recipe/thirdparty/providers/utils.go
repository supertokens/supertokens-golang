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

package providers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/derekstavis/go-qs"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func doGetRequest(req *http.Request) (interface{}, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("Provider API returned response with status `%s` and body `%s`", resp.Status, string(body)))
	}

	var result interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func postRequest(url string, params map[string]string) (map[string]interface{}, error) {
	querystring, err := getParamString(params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(querystring)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("accept", "application/json") // few providers like github don't send back json response by default

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getParamString(paramsMap map[string]string) (string, error) {
	params := map[string]interface{}{}
	for key, value := range paramsMap {
		params[key] = value
	}
	return qs.Marshal(params)
}

var jwksKeys = map[string]*keyfunc.JWKS{}
var jwksKeysLock = sync.Mutex{}

func getJWKSFromURL(url string) (*keyfunc.JWKS, error) {
	if jwks, ok := jwksKeys[url]; ok {
		return jwks, nil
	}

	jwksKeysLock.Lock()
	defer jwksKeysLock.Unlock()

	// Check again to see if it was added while we were waiting for the lock
	if jwks, ok := jwksKeys[url]; ok {
		return jwks, nil
	}

	options := keyfunc.Options{
		RefreshInterval: time.Hour,
	}
	jwks, err := keyfunc.Get(url, options)
	if err != nil {
		return nil, err
	}
	jwksKeys[url] = jwks
	return jwks, nil
}

var DevOauthClientIds = [...]string{
	"1060725074195-kmeum4crr01uirfl2op9kd5acmi9jutn.apps.googleusercontent.com", // google
	"467101b197249757c71f", // github
}

const (
	DevOauthAuthorisationUrl = "https://supertokens.io/dev/oauth/redirect-to-provider"
	DevOauthRedirectUrl      = "https://supertokens.io/dev/oauth/redirect-to-app"
	DevKeyIdentifier         = "4398792-"
)

func isUsingDevelopmentClientId(clientId string) bool {
	if strings.HasPrefix(clientId, DevKeyIdentifier) {
		return true
	} else {
		for _, devClientId := range DevOauthClientIds {
			if devClientId == clientId {
				return true
			}
		}
		return false
	}
}

func getActualClientIdFromDevelopmentClientId(clientId string) string {
	if strings.HasPrefix(clientId, DevKeyIdentifier) {
		return strings.Split(clientId, DevKeyIdentifier)[1]
	}
	return clientId
}

func getAuthRedirectForDev(clientId string, url string, queryParams map[string]interface{}) (string, map[string]interface{}, error) {
	if !isUsingDevelopmentClientId(clientId) {
		return url, queryParams, nil
	}
	queryParams["actual_redirect_uri"] = url
	return DevOauthAuthorisationUrl, queryParams, nil
}

func checkDevAndGetRedirectURI(clientId string, redirectURI string, userContext supertokens.UserContext) string {
	if isUsingDevelopmentClientId(clientId) {
		return DevOauthRedirectUrl
	}

	return redirectURI
}
