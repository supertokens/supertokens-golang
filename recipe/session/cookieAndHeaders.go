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

package session

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
	"time"

	sessionError "github.com/supertokens/supertokens-golang/recipe/session/errors"

	"github.com/supertokens/supertokens-golang/supertokens"

	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
)

const (
	authorizationHeaderKey = "authorization"
	accessTokenCookieKey   = "sAccessToken"
	accessTokenHeaderKey   = "st-access-token"
	refreshTokenCookieKey  = "sRefreshToken"
	refreshTokenHeaderKey  = "st-refresh-token"

	antiCsrfHeaderKey = "anti-csrf"
	ridHeaderKey      = "rid"

	frontTokenHeaderKey = "front-token"

	frontendSDKNameHeaderKey    = "supertokens-sdk-name"
	frontendSDKVersionHeaderKey = "supertokens-sdk-version"

	authModeHeaderKey = "st-auth-mode"
)

type TokenInfo struct {
	Uid string      `json:"uid"`
	Ate uint64      `json:"ate"`
	Up  interface{} `json:"up"`
}

func ClearSessionFromAllTokenTransferMethods(config sessmodels.TypeNormalisedInput, req *http.Request, res http.ResponseWriter, userContext supertokens.UserContext) error {
	// We are clearing the session in all transfermethods to be sure to override cookies in case they have been already added to the response.
	// This is done to handle the following use-case:
	// If the app overrides signInPOST to check the ban status of the user after the original implementation and throwing an UNAUTHORISED error
	// In this case: the SDK has attached cookies to the response, but none was sent with the request
	// We can't know which to clear since we can't reliably query or remove the set-cookie header added to the response (causes issues in some frameworks, i.e.: hapi)
	// The safe solution in this case is to overwrite all the response cookies/headers with an empty value, which is what we are doing here
	for _, transferMethod := range AvailableTokenTransferMethods {
		err := ClearSession(config, res, transferMethod, req, userContext)
		if err != nil {
			return err
		}
	}
	return nil
}

func ClearSession(config sessmodels.TypeNormalisedInput, res http.ResponseWriter, transferMethod sessmodels.TokenTransferMethod, request *http.Request, userContext supertokens.UserContext) error {
	// If we can be specific about which transferMethod we want to clear, there is no reason to clear the other ones
	tokenTypes := []sessmodels.TokenType{sessmodels.AccessToken, sessmodels.RefreshToken}
	for _, tokenType := range tokenTypes {
		err := setToken(config, res, tokenType, "", 0, transferMethod, request, userContext)
		if err != nil {
			return err
		}
	}

	res.Header().Del(antiCsrfHeaderKey)
	// This can be added multiple times in some cases, but that should be OK
	setHeader(res, frontTokenHeaderKey, "remove", false)
	setHeader(res, "Access-Control-Expose-Headers", frontTokenHeaderKey, true)
	return nil
}

func GetAntiCsrfTokenFromHeaders(req *http.Request) *string {
	return getHeader(req, antiCsrfHeaderKey)
}

func setAntiCsrfTokenInHeaders(res http.ResponseWriter, antiCsrfToken string) {
	setHeader(res, antiCsrfHeaderKey, antiCsrfToken, false)
	setHeader(res, "Access-Control-Expose-Headers", antiCsrfHeaderKey, true)
}

func BuildFrontToken(userId string, atExpiry uint64, jwtPayload interface{}) string {
	tokenInfo := &TokenInfo{
		Uid: userId,
		Ate: atExpiry,
		Up:  jwtPayload,
	}

	parsed, _ := json.Marshal(tokenInfo)
	data := []byte(parsed)

	return base64.StdEncoding.EncodeToString(data)
}

func setFrontTokenInHeaders(res http.ResponseWriter, frontToken string) {
	setHeader(res, frontTokenHeaderKey, frontToken, false)
	setHeader(res, "Access-Control-Expose-Headers", frontTokenHeaderKey, true)
}

func GetCORSAllowedHeaders() []string {
	return []string{
		antiCsrfHeaderKey, ridHeaderKey, authorizationHeaderKey, authModeHeaderKey,
	}
}

func getCookieNameFromTokenType(tokenType sessmodels.TokenType) (string, error) {
	if tokenType == sessmodels.AccessToken {
		return accessTokenCookieKey, nil
	}
	if tokenType == sessmodels.RefreshToken {
		return refreshTokenCookieKey, nil
	}
	return "", errors.New("Unknown token type, should never happen.")
}

