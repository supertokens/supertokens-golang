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
	defaultErrors "errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/supertokens/supertokens-golang/recipe/session/claims"

	"github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

// We are defining this here to reduce the scope of legacy code
const legacyIdRefreshTokenCookieName = "sIdRefreshToken"

func CreateNewSessionInRequest(req *http.Request, res http.ResponseWriter, tenantId string, config sessmodels.TypeNormalisedInput, appInfo supertokens.NormalisedAppinfo, recipeInstance Recipe, recipeImpl sessmodels.RecipeInterface, userID string, accessTokenPayload map[string]interface{}, sessionDataInDatabase map[string]interface{}, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
	supertokens.LogDebugMessage("createNewSession: Started")

	claimsAddedByOtherRecipes := recipeInstance.GetClaimsAddedByOtherRecipes()
	finalAccessTokenPayload := accessTokenPayload
	if finalAccessTokenPayload == nil {
		finalAccessTokenPayload = map[string]interface{}{}
	}

	issuer := appInfo.APIDomain.GetAsStringDangerous() + appInfo.APIBasePath.GetAsStringDangerous()
	finalAccessTokenPayload["iss"] = issuer

	for _, protectedProp := range protectedProps {
		delete(finalAccessTokenPayload, protectedProp)
	}

	for _, claim := range claimsAddedByOtherRecipes {
		_finalAccessTokenPayload, err := claim.Build(userID, tenantId, finalAccessTokenPayload, userContext)
		if err != nil {
			return nil, err
		}

		finalAccessTokenPayload = _finalAccessTokenPayload
	}

	supertokens.LogDebugMessage("createNewSession: Access token payload built")

	outputTokenTransferMethod := config.GetTokenTransferMethod(req, true, userContext)
	if outputTokenTransferMethod == sessmodels.AnyTransferMethod {
		authMode := GetAuthmodeFromHeader(req)
		if authMode != nil && *authMode == sessmodels.CookieTransferMethod {
			outputTokenTransferMethod = *authMode
		} else {
			outputTokenTransferMethod = sessmodels.HeaderTransferMethod
		}
	}

	supertokens.LogDebugMessage(fmt.Sprintf("createNewSession: using transfer method %s", outputTokenTransferMethod))

	isTopLevelAPIDomainIPAddress, err := supertokens.IsAnIPAddress(appInfo.TopLevelAPIDomain)
	if err != nil {
		return nil, err
	}

	topLevelWebsiteDomain, err := appInfo.GetTopLevelWebsiteDomain(req, userContext)
	if err != nil {
		return nil, err
	}

	isTopLevelWebsiteDomainIPAddress, err := supertokens.IsAnIPAddress(topLevelWebsiteDomain)
	if err != nil {
		return nil, err
	}

	cookieSameSite, err := config.GetCookieSameSite(req, userContext)
	if err != nil {
		return nil, err
	}

	if outputTokenTransferMethod == sessmodels.CookieTransferMethod &&
		cookieSameSite == "none" &&
		!config.CookieSecure &&
		!((appInfo.TopLevelAPIDomain == "localhost" || isTopLevelAPIDomainIPAddress) &&
			(topLevelWebsiteDomain == "localhost" || isTopLevelWebsiteDomainIPAddress)) {
		// We can allow insecure cookie when both website & API domain are localhost or an IP
		// When either of them is a different domain, API domain needs to have https and a secure cookie to work
		return nil, defaultErrors.New("Since your API and website domain are different, for sessions to work, please use https on your apiDomain and dont set cookieSecure to false.")
	}

	disableAntiCSRF := outputTokenTransferMethod == sessmodels.HeaderTransferMethod

	sessionResponse, err := (*recipeImpl.CreateNewSession)(userID, finalAccessTokenPayload, sessionDataInDatabase, &disableAntiCSRF, tenantId, userContext)

	if err != nil {
		return nil, err
	}

	supertokens.LogDebugMessage("createNewSession: Session created in core built")

	for _, tokenTransferMethod := range AvailableTokenTransferMethods {
		if tokenTransferMethod != outputTokenTransferMethod {
			token, err := GetToken(req, sessmodels.AccessToken, tokenTransferMethod)
			if err != nil {
				return nil, err
			}
			if token != nil {
				ClearSession(config, res, tokenTransferMethod, req, userContext)
			}
		}
	}

	supertokens.LogDebugMessage("createNewSession: Cleared old tokens")

	sessionResponse.AttachToRequestResponseWithContext(sessmodels.RequestResponseInfo{
		Res:                 res,
		Req:                 req,
		TokenTransferMethod: outputTokenTransferMethod,
	}, userContext)
	supertokens.LogDebugMessage("createNewSession: Attached new tokens to res")

	return sessionResponse, nil
}

