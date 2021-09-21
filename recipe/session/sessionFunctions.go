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

	"github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func createNewSessionHelper(recipeImplHandshakeInfo *sessmodels.HandshakeInfo, config sessmodels.TypeNormalisedInput, querier supertokens.Querier, userID string, JWTPayload, sessionData map[string]interface{}) (sessmodels.CreateOrRefreshAPIResponse, error) {
	if JWTPayload == nil {
		JWTPayload = map[string]interface{}{}
	}
	if sessionData == nil {
		sessionData = map[string]interface{}{}
	}
	requestBody := map[string]interface{}{
		"userId":             userID,
		"userDataInJWT":      JWTPayload,
		"userDataInDatabase": sessionData,
	}
	err := getHandshakeInfo(&recipeImplHandshakeInfo, config, querier, false)
	if err != nil {
		return sessmodels.CreateOrRefreshAPIResponse{}, err
	}
	requestBody["enableAntiCsrf"] = recipeImplHandshakeInfo.AntiCsrf == antiCSRF_VIA_TOKEN
	response, err := querier.SendPostRequest("/recipe/session", requestBody)
	if err != nil {
		return sessmodels.CreateOrRefreshAPIResponse{}, err
	}
	updateJwtSigningPublicKeyInfo(&recipeImplHandshakeInfo, response["jwtSigningPublicKey"].(string), uint64(response["jwtSigningPublicKeyExpiryTime"].(float64)))

	delete(response, "status")
	delete(response, "jwtSigningPublicKey")
	delete(response, "jwtSigningPublicKeyExpiryTime")

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

