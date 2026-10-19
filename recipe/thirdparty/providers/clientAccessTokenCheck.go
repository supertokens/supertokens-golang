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

package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

var errAccessTokenOfAnotherApp = errors.New("Access token does not belong to your application")
var errInvalidAccessToken = errors.New("Invalid access token")

// clientAccessTokenCheck checks with the provider that the access token was issued to one of the given client ids.
type clientAccessTokenCheck func(accessToken string, clientConfig tpmodels.ProviderConfigForClientType, clientIDs map[string]bool) error

// addClientAccessTokenCheck makes GetUserInfo check access tokens that the client sent (in oAuthTokens) with the
// provider, before the user info is fetched with them.
//
// The check is skipped:
//   - if ValidateAccessToken is set in the provider config (GetUserInfo runs that instead, as before)
//   - for the access token the SDK got itself, by exchanging an authorisation code with the provider. This works
//     because a new provider instance is created for every request.
func addClientAccessTokenCheck(impl *tpmodels.TypeProvider, providerConfig tpmodels.ProviderConfig, check clientAccessTokenCheck) {
	var accessTokenFromCodeExchange *string

	oExchangeAuthCodeForOAuthTokens := impl.ExchangeAuthCodeForOAuthTokens
	impl.ExchangeAuthCodeForOAuthTokens = func(redirectURIInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
		oAuthTokens, err := oExchangeAuthCodeForOAuthTokens(redirectURIInfo, userContext)
		if err == nil {
			if accessToken, ok := oAuthTokens["access_token"].(string); ok {
				accessTokenFromCodeExchange = &accessToken
			}
		}
		return oAuthTokens, err
	}

	oGetUserInfo := impl.GetUserInfo
	impl.GetUserInfo = func(oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
		accessToken, ok := oAuthTokens["access_token"].(string)
		isFromCodeExchange := accessTokenFromCodeExchange != nil && *accessTokenFromCodeExchange == accessToken
		if ok && impl.Config.ValidateAccessToken == nil && !isFromCodeExchange {
			err := check(accessToken, impl.Config, getClientIDsOfProvider(providerConfig, impl.Config))
			if err != nil {
				return tpmodels.TypeUserInfo{}, err
			}
		}
		return oGetUserInfo(oAuthTokens, userContext)
	}
}

// getClientIDsOfProvider returns the client ids of every client of the provider (for this tenant)
func getClientIDsOfProvider(providerConfig tpmodels.ProviderConfig, clientConfig tpmodels.ProviderConfigForClientType) map[string]bool {
	clientIDs := map[string]bool{}
	for _, client := range providerConfig.Clients {
		if client.ClientID != "" {
			clientIDs[getActualClientIdFromDevelopmentClientId(client.ClientID)] = true
		}
	}
	if clientConfig.ClientID != "" {
		clientIDs[getActualClientIdFromDevelopmentClientId(clientConfig.ClientID)] = true
	}
	return clientIDs
}

func isOneOfClientIDs(value interface{}, clientIDs map[string]bool) bool {
	switch v := value.(type) {
	case string:
		return clientIDs[v]
	case json.Number:
		return clientIDs[v.String()]
	case float64:
		return clientIDs[fmt.Sprintf("%.0f", v)]
	default:
		return false
	}
}

