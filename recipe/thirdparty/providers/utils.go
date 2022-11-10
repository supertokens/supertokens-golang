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
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/derekstavis/go-qs"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

// Network utils
func doGetRequest(url string, queryParams map[string]interface{}, headers map[string]string) (interface{}, error) {
	if queryParams != nil {
		querystring, err := qs.Marshal(queryParams)
		if err != nil {
			return nil, err
		}
		url = url + "?" + querystring
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("GET request to %s resulted in %d status with body %s", url, resp.StatusCode, string(body))
	}
	return result, nil
}

func doPostRequest(url string, params map[string]interface{}, headers map[string]interface{}) (map[string]interface{}, error) {
	postBody, err := qs.Marshal(params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(postBody)))
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value.(string))
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("accept", "application/json") // few providers like github don't send back json response by default

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("POST request to %s resulted in %d status with body %s", url, resp.StatusCode, string(body))
	}

	return result, nil
}

// JWKS utils
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

// OIDC utils
var oidcInfoMap = map[string]map[string]interface{}{}
var oidcInfoMapLock = sync.Mutex{}

func getOIDCDiscoveryInfo(issuer string) (map[string]interface{}, error) {
	if oidcInfo, ok := oidcInfoMap[issuer]; ok {
		return oidcInfo, nil
	}

	oidcInfoMapLock.Lock()
	defer oidcInfoMapLock.Unlock()

	// Check again to see if it was added while we were waiting for the lock
	if oidcInfo, ok := oidcInfoMap[issuer]; ok {
		return oidcInfo, nil
	}

	oidcInfo, err := doGetRequest(issuer+"/.well-known/openid-configuration", nil, nil)
	if err != nil {
		return nil, err
	}
	oidcInfoMap[issuer] = oidcInfo.(map[string]interface{})
	return oidcInfoMap[issuer], nil
}

// User map utils
func accessField(obj interface{}, key string) interface{} {
	keyParts := strings.Split(key, ".")
	for _, k := range keyParts {
		obj = obj.(map[string]interface{})[k]
	}
	return obj
}

func getSupertokensUserInfoFromRawUserInfo(idField string, emailField string, emailVerifiedField string, from tpmodels.TypeFrom) func(rawUserInfoResponse tpmodels.TypeRawUserInfoFromProvider, userContext supertokens.UserContext) (tpmodels.TypeSupertokensUserInfo, error) {
	return func(rawUserInfoResponse tpmodels.TypeRawUserInfoFromProvider, userContext supertokens.UserContext) (tpmodels.TypeSupertokensUserInfo, error) {
		var rawUserInfo map[string]interface{}

		if from == tpmodels.FromIdTokenPayload {
			if rawUserInfoResponse.FromIdTokenPayload == nil {
				return tpmodels.TypeSupertokensUserInfo{}, errors.New("rawUserInfoResponse.FromIdToken is not available")
			}
			rawUserInfo = rawUserInfoResponse.FromIdTokenPayload
		} else {
			if rawUserInfoResponse.FromUserInfoAPI == nil {
				return tpmodels.TypeSupertokensUserInfo{}, errors.New("rawUserInfoResponse.FromAccessToken is not available")
			}
			rawUserInfo = rawUserInfoResponse.FromUserInfoAPI
		}
		result := tpmodels.TypeSupertokensUserInfo{}
		result.ThirdPartyUserId = fmt.Sprint(accessField(rawUserInfo, idField))
		result.EmailInfo = &tpmodels.EmailStruct{
			ID: fmt.Sprint(accessField(rawUserInfo, emailField)),
		}
		if emailVerified, ok := accessField(rawUserInfo, emailVerifiedField).(bool); ok {
			result.EmailInfo.IsVerified = emailVerified
		} else if emailVerified, ok := accessField(rawUserInfo, emailVerifiedField).(string); ok {
			result.EmailInfo.IsVerified = emailVerified == "true"
		}
		return result, nil
	}
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

// PKCE related functions
// Ref: https://github.com/nirasan/go-oauth-pkce-code-verifier/blob/master/verifier.go

func randomBytes(length int) ([]byte, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	const csLen = byte(len(charset))
	output := make([]byte, 0, length)
	for {
		buf := make([]byte, length)
		if _, err := io.ReadFull(rand.Reader, buf); err != nil {
			return nil, fmt.Errorf("failed to read random bytes: %v", err)
		}
		for _, b := range buf {
			// Avoid bias by using a value range that's a multiple of 62
			if b < (csLen * 4) {
				output = append(output, charset[b%csLen])

				if len(output) == length {
					return output, nil
				}
			}
		}
	}
}

func encode(msg []byte) string {
	encoded := base64.StdEncoding.EncodeToString(msg)
	encoded = strings.Replace(encoded, "+", "-", -1)
	encoded = strings.Replace(encoded, "/", "_", -1)
	encoded = strings.Replace(encoded, "=", "", -1)
	return encoded
}

func generateCodeChallengeS256(length int) (codeChallenge string, codeVerifier string, err error) {
	buf, err := randomBytes(length)
	if err != nil {
		return "", "", err
	}

	codeVerifier = encode(buf)
	h := sha256.New()
	h.Write([]byte(codeVerifier))
	codeChallenge = encode(h.Sum(nil))
	err = nil
	return
}