func getSessionHelper(recipeImplHandshakeInfo *sessmodels.HandshakeInfo, config sessmodels.TypeNormalisedInput, querier supertokens.Querier, accessToken string, antiCsrfToken *string, doAntiCsrfCheck, containsCustomHeader bool) (sessmodels.GetSessionResponse, error) {
	err := getHandshakeInfo(&recipeImplHandshakeInfo, config, querier, false)
	if err != nil {
		return sessmodels.GetSessionResponse{}, err
	}

	if recipeImplHandshakeInfo.JWTSigningPublicKeyExpiryTime > getCurrTimeInMS() {
		accessTokenInfo, err := getInfoFromAccessToken(accessToken, recipeImplHandshakeInfo.JWTSigningPublicKey, recipeImplHandshakeInfo.AntiCsrf == antiCSRF_VIA_TOKEN && doAntiCsrfCheck)
		if err != nil {
			if !defaultErrors.As(err, &errors.TryRefreshTokenError{}) {
				return sessmodels.GetSessionResponse{}, err
			}

			payload, errFromPayload := getPayloadWithoutVerifying(accessToken)

			if errFromPayload != nil {
				// we want to return the original error..
				return sessmodels.GetSessionResponse{}, err
			}

			expiryTime := uint64(payload["expiryTime"].(float64))
			timeCreated := uint64(payload["timeCreated"].(float64))

			if expiryTime < getCurrTimeInMS() {
				return sessmodels.GetSessionResponse{}, err
			}

			if recipeImplHandshakeInfo.SigningKeyLastUpdated > timeCreated {
				return sessmodels.GetSessionResponse{}, err
			}
		}

		if doAntiCsrfCheck {
			if recipeImplHandshakeInfo.AntiCsrf == antiCSRF_VIA_TOKEN {
				if accessTokenInfo != nil {
					if antiCsrfToken == nil || *antiCsrfToken != *accessTokenInfo.antiCsrfToken {
						if antiCsrfToken == nil {
							return sessmodels.GetSessionResponse{}, errors.TryRefreshTokenError{Msg: "Provided antiCsrfToken is undefined. If you do not want anti-csrf check for this API, please set doAntiCsrfCheck to false for this API"}
						} else {
							return sessmodels.GetSessionResponse{}, errors.TryRefreshTokenError{Msg: "anti-csrf check failed"}
						}
					}
				}
			} else if recipeImplHandshakeInfo.AntiCsrf == antiCSRF_VIA_CUSTOM_HEADER {
				if !containsCustomHeader {
					return sessmodels.GetSessionResponse{}, errors.TryRefreshTokenError{Msg: "anti-csrf check failed. Please pass 'rid: \"session\"' header in the request, or set doAntiCsrfCheck to false for this API"}
				}
			}
		}

		if accessTokenInfo != nil &&
			!recipeImplHandshakeInfo.AccessTokenBlacklistingEnabled &&
			accessTokenInfo.parentRefreshTokenHash1 == nil {
			return sessmodels.GetSessionResponse{
				Session: sessmodels.SessionStruct{
					Handle:        accessTokenInfo.sessionHandle,
					UserID:        accessTokenInfo.userID,
					UserDataInJWT: accessTokenInfo.userData,
				},
			}, nil
		}
	}
	requestBody := map[string]interface{}{
		"accessToken":     accessToken,
		"doAntiCsrfCheck": doAntiCsrfCheck,
		"enableAntiCsrf":  recipeImplHandshakeInfo.AntiCsrf == antiCSRF_VIA_TOKEN,
	}
	if antiCsrfToken != nil {
		requestBody["antiCsrfToken"] = *antiCsrfToken
	}

	response, err := querier.SendPostRequest("/recipe/session/verify", requestBody)
	if err != nil {
		return sessmodels.GetSessionResponse{}, err
	}

	status := response["status"]
	if status.(string) == "OK" {
		updateJwtSigningPublicKeyInfo(&recipeImplHandshakeInfo, response["jwtSigningPublicKey"].(string), uint64(response["jwtSigningPublicKeyExpiryTime"].(float64)))
		delete(response, "status")
		delete(response, "jwtSigningPublicKey")
		delete(response, "jwtSigningPublicKeyExpiryTime")
		responseByte, err := json.Marshal(response)
		if err != nil {
			return sessmodels.GetSessionResponse{}, err
		}
		var result sessmodels.GetSessionResponse
		err = json.Unmarshal(responseByte, &result)
		if err != nil {
			return sessmodels.GetSessionResponse{}, err
		}
		return result, nil
	} else if response["status"].(string) == errors.UnauthorizedErrorStr {
		return sessmodels.GetSessionResponse{}, errors.UnauthorizedError{Msg: response["message"].(string)}
	} else {
		jwtSigningPublicKey, jwtSigningPublicKeyExist := response["jwtSigningPublicKey"]
		jwtSigningPublicKeyExpiryTime, jwtSigningPublicKeyExpiryTimeExist := response["jwtSigningPublicKeyExpiryTime"]
		if jwtSigningPublicKeyExist && jwtSigningPublicKeyExpiryTimeExist {
			updateJwtSigningPublicKeyInfo(&recipeImplHandshakeInfo, jwtSigningPublicKey.(string), uint64(jwtSigningPublicKeyExpiryTime.(float64)))
		} else {
			// we ignore any errors produced by this function..
			getHandshakeInfo(&recipeImplHandshakeInfo, config, querier, true)
		}
		return sessmodels.GetSessionResponse{}, errors.TryRefreshTokenError{Msg: response["message"].(string)}
	}
}

func getSessionInformationHelper(querier supertokens.Querier, sessionHandle string) (sessmodels.SessionInformation, error) {
	response, err := querier.SendGetRequest("/recipe/session",
		map[string]string{
			"sessionHandle": sessionHandle,
		})
	if err != nil {
		return sessmodels.SessionInformation{}, err
	}
	if response["status"] == "OK" {
		return sessmodels.SessionInformation{
			SessionHandle: response["sessionHandle"].(string),
			UserId:        response["userId"].(string),
			SessionData:   response["userDataInDatabase"].(map[string]interface{}),
			Expiry:        uint64(response["expiry"].(float64)),
			TimeCreated:   uint64(response["timeCreated"].(float64)),
			JwtPayload:    response["userDataInJWT"].(map[string]interface{}),
		}, nil
	}
	return sessmodels.SessionInformation{}, errors.UnauthorizedError{Msg: response["message"].(string)}
}

