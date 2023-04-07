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
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func validateAndNormaliseUserInput(appInfo supertokens.NormalisedAppinfo, config *sessmodels.TypeInput) (sessmodels.TypeNormalisedInput, error) {
	var (
		cookieDomain *string = nil
		err          error
	)

	if config != nil && config.CookieDomain != nil {
		cookieDomain, err = normaliseSessionScopeOrThrowError(*config.CookieDomain)
		if err != nil {
			return sessmodels.TypeNormalisedInput{}, err
		}
	}

	apiDomainScheme, err := GetURLScheme(appInfo.APIDomain.GetAsStringDangerous())
	if err != nil {
		return sessmodels.TypeNormalisedInput{}, err
	}
	websiteDomainScheme, err := GetURLScheme(appInfo.WebsiteDomain.GetAsStringDangerous())
	if err != nil {
		return sessmodels.TypeNormalisedInput{}, err
	}

	cookieSameSite := cookieSameSite_LAX
	if apiDomainScheme != websiteDomainScheme || appInfo.TopLevelAPIDomain != appInfo.TopLevelWebsiteDomain {
		cookieSameSite = cookieSameSite_NONE
	}

	if config != nil && config.CookieSameSite != nil {
		cookieSameSite, err = normaliseSameSiteOrThrowError(*config.CookieSameSite)
		if err != nil {
			return sessmodels.TypeNormalisedInput{}, err
		}
	}

	cookieSecure := false
	if config == nil || config.CookieSecure == nil {
		cookieSecure = strings.HasPrefix(appInfo.APIDomain.GetAsStringDangerous(), "https")
	} else {
		cookieSecure = *config.CookieSecure
	}

	sessionExpiredStatusCode := 401
	if config != nil && config.SessionExpiredStatusCode != nil {
		sessionExpiredStatusCode = *config.SessionExpiredStatusCode
	}

	invalidClaimStatusCode := 403
	if config != nil && config.InvalidClaimStatusCode != nil {
		invalidClaimStatusCode = *config.InvalidClaimStatusCode
	}

	if sessionExpiredStatusCode == invalidClaimStatusCode {
		return sessmodels.TypeNormalisedInput{}, errors.New("SessionExpiredStatusCode and InvalidClaimStatusCode cannot have the same value")
	}

	if config != nil && config.AntiCsrf != nil {
		if *config.AntiCsrf != antiCSRF_NONE && *config.AntiCsrf != antiCSRF_VIA_CUSTOM_HEADER && *config.AntiCsrf != antiCSRF_VIA_TOKEN {
			return sessmodels.TypeNormalisedInput{}, errors.New("antiCsrf config must be one of 'NONE' or 'VIA_CUSTOM_HEADER' or 'VIA_TOKEN'")
		}
	}

	antiCsrf := antiCSRF_NONE
	if config == nil || config.AntiCsrf == nil {
		if cookieSameSite == cookieSameSite_NONE {
			antiCsrf = antiCSRF_VIA_CUSTOM_HEADER
		} else {
			antiCsrf = antiCSRF_NONE
		}
	} else {
		antiCsrf = *config.AntiCsrf
	}

	useDynamicAccessTokenSigningKey := true

	if config != nil && config.UseDynamicAccessTokenSigningKey != nil {
		useDynamicAccessTokenSigningKey = *config.UseDynamicAccessTokenSigningKey
	}

	errorHandlers := sessmodels.NormalisedErrorHandlers{
		OnTokenTheftDetected: func(sessionHandle string, userID string, req *http.Request, res http.ResponseWriter) error {
			recipeInstance, err := getRecipeInstanceOrThrowError()
			if err != nil {
				return err
			}
			return sendTokenTheftDetectedResponse(*recipeInstance, sessionHandle, userID, req, res)
		},
		OnTryRefreshToken: func(message string, req *http.Request, res http.ResponseWriter) error {
			recipeInstance, err := getRecipeInstanceOrThrowError()
			if err != nil {
				return err
			}
			return sendTryRefreshTokenResponse(*recipeInstance, message, req, res)
		},
		OnUnauthorised: func(message string, req *http.Request, res http.ResponseWriter) error {
			recipeInstance, err := getRecipeInstanceOrThrowError()
			if err != nil {
				return err
			}
			return sendUnauthorisedResponse(*recipeInstance, message, req, res)
		},
		OnInvalidClaim: func(validationErrors []claims.ClaimValidationError, req *http.Request, res http.ResponseWriter) error {
			recipeInstance, err := getRecipeInstanceOrThrowError()
			if err != nil {
				return err
			}
			return sendInvalidClaimResponse(*recipeInstance, validationErrors, req, res)
		},
	}

	if config != nil && config.ErrorHandlers != nil {
		if config.ErrorHandlers.OnTokenTheftDetected != nil {
			errorHandlers.OnTokenTheftDetected = config.ErrorHandlers.OnTokenTheftDetected
		}
		if config.ErrorHandlers.OnUnauthorised != nil {
			errorHandlers.OnUnauthorised = config.ErrorHandlers.OnUnauthorised
		}
		if config.ErrorHandlers.OnInvalidClaim != nil {
			errorHandlers.OnInvalidClaim = config.ErrorHandlers.OnInvalidClaim
		}
	}

	refreshAPIPath, err := supertokens.NewNormalisedURLPath(refreshAPIPath)
	if err != nil {
		return sessmodels.TypeNormalisedInput{}, err
	}

	if config == nil {
		config = &sessmodels.TypeInput{}
	}

	if config.GetTokenTransferMethod == nil {
		config.GetTokenTransferMethod = defaultGetTokenTransferMethod
	}

	typeNormalisedInput := sessmodels.TypeNormalisedInput{
		RefreshTokenPath:                appInfo.APIBasePath.AppendPath(refreshAPIPath),
		CookieDomain:                    cookieDomain,
		CookieSameSite:                  cookieSameSite,
		CookieSecure:                    cookieSecure,
		SessionExpiredStatusCode:        sessionExpiredStatusCode,
		InvalidClaimStatusCode:          invalidClaimStatusCode,
		AntiCsrf:                        antiCsrf,
		UseDynamicAccessTokenSigningKey: useDynamicAccessTokenSigningKey,
		ErrorHandlers:                   errorHandlers,
		GetTokenTransferMethod:          config.GetTokenTransferMethod,
		Override: sessmodels.OverrideStruct{
			Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
				return originalImplementation
			}, APIs: func(originalImplementation sessmodels.APIInterface) sessmodels.APIInterface {
				return originalImplementation
			},
			OpenIdFeature: nil},
	}

	if config != nil && config.Override != nil {
		if config.Override.Functions != nil {
			typeNormalisedInput.Override.Functions = config.Override.Functions
		}
		if config.Override.APIs != nil {
			typeNormalisedInput.Override.APIs = config.Override.APIs
		}
		typeNormalisedInput.Override.OpenIdFeature = config.Override.OpenIdFeature
	}

	return typeNormalisedInput, nil
}
func normaliseSameSiteOrThrowError(sameSite string) (string, error) {
	sameSite = strings.TrimSpace(sameSite)
	sameSite = strings.ToLower(sameSite)
	if sameSite != cookieSameSite_STRICT && sameSite != cookieSameSite_LAX && sameSite != cookieSameSite_NONE {
		return "", errors.New(`cookie same site must be one of "strict", "lax", or "none"`)
	}
	return sameSite, nil
}

