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
	"encoding/json"
	defaultErrors "errors"
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func createNewSessionHelper(config sessmodels.TypeNormalisedInput, querier supertokens.Querier, userID string, disableAntiCsrf bool, AccessTokenPayload, sessionDataInDatabase map[string]interface{}) (sessmodels.CreateOrRefreshAPIResponse, error) {
	if AccessTokenPayload == nil {
		AccessTokenPayload = map[string]interface{}{}
	}
	if sessionDataInDatabase == nil {
		sessionDataInDatabase = map[string]interface{}{}
	}
	requestBody := map[string]interface{}{
		"userId":               userID,
		"userDataInJWT":        AccessTokenPayload,
		"userDataInDatabase":   sessionDataInDatabase,
		"enableAntiCsrf":       !disableAntiCsrf && config.AntiCsrf == AntiCSRF_VIA_TOKEN,
		"useDynamicSigningKey": config.UseDynamicAccessTokenSigningKey,
	}

	response, err := querier.SendPostRequest("/recipe/session", requestBody)
	if err != nil {
		return sessmodels.CreateOrRefreshAPIResponse{}, err
	}

	responseByte, err := json.Marshal(response)
	if err != nil {
		return sessmodels.CreateOrRefreshAPIResponse{}, err
	}
	var resp sessmodels.CreateOrRefreshAPIResponse
	err = json.Unmarshal(responseByte, &resp)
	if err != nil {
		return sessmodels.CreateOrRefreshAPIResponse{}, err
	}

	return resp, nil
}