func GetSessionFromRequest(req *http.Request, res http.ResponseWriter, config sessmodels.TypeNormalisedInput, options *sessmodels.VerifySessionOptions, recipeImpl sessmodels.RecipeInterface, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
	idRefreshToken := GetCookieValue(req, legacyIdRefreshTokenCookieName)
	if idRefreshToken != nil {
		supertokens.LogDebugMessage("GetSessionFromRequest: Returning TryRefreshTokenError because the request is using a legacy session and should be refreshed")
		return nil, errors.TryRefreshTokenError{
			Msg: "using legacy session, please call the refresh API",
		}
	}

	sessionOptional := options != nil && options.SessionRequired != nil && !*options.SessionRequired
	supertokens.LogDebugMessage(fmt.Sprintf("getSession: optional validation %v", sessionOptional))

	accessTokens := map[sessmodels.TokenTransferMethod]*sessmodels.ParsedJWTInfo{}

	// We check all token transfer methods for available access tokens
	for _, tokenTransferMethod := range AvailableTokenTransferMethods {
		token, err := GetToken(req, sessmodels.AccessToken, tokenTransferMethod)
		if err != nil {
			return nil, err
		}
		if token != nil {
			parsedToken, err := ParseJWTWithoutSignatureVerification(*token)
			if err != nil {
				supertokens.LogDebugMessage(fmt.Sprintf("getSession: ignoring token in %s, because token parsing failed", tokenTransferMethod))
			} else {
				err := ValidateAccessTokenStructure(parsedToken.Payload, parsedToken.Version)
				if err != nil {
					supertokens.LogDebugMessage(fmt.Sprintf("getSession: ignoring token in %s, because it doesn't match our access token structure", tokenTransferMethod))
				} else {
					supertokens.LogDebugMessage(fmt.Sprintf("getSession: got access token from %s", tokenTransferMethod))
					accessTokens[tokenTransferMethod] = &parsedToken
				}
			}
		}
	}

	allowedTokenTransferMethod := config.GetTokenTransferMethod(req, false, userContext)

	var requestTokenTransferMethod *sessmodels.TokenTransferMethod
	var accessToken *sessmodels.ParsedJWTInfo

	if (allowedTokenTransferMethod == sessmodels.AnyTransferMethod || allowedTokenTransferMethod == sessmodels.HeaderTransferMethod) && (accessTokens[sessmodels.HeaderTransferMethod] != nil) {
		supertokens.LogDebugMessage("getSession: using header transfer method")
		headerMethod := sessmodels.HeaderTransferMethod
		requestTokenTransferMethod = &headerMethod
		accessToken = accessTokens[sessmodels.HeaderTransferMethod]
	} else if (allowedTokenTransferMethod == sessmodels.AnyTransferMethod || allowedTokenTransferMethod == sessmodels.CookieTransferMethod) && (accessTokens[sessmodels.CookieTransferMethod] != nil) {
		supertokens.LogDebugMessage("getSession: using cookie transfer method")
		cookieMethod := sessmodels.CookieTransferMethod
		requestTokenTransferMethod = &cookieMethod
		accessToken = accessTokens[sessmodels.CookieTransferMethod]
	}

	antiCsrfToken := GetAntiCsrfTokenFromHeaders(req)
	var doAntiCsrfCheck *bool

	if options != nil {
		doAntiCsrfCheck = options.AntiCsrfCheck
	}

	if doAntiCsrfCheck == nil {
		doAntiCsrfCheckBool := req.Method != http.MethodGet
		doAntiCsrfCheck = &doAntiCsrfCheckBool
	}

	False := false
	if requestTokenTransferMethod != nil && *requestTokenTransferMethod == sessmodels.HeaderTransferMethod {
		doAntiCsrfCheck = &False
	}

	if accessToken == nil {
		doAntiCsrfCheck = &False
	}

	antiCsrf := config.AntiCsrfFunctionOrString.StrValue
	if antiCsrf == "" {
		antiCsrfTemp, err := config.AntiCsrfFunctionOrString.FunctionValue(req, userContext)
		if err != nil {
			return nil, err
		}
		antiCsrf = antiCsrfTemp
	}

	if *doAntiCsrfCheck && antiCsrf == AntiCSRF_VIA_CUSTOM_HEADER {
		if antiCsrf == AntiCSRF_VIA_CUSTOM_HEADER {
			if GetRidFromHeader(req) == nil {
				supertokens.LogDebugMessage("getSession: Returning TRY_REFRESH_TOKEN because custom header (rid) was not passed")
				return nil, errors.TryRefreshTokenError{
					Msg: "anti-csrf check failed. Please pass 'rid: \"session\"' header in the request, or set doAntiCsrfCheck to false for this API",
				}
			}

			supertokens.LogDebugMessage("getSession: VIA_CUSTOM_HEADER anti-csrf check passed")
			False := false
			doAntiCsrfCheck = &False
		}
	}

	supertokens.LogDebugMessage("getSession: Value of doAntiCsrfCheck is: " + strconv.FormatBool(*doAntiCsrfCheck))

	_verifySessionOptionsToPass := sessmodels.VerifySessionOptions{
		AntiCsrfCheck: doAntiCsrfCheck,
	}

	if options != nil {
		_verifySessionOptionsToPass = sessmodels.VerifySessionOptions{
			AntiCsrfCheck:                 _verifySessionOptionsToPass.AntiCsrfCheck,
			SessionRequired:               options.SessionRequired,
			CheckDatabase:                 options.CheckDatabase,
			OverrideGlobalClaimValidators: options.OverrideGlobalClaimValidators,
		}
	}

	var rawTokenString *string

	if accessToken != nil {
		rawTokenString = &accessToken.RawTokenString
	}

	result, err := (*recipeImpl.GetSession)(rawTokenString, antiCsrfToken, &_verifySessionOptionsToPass, userContext)

	if err != nil {
		return nil, err
	}

	sessionResult := result
	if sessionResult != nil {
		var overrideGlobalClaimValidators func(globalClaimValidators []claims.SessionClaimValidator, sessionContainer sessmodels.SessionContainer, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error) = nil
		if options != nil {
			overrideGlobalClaimValidators = options.OverrideGlobalClaimValidators
		}

		if err != nil {
			return nil, err
		}
		claimValidators, err := GetRequiredClaimValidators(sessionResult, overrideGlobalClaimValidators, userContext)

		if err != nil {
			return nil, err
		}

		err = (*sessionResult).AssertClaimsWithContext(claimValidators, userContext)
		if err != nil {
			return nil, err
		}

		// requestTransferMethod can only be nil here if the user has overridden GetSession
		// to load the session by a custom method in that (very niche) case they also need to
		// override how the session is attached to the response.
		// In that scenario the transferMethod passed to attachToRequestResponse likely doesn't
		// matter, still, we follow the general fallback logic
		transferMethod := sessmodels.HeaderTransferMethod

		if requestTokenTransferMethod != nil {
			transferMethod = *requestTokenTransferMethod
		} else if allowedTokenTransferMethod != sessmodels.AnyTransferMethod {
			transferMethod = allowedTokenTransferMethod
		}

		err = (*sessionResult).AttachToRequestResponseWithContext(sessmodels.RequestResponseInfo{
			Res:                 res,
			Req:                 req,
			TokenTransferMethod: transferMethod,
		}, userContext)

		if err != nil {
			return nil, err
		}
	}

	return sessionResult, nil
}