func GetURLScheme(URL string) (string, error) {
	urlObj, err := url.Parse(URL)
	if err != nil {
		return "", err
	}
	return urlObj.Scheme, nil
}

func normaliseSessionScopeOrThrowError(sessionScope string) (*string, error) {
	sessionScope = strings.TrimSpace(sessionScope)
	sessionScope = strings.ToLower(sessionScope)

	sessionScope = strings.TrimPrefix(sessionScope, ".")

	if !strings.HasPrefix(sessionScope, "http://") && !strings.HasPrefix(sessionScope, "https://") {
		sessionScope = "http://" + sessionScope
	}

	urlObj, err := url.Parse(sessionScope)
	if err != nil {
		return nil, errors.New("Please provide a valid sessionScope")
	}

	sessionScope = urlObj.Hostname()
	sessionScope = strings.TrimPrefix(sessionScope, ".")

	noDotNormalised := sessionScope

	isAnIP, err := supertokens.IsAnIPAddress(sessionScope)
	if err != nil {
		return nil, err
	}
	if sessionScope == "localhost" || isAnIP {
		noDotNormalised = sessionScope
	}
	if strings.HasPrefix(sessionScope, ".") {
		noDotNormalised = "." + sessionScope
	}
	return &noDotNormalised, nil
}