func getSessionHelper(config sessmodels.TypeNormalisedInput, querier supertokens.Querier, parsedAccessToken sessmodels.ParsedJWTInfo, antiCsrfToken *string, doAntiCsrfCheck, alwaysCheckCore bool) (sessmodels.GetSessionResponse, error) {
	var accessTokenInfo *AccessTokenInfoStruct = nil
	var err error = nil
	combinedJwks, jwksError := sessmodels.GetCombinedJWKS()
	if jwksError != nil {
		supertokens.LogDebugMessage("getSessionHelper: Returning TryRefreshTokenError because there was an error fetching JWKs")
		if !defaultErrors.As(jwksError, &errors.TryRefreshTokenError{}) {
			return sessmodels.GetSessionResponse{}, jwksError
		}
	}

	accessTokenInfo, err = GetInfoFromAccessToken(parsedAccessToken, *combinedJwks, config.AntiCsrf == AntiCSRF_VIA_TOKEN && doAntiCsrfCheck)
	if err != nil {
		if !defaultErrors.As(err, &errors.TryRefreshTokenError{}) {
			supertokens.LogDebugMessage("getSessionHelper: Returning TryRefreshTokenError because GetInfoFromAccessToken returned an error")
			return sessmodels.GetSessionResponse{}, err
		}

		if parsedAccessToken.Version < 3 {
			payload := parsedAccessToken.Payload
			expiryTimeInPayload, expiryOk := payload["expiryTime"]
			timeCreatedInPayload, timeCreatedOk := payload["timeCreated"]

			if !expiryOk || !timeCreatedOk {
				return sessmodels.GetSessionResponse{}, err
			}

			expiryTime := uint64(expiryTimeInPayload.(float64))
			timeCreated := uint64(timeCreatedInPayload.(float64))

			if expiryTime < GetCurrTimeInMS() {
				return sessmodels.GetSessionResponse{}, err
			}

			// We check if the token was created since the last time we refreshed the keys from the core
			// Since we do not know the exact timing of the last refresh, we check against the max age
			if timeCreated <= (GetCurrTimeInMS() - uint64(sessmodels.JWKCacheMaxAgeInMs)) {
				return sessmodels.GetSessionResponse{}, err
			}
		} else {
			// Since v3 (and above) tokens contain a kid we can trust the cache-refresh mechanism of the jwt library.
			// This means we do not need to call the core since the signature wouldn't pass verification anyway.
			return sessmodels.GetSessionResponse{}, err
		}
	}

	if parsedAccessToken.Version >= 3 {
		tokenUsesDynamicKey := false
		kid := parsedAccessToken.KID

		if kid != nil && strings.HasPrefix(*kid, "d-") {
			tokenUsesDynamicKey = true
		}

		if tokenUsesDynamicKey != config.UseDynamicAccessTokenSigningKey {
			supertokens.LogDebugMessage("getSession: Returning TRY_REFRESH_TOKEN because the access token doesn't match the useDynamicAccessTokenSigningKey in the config")

			return sessmodels.GetSessionResponse{}, errors.TryRefreshTokenError{Msg: "The access token doesn't match the useDynamicAccessTokenSigningKey setting"}
		}
	}

	if doAntiCsrfCheck {
		if config.AntiCsrf == AntiCSRF_VIA_TOKEN {
			if accessTokenInfo != nil {
				if antiCsrfToken == nil || *antiCsrfToken != *accessTokenInfo.AntiCsrfToken {
					if antiCsrfToken == nil {
						supertokens.LogDebugMessage("getSession: Returning TRY_REFRESH_TOKEN because antiCsrfToken is missing from request")
						return sessmodels.GetSessionResponse{}, errors.TryRefreshTokenError{Msg: "Provided antiCsrfToken is undefined. If you do not want anti-csrf check for this API, please set doAntiCsrfCheck to false for this API"}
					} else {
						supertokens.LogDebugMessage("getSession: Returning TRY_REFRESH_TOKEN because the passed antiCsrfToken is not the same as in the access token")
						return sessmodels.GetSessionResponse{}, errors.TryRefreshTokenError{Msg: "anti-csrf check failed"}
					}
				}
			}
		} else if config.AntiCsrf == AntiCSRF_VIA_CUSTOM_HEADER {
			return sessmodels.GetSessionResponse{}, defaultErrors.New("Please either use VIA_TOKEN, NONE or call with doAntiCsrfCheck false")
		}
	}

	if accessTokenInfo != nil && !alwaysCheckCore && accessTokenInfo.ParentRefreshTokenHash1 == nil {
		return sessmodels.GetSessionResponse{
			Session: sessmodels.SessionStruct{
				Handle:                accessTokenInfo.SessionHandle,
				UserID:                accessTokenInfo.UserID,
				UserDataInAccessToken: accessTokenInfo.UserData,
				ExpiryTime:            accessTokenInfo.ExpiryTime,
			},
		}, nil
	}
	requestBody := map[string]interface{}{
		"accessToken":     parsedAccessToken.RawTokenString,
		"doAntiCsrfCheck": doAntiCsrfCheck,
		"enableAntiCsrf":  config.AntiCsrf == AntiCSRF_VIA_TOKEN,
		"checkDatabase":   alwaysCheckCore,
	}
	if antiCsrfToken != nil {
		requestBody["antiCsrfToken"] = *antiCsrfToken
	}

	if supertokens.IsRunningInTestMode() {
		didGetSessionCallCore = true
	}
	response, err := querier.SendPostRequest("/recipe/session/verify", requestBody)
	if err != nil {
		return sessmodels.GetSessionResponse{}, err
	}

	status := response["status"]
	if status.(string) == "OK" {
		delete(response, "status")

		responseByte, err := json.Marshal(response)
		if err != nil {
			return sessmodels.GetSessionResponse{}, err
		}
		var result sessmodels.GetSessionResponse
		err = json.Unmarshal(responseByte, &result)
		if err != nil {
			return sessmodels.GetSessionResponse{}, err
		}

		var expiryToSet uint64

		if result.AccessToken.Token != "" {
			expiryToSet = result.AccessToken.Expiry
		} else if accessTokenInfo != nil {
			expiryToSet = accessTokenInfo.ExpiryTime
		} else {
			expiryToSet = *sanitizeNumberInputAsUint64(parsedAccessToken.Payload["expiryTime"])
		}

		result.Session.ExpiryTime = expiryToSet
		return result, nil
	} else if response["status"].(string) == errors.UnauthorizedErrorStr {
		supertokens.LogDebugMessage("getSession: Returning UNAUTHORISED because of core response")
		return sessmodels.GetSessionResponse{}, errors.UnauthorizedError{Msg: response["message"].(string)}
	} else {
		supertokens.LogDebugMessage("getSession: Returning TRY_REFRESH_TOKEN because of core response")
		return sessmodels.GetSessionResponse{}, errors.TryRefreshTokenError{Msg: response["message"].(string)}
	}
}

