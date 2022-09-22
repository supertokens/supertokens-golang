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
	"github.com/supertokens/supertokens-golang/recipe/session/cookiesandheaders"
	"github.com/supertokens/supertokens-golang/recipe/session/sessionwithjwt"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"golang.org/x/net/publicsuffix"
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

	topLevelAPIDomain, err := GetTopLevelDomainForSameSiteResolution(appInfo.APIDomain.GetAsStringDangerous())
	if err != nil {
		return sessmodels.TypeNormalisedInput{}, err
	}
	topLevelWebsiteDomain, err := GetTopLevelDomainForSameSiteResolution(appInfo.WebsiteDomain.GetAsStringDangerous())
	if err != nil {
		return sessmodels.TypeNormalisedInput{}, err
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
	if topLevelAPIDomain != topLevelWebsiteDomain {
		cookieSameSite = cookieSameSite_NONE
	}
	if apiDomainScheme != websiteDomainScheme {
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

	IsAnIPAPIDomain, err := supertokens.IsAnIPAddress(topLevelAPIDomain)
	if err != nil {
		return sessmodels.TypeNormalisedInput{}, err
	}
	IsAnIPWebsiteDomain, err := supertokens.IsAnIPAddress(topLevelWebsiteDomain)
	if err != nil {
		return sessmodels.TypeNormalisedInput{}, err
	}

	if cookieSameSite == cookieSameSite_NONE &&
		!cookieSecure && !((topLevelAPIDomain == "localhost" || IsAnIPAPIDomain) &&
		(topLevelWebsiteDomain == "localhost" || IsAnIPWebsiteDomain)) {
		// We can allow insecure cookie when both website & API domain are localhost or an IP
		// When either of them is a different domain, API domain needs to have https and a secure cookie to work
		return sessmodels.TypeNormalisedInput{}, errors.New("Since your API and website domain are different, for sessions to work, please use https on your apiDomain and dont set cookieSecure to false.")
	}

	refreshAPIPath, err := supertokens.NewNormalisedURLPath(refreshAPIPath)
	if err != nil {
		return sessmodels.TypeNormalisedInput{}, err
	}

	Jwt := sessmodels.JWTNormalisedConfig{Enable: false, PropertyNameInAccessTokenPayload: "jwt"}
	if config != nil && config.Jwt != nil {
		Jwt.Enable = config.Jwt.Enable
		Jwt.Issuer = config.Jwt.Issuer
		if config.Jwt.PropertyNameInAccessTokenPayload != nil {
			Jwt.PropertyNameInAccessTokenPayload = *config.Jwt.PropertyNameInAccessTokenPayload
		}
	}
	if sessionwithjwt.ACCESS_TOKEN_PAYLOAD_JWT_PROPERTY_NAME_KEY == Jwt.PropertyNameInAccessTokenPayload {
		return sessmodels.TypeNormalisedInput{}, errors.New(sessionwithjwt.ACCESS_TOKEN_PAYLOAD_JWT_PROPERTY_NAME_KEY + " is a reserved property name, please use a different key name for the jwt")
	}

	typeNormalisedInput := sessmodels.TypeNormalisedInput{
		RefreshTokenPath:         appInfo.APIBasePath.AppendPath(refreshAPIPath),
		CookieDomain:             cookieDomain,
		CookieSameSite:           cookieSameSite,
		CookieSecure:             cookieSecure,
		SessionExpiredStatusCode: sessionExpiredStatusCode,
		InvalidClaimStatusCode:   invalidClaimStatusCode,
		AntiCsrf:                 antiCsrf,
		ErrorHandlers:            errorHandlers,
		Jwt:                      Jwt,
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

func GetTopLevelDomainForSameSiteResolution(URL string) (string, error) {
	urlObj, err := url.Parse(URL)
	if err != nil {
		return "", err
	}
	hostname := urlObj.Hostname()
	isAnIP, err := supertokens.IsAnIPAddress(hostname)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(hostname, "localhost") || strings.HasPrefix(hostname, "localhost.org") || isAnIP {
		return "localhost", nil
	}
	parsedURL, err := publicsuffix.EffectiveTLDPlusOne(hostname)
	if err != nil {
		return "", errors.New("Please make sure that the apiDomain and websiteDomain have correct values")
	}
	return parsedURL, nil
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

func attachCreateOrRefreshSessionResponseToRes(config sessmodels.TypeNormalisedInput, res http.ResponseWriter, response sessmodels.CreateOrRefreshAPIResponse) {
	accessToken := response.AccessToken
	refreshToken := response.RefreshToken
	idRefreshToken := response.IDRefreshToken
	cookiesandheaders.SetFrontTokenInHeaders(res, response.Session.UserID, response.AccessToken.Expiry, response.Session.UserDataInAccessToken)
	cookiesandheaders.AttachAccessTokenToCookie(config, res, accessToken.Token, accessToken.Expiry)
	cookiesandheaders.AttachRefreshTokenToCookie(config, res, refreshToken.Token, refreshToken.Expiry)
	cookiesandheaders.SetIDRefreshTokenInHeaderAndCookie(config, res, idRefreshToken.Token, idRefreshToken.Expiry)
	if response.AntiCsrfToken != nil {
		cookiesandheaders.SetAntiCsrfTokenInHeaders(res, *response.AntiCsrfToken)
	}
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

func frontendHasInterceptor(req *http.Request) bool {
	return cookiesandheaders.GetRidFromHeader(req) != nil
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