func RefreshSessionInRequest(req *http.Request, res http.ResponseWriter, config sessmodels.TypeNormalisedInput, recipeImpl sessmodels.RecipeInterface, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
	supertokens.LogDebugMessage("refreshSession: Started")

	refreshTokens := map[sessmodels.TokenTransferMethod]*string{}
	// We check all token transfer methods for available refresh tokens
	// We do this so that we can later clear all we are not overwriting
	for _, tokenTransferMethod := range AvailableTokenTransferMethods {
		token, err := GetToken(req, sessmodels.RefreshToken, tokenTransferMethod)
		if err != nil {
			return nil, err
		}
		refreshTokens[tokenTransferMethod] = token
		if token != nil {
			supertokens.LogDebugMessage("refreshSession: got refresh token from " + string(tokenTransferMethod))
		}
	}

	allowedTokenTransferMethod := config.GetTokenTransferMethod(req, false, userContext)
	supertokens.LogDebugMessage("refreshSession: getTokenTransferMethod returned " + string(allowedTokenTransferMethod))

	var requestTokenTransferMethod sessmodels.TokenTransferMethod
	var refreshToken *string

	if (allowedTokenTransferMethod == sessmodels.AnyTransferMethod || allowedTokenTransferMethod == sessmodels.HeaderTransferMethod) && refreshTokens[sessmodels.HeaderTransferMethod] != nil {
		supertokens.LogDebugMessage("refreshSession: using header transfer method")
		requestTokenTransferMethod = sessmodels.HeaderTransferMethod
		refreshToken = refreshTokens[sessmodels.HeaderTransferMethod]
	} else if (allowedTokenTransferMethod == sessmodels.AnyTransferMethod || allowedTokenTransferMethod == sessmodels.CookieTransferMethod) && refreshTokens[sessmodels.CookieTransferMethod] != nil {
		supertokens.LogDebugMessage("refreshSession: using cookie transfer method")
		requestTokenTransferMethod = sessmodels.CookieTransferMethod
		refreshToken = refreshTokens[sessmodels.CookieTransferMethod]
	} else {
		if GetCookieValue(req, legacyIdRefreshTokenCookieName) != nil {
			supertokens.LogDebugMessage("refreshSession: cleared legacy id refresh token because refresh token was not found")
			setCookie(config, res, legacyIdRefreshTokenCookieName, "", 0, "accessTokenPath", req, userContext)
		}

		supertokens.LogDebugMessage("refreshSession: UNAUTHORISED because refresh token in request is undefined")
		False := false
		return nil, errors.UnauthorizedError{
			Msg:         "Refresh token not found. Are you sending the refresh token in the request as a cookie?",
			ClearTokens: &False,
		}
	}

	antiCsrfToken := GetAntiCsrfTokenFromHeaders(req)
	disableAntiCSRF := requestTokenTransferMethod == sessmodels.HeaderTransferMethod
	antiCsrf := config.AntiCsrfFunctionOrString.StrValue
	if antiCsrf == "" {
		antiCsrfTemp, err := config.AntiCsrfFunctionOrString.FunctionValue(req, userContext)
		if err != nil {
			return nil, err
		}
		antiCsrf = antiCsrfTemp
	}

	if antiCsrf == AntiCSRF_VIA_CUSTOM_HEADER && !disableAntiCSRF {
		ridFromHeader := GetRidFromHeader(req)

		if ridFromHeader == nil {
			supertokens.LogDebugMessage("refreshSession: Returning UNAUTHORISED because custom header (rid) was not passed")
			clearTokens := false
			return nil, errors.UnauthorizedError{
				Msg:         "anti-csrf check failed. Please pass 'rid: \"session\"' header in the request.",
				ClearTokens: &clearTokens,
			}
		}

		disableAntiCSRF = true
	}

	result, err := (*recipeImpl.RefreshSession)(*refreshToken, antiCsrfToken, disableAntiCSRF, userContext)

	if err != nil {
		unauthorisedErr := errors.UnauthorizedError{}
		isUnauthorisedErr := defaultErrors.As(err, &unauthorisedErr)
		isTokenTheftDetectedErr := defaultErrors.As(err, &errors.TokenTheftDetectedError{})

		// This token isn't handled by getToken/setToken to limit the scope of this legacy/migration code
		if (isTokenTheftDetectedErr) || (isUnauthorisedErr && unauthorisedErr.ClearTokens != nil && *unauthorisedErr.ClearTokens) {
			if GetCookieValue(req, legacyIdRefreshTokenCookieName) != nil {
				supertokens.LogDebugMessage("refreshSession: cleared legacy id refresh token because refresh is clearing other tokens")
				setCookie(config, res, legacyIdRefreshTokenCookieName, "", 0, "accessTokenPath", req, userContext)
			}
		}

		if isUnauthorisedErr {
			supertokens.LogDebugMessage("RefreshSessionInRequest: Returning UnauthorizedError because RefreshSession returned an error")
		}

		return nil, err
	}

	supertokens.LogDebugMessage("refreshSession: Attaching refreshed session info as " + string(requestTokenTransferMethod))

	for _, tokenTransferMethod := range AvailableTokenTransferMethods {
		if tokenTransferMethod != requestTokenTransferMethod && refreshTokens[tokenTransferMethod] != nil {
			ClearSession(config, res, tokenTransferMethod, req, userContext)
		}
	}

	(*result).AttachToRequestResponseWithContext(sessmodels.RequestResponseInfo{
		Res:                 res,
		Req:                 req,
		TokenTransferMethod: requestTokenTransferMethod,
	}, userContext)

	supertokens.LogDebugMessage("refreshSession: Success!")

	if GetCookieValue(req, legacyIdRefreshTokenCookieName) != nil {
		supertokens.LogDebugMessage("refreshSession: cleared legacy id refresh token after successful refresh")
		setCookie(config, res, legacyIdRefreshTokenCookieName, "", 0, "accessTokenPath", req, userContext)
	}

	return result, nil
}