func getResponseHeaderNameForTokenType(tokenType sessmodels.TokenType) (string, error) {
	if tokenType == sessmodels.AccessToken {
		return accessTokenHeaderKey, nil
	}
	if tokenType == sessmodels.RefreshToken {
		return refreshTokenHeaderKey, nil
	}
	return "", errors.New("Unknown token type, should never happen.")
}

func GetToken(req *http.Request, tokenType sessmodels.TokenType, transferMethod sessmodels.TokenTransferMethod) (*string, error) {
	if transferMethod == sessmodels.CookieTransferMethod {
		cookieName, err := getCookieNameFromTokenType(tokenType)
		if err != nil {
			return nil, err
		}
		return GetCookieValue(req, cookieName), nil
	} else if transferMethod == sessmodels.HeaderTransferMethod {
		headerValue := getHeader(req, authorizationHeaderKey)
		if headerValue == nil || !strings.HasPrefix(*headerValue, "Bearer ") {
			return nil, nil
		}

		token := strings.TrimSpace(strings.ReplaceAll(*headerValue, "Bearer ", ""))
		return &token, nil
	}
	return nil, errors.New("Should never happen")
}

func setToken(config sessmodels.TypeNormalisedInput, res http.ResponseWriter, tokenType sessmodels.TokenType, value string, expires uint64, transferMethod sessmodels.TokenTransferMethod, request *http.Request, userContext supertokens.UserContext) error {
	supertokens.LogDebugMessage(fmt.Sprint("setToken: Setting ", tokenType, " token as ", transferMethod))
	if transferMethod == sessmodels.CookieTransferMethod {
		cookieName, err := getCookieNameFromTokenType(tokenType)
		if err != nil {
			return err
		}
		pathType := ""
		if tokenType == sessmodels.AccessToken {
			pathType = "accessTokenPath"
		} else if tokenType == sessmodels.RefreshToken {
			pathType = "refreshTokenPath"
		}
		setCookie(config, res, cookieName, value, expires, pathType, request, userContext)
	} else if transferMethod == sessmodels.HeaderTransferMethod {
		headerName, err := getResponseHeaderNameForTokenType(tokenType)
		if err != nil {
			return err
		}

		setHeader(res, headerName, value, false)
		setHeader(res, "Access-Control-Expose-Headers", headerName, true)
	}
	return nil
}

func setHeader(res http.ResponseWriter, key, value string, allowDuplicateKey bool) {
	existingValue := res.Header().Get(strings.ToLower(key))
	if existingValue == "" {
		res.Header().Set(key, value)
	} else if allowDuplicateKey {
		/**
		We only want to append if it does not already exist

		For example if the caller is trying to add front token to the access control exposed headers property
		we do not want to append if something else had already added it
		*/
		if !strings.Contains(existingValue, value) {
			res.Header().Set(key, existingValue+", "+value)
		}
	} else {
		res.Header().Set(key, value)
	}
}

func setCookie(config sessmodels.TypeNormalisedInput, res http.ResponseWriter, name string, value string, expires uint64, pathType string, request *http.Request, userContext supertokens.UserContext) error {
	var domain string
	if config.CookieDomain != nil {
		domain = *config.CookieDomain
	}
	secure := config.CookieSecure
	sameSite, err := config.GetCookieSameSite(request, userContext)

	if err != nil {
		return err
	}

	path := ""
	if pathType == "refreshTokenPath" {
		path = config.RefreshTokenPath.GetAsStringDangerous()
	} else if pathType == "accessTokenPath" {
		path = "/"
	}

	var sameSiteField = http.SameSiteNoneMode
	if sameSite == "lax" {
		sameSiteField = http.SameSiteLaxMode
	} else if sameSite == "strict" {
		sameSiteField = http.SameSiteStrictMode
	}

	httpOnly := true

	cookie := &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
		Expires:  time.Unix(int64(expires/1000), 0),
		Path:     path,
		SameSite: sameSiteField,
	}
	setCookieValue(res, cookie)
	return nil
}

func GetAuthmodeFromHeader(req *http.Request) *sessmodels.TokenTransferMethod {
	val := getHeader(req, authModeHeaderKey)
	if val == nil {
		return nil
	}
	valLcase := sessmodels.TokenTransferMethod(strings.ToLower(*val))
	return &valLcase
}

func getHeader(request *http.Request, key string) *string {
	value := request.Header.Get(key)
	if value == "" {
		return nil
	}
	return &value
}