func getCurrTimeInMS() uint64 {
	return uint64(time.Now().UnixNano() / 1000000)
}

func attachCreateOrRefreshSessionResponseToRes(config sessmodels.TypeNormalisedInput, res http.ResponseWriter, response sessmodels.CreateOrRefreshAPIResponse, tokenTransferMethod sessmodels.TokenTransferMethod) {
	accessToken := response.AccessToken
	refreshToken := response.RefreshToken
	setFrontTokenInHeaders(res, response.Session.UserID, response.AccessToken.Expiry, response.Session.UserDataInAccessToken)
	setToken(
		config,
		res,
		sessmodels.AccessToken,
		accessToken.Token,
		// We set the expiration to 100 years, because we can't really access the expiration of the refresh token everywhere we are setting it.
		// This should be safe to do, since this is only the validity of the cookie (set here or on the frontend) but we check the expiration of the JWT anyway.
		// Even if the token is expired the presence of the token indicates that the user could have a valid refresh
		// Setting them to infinity would require special case handling on the frontend and just adding 10 years seems enough.
		getCurrTimeInMS()+3153600000000,
		tokenTransferMethod,
	)
	setToken(config, res, sessmodels.RefreshToken, refreshToken.Token, refreshToken.Expiry, tokenTransferMethod)

	if response.AntiCsrfToken != nil {
		setAntiCsrfTokenInHeaders(res, *response.AntiCsrfToken)
	}
}

func SetAccessTokenInResponse(config sessmodels.TypeNormalisedInput, res http.ResponseWriter, accessToken sessmodels.CreateOrRefreshAPIResponseToken, session sessmodels.SessionStruct, tokenTransferMethod sessmodels.TokenTransferMethod) error {
	responseToken, err := parseJWTWithoutSignatureVerification(accessToken.Token)
	if err != nil {
		return err
	}

	var payload map[string]interface{}

	if responseToken.Version == 3 {
		payload = responseToken.Payload
	} else {
		payload = session.UserDataInAccessToken
	}

	setFrontTokenInHeaders(res, session.UserID, accessToken.Expiry, payload)
	setToken(config, res, sessmodels.AccessToken, accessToken.Token, getCurrTimeInMS()+3153600000000, tokenTransferMethod)
	return nil
}

func sendTryRefreshTokenResponse(recipeInstance Recipe, _ string, _ *http.Request, response http.ResponseWriter) error {
	return supertokens.SendNon200ResponseWithMessage(response, "try refresh token", recipeInstance.Config.SessionExpiredStatusCode)
}

func sendUnauthorisedResponse(recipeInstance Recipe, _ string, _ *http.Request, response http.ResponseWriter) error {
	return supertokens.SendNon200ResponseWithMessage(response, "unauthorised", recipeInstance.Config.SessionExpiredStatusCode)
}