func getSessionInformationHelper(querier supertokens.Querier, sessionHandle string) (*sessmodels.SessionInformation, error) {
	response, err := querier.SendGetRequest("/recipe/session",
		map[string]string{
			"sessionHandle": sessionHandle,
		})
	if err != nil {
		return nil, err
	}
	if response["status"] == "OK" {
		return &sessmodels.SessionInformation{
			SessionHandle:                    response["sessionHandle"].(string),
			UserId:                           response["userId"].(string),
			SessionDataInDatabase:            response["userDataInDatabase"].(map[string]interface{}),
			Expiry:                           uint64(response["expiry"].(float64)),
			TimeCreated:                      uint64(response["timeCreated"].(float64)),
			CustomClaimsInAccessTokenPayload: response["userDataInJWT"].(map[string]interface{}),
		}, nil
	}
	return nil, nil
}

func refreshSessionHelper(config sessmodels.TypeNormalisedInput, querier supertokens.Querier, refreshToken string, antiCsrfToken *string, disableAntiCsrf bool) (sessmodels.CreateOrRefreshAPIResponse, error) {
	requestBody := map[string]interface{}{
		"refreshToken":   refreshToken,
		"enableAntiCsrf": !disableAntiCsrf && config.AntiCsrf == AntiCSRF_VIA_TOKEN,
	}
	if antiCsrfToken != nil {
		requestBody["antiCsrfToken"] = *antiCsrfToken
	}

	if config.AntiCsrf == AntiCSRF_VIA_CUSTOM_HEADER && !disableAntiCsrf {
		return sessmodels.CreateOrRefreshAPIResponse{}, defaultErrors.New("Please either use VIA_TOKEN, NONE or call with doAntiCsrfCheck false")
	}

	response, err := querier.SendPostRequest("/recipe/session/refresh", requestBody)
	if err != nil {
		supertokens.LogDebugMessage("refreshSessionHelper: Call to /recipe/session/refresh API failed")
		return sessmodels.CreateOrRefreshAPIResponse{}, err
	}
	if response["status"] == "OK" {
		delete(response, "status")
		responseByte, err := json.Marshal(response)
		if err != nil {
			supertokens.LogDebugMessage("refreshSessionHelper: Could not parse response from /recipe/session/refresh API")
			return sessmodels.CreateOrRefreshAPIResponse{}, err
		}
		var result sessmodels.CreateOrRefreshAPIResponse
		err = json.Unmarshal(responseByte, &result)
		if err != nil {
			supertokens.LogDebugMessage("refreshSessionHelper: Could not decode response from /recipe/session/refresh API")
			return sessmodels.CreateOrRefreshAPIResponse{}, err
		}
		return result, nil
	} else if response["status"].(string) == errors.UnauthorizedErrorStr {
		supertokens.LogDebugMessage("refreshSession: Returning UNAUTHORISED because of core response")
		return sessmodels.CreateOrRefreshAPIResponse{}, errors.UnauthorizedError{Msg: response["message"].(string)}
	} else {
		sessionInfo := errors.TokenTheftDetectedErrorPayload{
			SessionHandle: (response["session"].(map[string]interface{}))["handle"].(string),
			UserID:        (response["session"].(map[string]interface{}))["userId"].(string),
		}

		supertokens.LogDebugMessage("refreshSession: Returning TOKEN_THEFT_DETECTED because of core response")
		return sessmodels.CreateOrRefreshAPIResponse{}, errors.TokenTheftDetectedError{
			Msg:     "Token theft detected",
			Payload: sessionInfo,
		}
	}
}

func revokeAllSessionsForUserHelper(querier supertokens.Querier, userID string) ([]string, error) {
	response, err := querier.SendPostRequest("/recipe/session/remove", map[string]interface{}{
		"userId": userID,
	})
	if err != nil {
		return nil, err
	}

	revokedSessionHandlesAsSliceOfInterfaces := response["sessionHandlesRevoked"].([]interface{})

	var result []string

	for _, val := range revokedSessionHandlesAsSliceOfInterfaces {
		result = append(result, val.(string))
	}

	return result, nil
}