// doAccessTokenCheckRequest sends a request to a provider's token check API. It doesn't log the query params, the
// headers or the response body, since they contain the access token (and for Facebook, the client secret).
// The returned status is -1 if the request itself failed.
func doAccessTokenCheckRequest(method string, requestURL string, queryParams map[string]string, formParams map[string]string, headers map[string]string) (map[string]interface{}, int, error) {
	supertokens.LogDebugMessage(fmt.Sprintf("%s request to %s to check an access token sent by the client", method, requestURL))

	parsedURL, err := url.Parse(requestURL)
	if err != nil {
		return nil, -1, err
	}
	if queryParams != nil {
		query := url.Values{}
		for key, value := range queryParams {
			query.Set(key, value)
		}
		parsedURL.RawQuery = query.Encode()
	}

	var req *http.Request
	if formParams != nil {
		form := url.Values{}
		for key, value := range formParams {
			form.Set(key, value)
		}
		req, err = http.NewRequest(method, parsedURL.String(), strings.NewReader(form.Encode()))
		if err == nil {
			req.Header.Set("content-type", "application/x-www-form-urlencoded")
		}
	} else {
		req, err = http.NewRequest(method, parsedURL.String(), nil)
	}
	if err != nil {
		return nil, -1, err
	}
	req.Header.Set("accept", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, -1, err
	}
	defer resp.Body.Close()

	supertokens.LogDebugMessage(fmt.Sprintf("Received response with status %d from %s", resp.StatusCode, requestURL))

	if resp.StatusCode >= 400 {
		return nil, resp.StatusCode, nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	decoder := json.NewDecoder(strings.NewReader(string(body)))
	decoder.UseNumber()
	var result map[string]interface{}
	if err := decoder.Decode(&result); err != nil {
		return nil, resp.StatusCode, nil
	}
	return result, resp.StatusCode, nil
}

func failedAccessTokenCheckError(providerName string, status int, err error, invalidTokenStatuses ...int) error {
	if status == -1 {
		// The request itself failed (e.g. a network error), so this says nothing about the token
		return err
	}
	for _, invalidTokenStatus := range invalidTokenStatuses {
		if status == invalidTokenStatus {
			return errInvalidAccessToken
		}
	}
	return fmt.Errorf("Could not check the access token with %s (status %d)", providerName, status)
}

// checkGoogleAccessToken uses Google's token info API, which returns the client the token is for (aud) and the
// client that requested it (azp). They can differ, e.g. for tokens requested on Android or iOS for another client
// of the same Google Cloud project.
func checkGoogleAccessToken(accessToken string, clientConfig tpmodels.ProviderConfigForClientType, clientIDs map[string]bool) error {
	resp, status, err := doAccessTokenCheckRequest("POST", "https://oauth2.googleapis.com/tokeninfo", nil, map[string]string{
		"access_token": accessToken,
	}, nil)
	if err != nil || status != 200 || resp == nil {
		return failedAccessTokenCheckError("Google", status, err, 400, 401)
	}
	if !isOneOfClientIDs(resp["aud"], clientIDs) && !isOneOfClientIDs(resp["azp"], clientIDs) {
		return errAccessTokenOfAnotherApp
	}
	return nil
}

// checkFacebookAccessToken uses Facebook's debug token API, which needs an app access token, made from the client id
// and the client secret.
func checkFacebookAccessToken(accessToken string, clientConfig tpmodels.ProviderConfigForClientType, clientIDs map[string]bool) error {
	if clientConfig.ClientSecret == "" {
		return errors.New("The Facebook provider needs the ClientSecret of the client to check access tokens sent in oAuthTokens. Please add it to the provider config, or set a custom ValidateAccessToken.")
	}
	resp, status, err := doAccessTokenCheckRequest("GET", "https://graph.facebook.com/debug_token", map[string]string{
		"input_token":  accessToken,
		"access_token": getActualClientIdFromDevelopmentClientId(clientConfig.ClientID) + "|" + clientConfig.ClientSecret,
	}, nil, nil)
	if status == -1 {
		return err
	}
	data, dataOk := resp["data"].(map[string]interface{})
	if err != nil || status != 200 || !dataOk {
		return fmt.Errorf("Could not check the access token with Facebook (status %d). Please check the ClientID and ClientSecret of the Facebook provider.", status)
	}
	if isValid, ok := data["is_valid"].(bool); !ok || !isValid {
		return errInvalidAccessToken
	}
	if !isOneOfClientIDs(data["app_id"], clientIDs) {
		return errAccessTokenOfAnotherApp
	}
	return nil
}

// checkGitlabAccessToken uses GitLab's token info API, which returns the application (client id) the token was
// issued to.
func checkGitlabAccessToken(accessToken string, clientConfig tpmodels.ProviderConfigForClientType, clientIDs map[string]bool) error {
	resp, status, err := doAccessTokenCheckRequest("GET", getGitlabBaseURL(clientConfig)+"/oauth/token/info", nil, nil, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	if err != nil || status != 200 || resp == nil {
		return failedAccessTokenCheckError("GitLab", status, err, 401)
	}
	application, _ := resp["application"].(map[string]interface{})
	if application == nil || !isOneOfClientIDs(application["uid"], clientIDs) {
		return errAccessTokenOfAnotherApp
	}
	return nil
}

func getGitlabBaseURL(clientConfig tpmodels.ProviderConfigForClientType) string {
	wellKnownSuffix := "/.well-known/openid-configuration"
	if strings.HasSuffix(clientConfig.OIDCDiscoveryEndpoint, wellKnownSuffix) {
		return strings.TrimSuffix(clientConfig.OIDCDiscoveryEndpoint, wellKnownSuffix)
	}
	return "https://gitlab.com"
}

// checkDiscordAccessToken uses Discord's current authorization API, which returns the application the token was
// issued to. Discord's client id is the application's id.
func checkDiscordAccessToken(accessToken string, clientConfig tpmodels.ProviderConfigForClientType, clientIDs map[string]bool) error {
	resp, status, err := doAccessTokenCheckRequest("GET", "https://discord.com/api/oauth2/@me", nil, nil, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	if err != nil || status != 200 || resp == nil {
		return failedAccessTokenCheckError("Discord", status, err, 401)
	}
	application, _ := resp["application"].(map[string]interface{})
	if application == nil || !isOneOfClientIDs(application["id"], clientIDs) {
		return errAccessTokenOfAnotherApp
	}
	return nil
}