func sendInvalidClaimResponse(recipeInstance Recipe, claimValidationErrors []claims.ClaimValidationError, _ *http.Request, response http.ResponseWriter) error {
	return supertokens.SendNon200Response(response, recipeInstance.Config.InvalidClaimStatusCode, map[string]interface{}{
		"message":               "invalid claim",
		"claimValidationErrors": claimValidationErrors,
	})
}

func sendTokenTheftDetectedResponse(recipeInstance Recipe, sessionHandle string, _ string, _ *http.Request, response http.ResponseWriter) error {
	_, err := (*recipeInstance.RecipeImpl.RevokeSession)(sessionHandle, &map[string]interface{}{})
	if err != nil {
		return err
	}
	return supertokens.SendNon200ResponseWithMessage(response, "token theft detected", recipeInstance.Config.SessionExpiredStatusCode)
}

func getKeyInfoFromJson(response map[string]interface{}) []sessmodels.KeyInfo {
	keyList := []sessmodels.KeyInfo{}

	_, ok := response["jwtSigningPublicKeyList"]
	if ok {
		for _, k := range response["jwtSigningPublicKeyList"].([]interface{}) {
			keyList = append(keyList, sessmodels.KeyInfo{
				PublicKey:  (k.((map[string]interface{})))["publicKey"].(string),
				ExpiryTime: uint64((k.((map[string]interface{})))["expiryTime"].(float64)),
				CreatedAt:  uint64((k.((map[string]interface{})))["createdAt"].(float64)),
			})
		}
	}

	return keyList
}

func validateClaimsInPayload(claimValidators []claims.SessionClaimValidator, newAccessTokenPayload map[string]interface{}, userContext supertokens.UserContext) []claims.ClaimValidationError {
	validationErrors := []claims.ClaimValidationError{}

	for _, validator := range claimValidators {
		claimValidationResult := validator.Validate(newAccessTokenPayload, userContext)
		supertokens.LogDebugMessage(fmt.Sprint("validateClaimsInPayload ", validator.ID, " validation res ", claimValidationResult))
		if !claimValidationResult.IsValid {
			validationErrors = append(validationErrors, claims.ClaimValidationError{
				ID:     validator.ID,
				Reason: claimValidationResult.Reason,
			})
		}
	}
	return validationErrors
}

func getRequiredClaimValidators(
	sessionContainer sessmodels.SessionContainer,
	overrideGlobalClaimValidators func(globalClaimValidators []claims.SessionClaimValidator, sessionContainer sessmodels.SessionContainer, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error),
	userContext supertokens.UserContext,
) ([]claims.SessionClaimValidator, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	claimValidatorsAddedByOtherRecipes := instance.getClaimValidatorsAddedByOtherRecipes()
	globalClaimValidators, err := (*instance.RecipeImpl.GetGlobalClaimValidators)(sessionContainer.GetUserID(), claimValidatorsAddedByOtherRecipes, userContext)
	if err != nil {
		return nil, err
	}
	if overrideGlobalClaimValidators != nil {
		globalClaimValidators, err = overrideGlobalClaimValidators(globalClaimValidators, sessionContainer, userContext)
		if err != nil {
			return nil, err
		}
	}
	return globalClaimValidators, nil
}

func defaultGetTokenTransferMethod(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
	// We allow fallback (checking headers then cookies) by default when validating

	if !forCreateNewSession {
		return sessmodels.AnyTransferMethod
	}

	// In create new session we respect the frontend preference by default
	authMode := getAuthmodeFromHeader(req)
	if authMode == nil {
		return sessmodels.AnyTransferMethod
	}
	switch *authMode {
	case sessmodels.CookieTransferMethod:
		return sessmodels.CookieTransferMethod
	case sessmodels.HeaderTransferMethod:
		return sessmodels.HeaderTransferMethod
	default:
		return sessmodels.AnyTransferMethod
	}
}

func getRidFromHeader(req *http.Request) *string {
	rid := req.Header.Get("rid")
	if rid == "" {
		return nil
	}
	return &rid
}