func refreshSessionHelper(recipeImplHandshakeInfo *sessmodels.HandshakeInfo, config sessmodels.TypeNormalisedInput, querier supertokens.Querier, refreshToken string, antiCsrfToken *string, containsCustomHeader bool) (sessmodels.CreateOrRefreshAPIResponse, error) {
	err := getHandshakeInfo(&recipeImplHandshakeInfo, config, querier, false)
	if err != nil {
		return sessmodels.CreateOrRefreshAPIResponse{}, err
	}

	if recipeImplHandshakeInfo.AntiCsrf == antiCSRF_VIA_CUSTOM_HEADER {
		if !containsCustomHeader {
			clearCookies := false
			return sessmodels.CreateOrRefreshAPIResponse{}, errors.UnauthorizedError{
				Msg:          "anti-csrf check failed. Please pass 'rid: \"session\"' header in the request.",
				ClearCookies: &clearCookies,
			}
		}
	}

	requestBody := map[string]interface{}{
		"refreshToken":   refreshToken,
		"enableAntiCsrf": recipeImplHandshakeInfo.AntiCsrf == antiCSRF_VIA_TOKEN,
	}
	if antiCsrfToken != nil {
		requestBody["antiCsrfToken"] = *antiCsrfToken
	}
	response, err := querier.SendPostRequest("/recipe/session/refresh", requestBody)
	if err != nil {
		return sessmodels.CreateOrRefreshAPIResponse{}, err
	}
	if response["status"] == "OK" {
		delete(response, "status")
		responseByte, err := json.Marshal(response)
		if err != nil {
			return sessmodels.CreateOrRefreshAPIResponse{}, err
		}
		var result sessmodels.CreateOrRefreshAPIResponse
		err = json.Unmarshal(responseByte, &result)
		if err != nil {
			return sessmodels.CreateOrRefreshAPIResponse{}, err
		}
		return result, nil
	} else if response["status"].(string) == errors.UnauthorizedErrorStr {
		return sessmodels.CreateOrRefreshAPIResponse{}, errors.UnauthorizedError{Msg: response["message"].(string)}
	} else {
		sessionInfo := errors.TokenTheftDetectedErrorPayload{
			SessionHandle: (response["session"].(map[string]interface{}))["handle"].(string),
			UserID:        (response["session"].(map[string]interface{}))["userId"].(string),
		}
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
	return response["sessionHandlesRevoked"].([]string), nil
}

func getAllSessionHandlesForUserHelper(querier supertokens.Querier, userID string) ([]string, error) {
	response, err := querier.SendGetRequest("/recipe/session/user", map[string]string{
		"userId": userID,
	})
	if err != nil {
		return nil, err
	}
	return response["sessionHandles"].([]string), nil
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
	return response["sessionHandlesRevoked"].([]string), nil
}

func updateSessionDataHelper(querier supertokens.Querier, sessionHandle string, newSessionData map[string]interface{}) error {
	if newSessionData == nil {
		newSessionData = map[string]interface{}{}
	}
	response, err := querier.SendPutRequest("/recipe/session/data",
		map[string]interface{}{
			"sessionHandle":      sessionHandle,
			"userDataInDatabase": newSessionData,
		})
	if err != nil {
		return err
	}
	if response["status"].(string) == errors.UnauthorizedErrorStr {
		return errors.UnauthorizedError{Msg: response["message"].(string)}
	}
	return nil
}

func updateJWTPayloadHelper(querier supertokens.Querier, sessionHandle string, newJWTPayload map[string]interface{}) error {
	if newJWTPayload == nil {
		newJWTPayload = map[string]interface{}{}
	}
	response, err := querier.SendPutRequest("/recipe/jwt/data", map[string]interface{}{
		"sessionHandle": sessionHandle,
		"userDataInJWT": newJWTPayload,
	})
	if err != nil {
		return err
	}
	if response["status"].(string) == errors.UnauthorizedErrorStr {
		return errors.UnauthorizedError{Msg: response["message"].(string)}
	}
	return nil
}