func getAllSessionHandlesForUserHelper(querier supertokens.Querier, userID string) ([]string, error) {
	response, err := querier.SendGetRequest("/recipe/session/user", map[string]string{
		"userId": userID,
	})
	if err != nil {
		return nil, err
	}

	sessionHandlesAsSliceOfInterfaces := response["sessionHandles"].([]interface{})

	var result []string

	for _, val := range sessionHandlesAsSliceOfInterfaces {
		result = append(result, val.(string))
	}

	return result, nil
}

func revokeSessionHelper(querier supertokens.Querier, sessionHandle string) (bool, error) {
	response, err := querier.SendPostRequest("/recipe/session/remove",
		map[string]interface{}{
			"sessionHandles": [1]string{sessionHandle},
		})
	if err != nil {
		return false, err
	}
	return len(response["sessionHandlesRevoked"].([]interface{})) == 1, nil
}

func revokeMultipleSessionsHelper(querier supertokens.Querier, sessionHandles []string) ([]string, error) {
	response, err := querier.SendPostRequest("/recipe/session/remove",
		map[string]interface{}{
			"sessionHandles": sessionHandles,
		})
	if err != nil {
		return nil, err
	}
	revokedSessionHandlesAsSliceOfInterfaces := response["sessionHandlesRevoked"].([]interface{})

	var result []string

	for _, val := range revokedSessionHandlesAsSliceOfInterfaces {
		result = append(result, val.(string))
	}

	return result, nil
}

func updateSessionDataInDatabaseHelper(querier supertokens.Querier, sessionHandle string, newSessionData map[string]interface{}) (bool, error) {
	if newSessionData == nil {
		newSessionData = map[string]interface{}{}
	}
	response, err := querier.SendPutRequest("/recipe/session/data",
		map[string]interface{}{
			"sessionHandle":      sessionHandle,
			"userDataInDatabase": newSessionData,
		})
	if err != nil {
		return false, err
	}
	if response["status"].(string) == errors.UnauthorizedErrorStr {
		return false, nil
	}
	return true, nil
}

func updateAccessTokenPayloadHelper(querier supertokens.Querier, sessionHandle string, newAccessTokenPayload map[string]interface{}) (bool, error) {
	if newAccessTokenPayload == nil {
		newAccessTokenPayload = map[string]interface{}{}
	}
	response, err := querier.SendPutRequest("/recipe/jwt/data", map[string]interface{}{
		"sessionHandle": sessionHandle,
		"userDataInJWT": newAccessTokenPayload,
	})
	if err != nil {
		return false, err
	}
	if response["status"].(string) == errors.UnauthorizedErrorStr {
		return false, nil
	}
	return true, nil
}

func regenerateAccessTokenHelper(querier supertokens.Querier, newAccessTokenPayload *map[string]interface{}, accessToken string) (*sessmodels.RegenerateAccessTokenResponse, error) {
	if newAccessTokenPayload == nil {
		newAccessTokenPayload = &map[string]interface{}{}
	}
	response, err := querier.SendPostRequest("/recipe/session/regenerate", map[string]interface{}{
		"accessToken":   accessToken,
		"userDataInJWT": newAccessTokenPayload,
	})
	if err != nil {
		supertokens.LogDebugMessage("regenerateAccessTokenHelper: Call to /recipe/session/regenerate failed")
		return nil, err
	}
	if response["status"].(string) == errors.UnauthorizedErrorStr {
		return nil, nil
	}
	responseByte, err := json.Marshal(response)
	if err != nil {
		supertokens.LogDebugMessage("regenerateAccessTokenHelper: Failed to parse response from core")
		return nil, err
	}
	var resp sessmodels.RegenerateAccessTokenResponse
	err = json.Unmarshal(responseByte, &resp)
	if err != nil {
		supertokens.LogDebugMessage("regenerateAccessTokenHelper: Failed to decode response from core")
		return nil, err
	}
	return &resp, nil
}