func GetCookieValue(request *http.Request, key string) *string {
	cookies := request.Cookies()
	if len(cookies) == 0 {
		return nil
	}

	for _, cookie := range cookies {
		if cookie.Name == key {
			val, err := url.QueryUnescape(cookie.Value)
			if err != nil {
				return nil
			}
			return &val
		}
	}
	return nil
}

// setCookieValue replaces cookie.go SetCookie, it replaces the cookie values instead of appending them
func setCookieValue(w http.ResponseWriter, cookie *http.Cookie) {
	cookieHeader := w.Header().Values("Set-Cookie")
	if len(cookieHeader) == 0 {
		w.Header().Set("Set-Cookie", cookie.String())
		return
	}
	existingCookies := make(map[string]string, len(cookieHeader))
	// map existing cookies by cookie name
	for _, ch := range cookieHeader {
		existingCookies[getCookieName(ch)] = ch
	}
	// replace if already existing
	existingCookies[getCookieName(cookie.String())] = cookie.String()
	// clear previous cookies from the headers
	w.Header().Del("Set-Cookie")
	// and add them back
	for _, ck := range existingCookies {
		w.Header().Add("Set-Cookie", ck)
	}
}

func getCookieName(cookie string) string {
	parts := strings.Split(textproto.TrimString(cookie), ";")
	if len(parts) == 1 && parts[0] == "" {
		return ""
	}
	parts[0] = textproto.TrimString(parts[0])
	kv := strings.Split(parts[0], "=")
	if len(kv) == 0 {
		return ""
	}
	return kv[0]
}

// ClearSessionCookiesFromOlderCookieDomain addresses an edge case where changing the cookieDomain config on the server can
// lead to session integrity issues. For instance, if the API server URL is 'api.example.com'
// with a cookie domain of '.example.com', and the server updates the cookie domain to 'api.example.com',
// the client may retain cookies with both '.example.com' and 'api.example.com' domains.
//
// Consequently, if the server chooses the older cookie, session invalidation occurs, potentially
// resulting in an infinite refresh loop. To fix this, users are asked to specify "OlderCookieDomain" in
// the config.
//
// This function checks for multiple cookies with the same name and clears the cookies for the older domain.
func ClearSessionCookiesFromOlderCookieDomain(req *http.Request, res http.ResponseWriter, config sessmodels.TypeNormalisedInput, userContext supertokens.UserContext) error {
	allowedTransferMethod := config.GetTokenTransferMethod(req, false, userContext)

	// If the transfer method is 'header', there's no need to clear cookies immediately, even if there are multiple in the request.
	if allowedTransferMethod == sessmodels.HeaderTransferMethod {
		return nil
	}

	didClearCookies := false

	tokenTypes := []sessmodels.TokenType{sessmodels.AccessToken, sessmodels.RefreshToken}
	for _, token := range tokenTypes {
		if hasMultipleCookiesForTokenType(req, token) {
			// If a request has multiple session cookies and 'olderCookieDomain' is
			// unset, we can't identify the correct cookie for refreshing the session.
			// Using the wrong cookie can cause an infinite refresh loop. To avoid this,
			// we throw a 500 error asking the user to set 'olderCookieDomain'.
			if config.OlderCookieDomain == nil {
				return errors.New(`The request contains multiple session cookies. This may happen if you've changed the 'cookieDomain' value in your configuration. To clear tokens from the previous domain, set 'olderCookieDomain' in your config.`)
			}

			supertokens.LogDebugMessage(fmt.Sprint("ClearSessionCookiesFromOlderCookieDomain: Clearing duplicate ", token, " cookie with domain ", config.OlderCookieDomain))
			config.CookieDomain = config.OlderCookieDomain
			setToken(config, res, token, "", 0, sessmodels.CookieTransferMethod, req, userContext)

			didClearCookies = true
		}
	}

	if didClearCookies {
		return sessionError.ClearDuplicateSessionCookiesError{
			Msg: "The request contains multiple session cookies. We are clearing the cookie from OlderCookieDomain. Session will be refreshed in the next refresh call.",
		}
	}

	return nil
}

func hasMultipleCookiesForTokenType(req *http.Request, tokenType sessmodels.TokenType) bool {
	// Count of cookies with the specified token type
	count := 0

	// Loop through each cookie in the request
	for _, cookie := range req.Cookies() {
		// Check if the cookie's name matches the token type
		cookieName, _ := getCookieNameFromTokenType(tokenType)
		if cookie.Name == cookieName {
			count++
		}
	}

	// If count is greater than 1, then there are multiple cookies with the given token type
	return count > 1
}
