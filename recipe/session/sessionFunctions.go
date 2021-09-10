package session

import (
	"encoding/json"
	defaultErrors "errors"

	"github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func createNewSessionHelper(recipeImplHandshakeInfo *models.HandshakeInfo, config models.TypeNormalisedInput, querier supertokens.Querier, userID string, JWTPayload, sessionData interface{}) (models.CreateOrRefreshAPIResponse, error) {
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
		return models.CreateOrRefreshAPIResponse{}, err
	}
	requestBody["enableAntiCsrf"] = recipeImplHandshakeInfo.AntiCsrf == antiCSRF_VIA_TOKEN
	response, err := querier.SendPostRequest("/recipe/session", requestBody)
	if err != nil {
		return models.CreateOrRefreshAPIResponse{}, err
	}
	updateJwtSigningPublicKeyInfo(&recipeImplHandshakeInfo, response["jwtSigningPublicKey"].(string), uint64(response["jwtSigningPublicKeyExpiryTime"].(float64)))

	delete(response, "status")
	delete(response, "jwtSigningPublicKey")
	delete(response, "jwtSigningPublicKeyExpiryTime")

	responseByte, err := json.Marshal(response)
	if err != nil {
		return models.CreateOrRefreshAPIResponse{}, err
	}
	var resp models.CreateOrRefreshAPIResponse
	err = json.Unmarshal(responseByte, &resp)
	if err != nil {
		return models.CreateOrRefreshAPIResponse{}, err
	}
	return resp, nil
}

func getSessionHelper(recipeImplHandshakeInfo *models.HandshakeInfo, config models.TypeNormalisedInput, querier supertokens.Querier, accessToken string, antiCsrfToken *string, doAntiCsrfCheck, containsCustomHeader bool) (models.GetSessionResponse, error) {
	err := getHandshakeInfo(&recipeImplHandshakeInfo, config, querier, false)
	if err != nil {
		return models.GetSessionResponse{}, err
	}

	if recipeImplHandshakeInfo.JWTSigningPublicKeyExpiryTime > getCurrTimeInMS() && false {
		accessTokenInfo, err := getInfoFromAccessToken(accessToken, recipeImplHandshakeInfo.JWTSigningPublicKey, recipeImplHandshakeInfo.AntiCsrf == antiCSRF_VIA_TOKEN && doAntiCsrfCheck)
		if err != nil {
			if !defaultErrors.As(err, &errors.TryRefreshTokenError{}) {
				return models.GetSessionResponse{}, err
			}

			payload, errFromPayload := getPayloadWithoutVerifying(accessToken)

			if errFromPayload != nil {
				// we want to return the original error..
				return models.GetSessionResponse{}, err
			}

			expiryTime := uint64(payload["expiryTime"].(float64))
			timeCreated := uint64(payload["timeCreated"].(float64))

			if expiryTime < getCurrTimeInMS() {
				return models.GetSessionResponse{}, err
			}

			if recipeImplHandshakeInfo.SigningKeyLastUpdated > timeCreated {
				return models.GetSessionResponse{}, err
			}
		}

		if doAntiCsrfCheck {
			if recipeImplHandshakeInfo.AntiCsrf == antiCSRF_VIA_TOKEN {
				if accessTokenInfo != nil {
					if antiCsrfToken == nil || *antiCsrfToken == *accessTokenInfo.antiCsrfToken {
						if antiCsrfToken == nil {
							return models.GetSessionResponse{}, errors.TryRefreshTokenError{Msg: "Provided antiCsrfToken is undefined. If you do not want anti-csrf check for this API, please set doAntiCsrfCheck to false for this API"}
						} else {
							return models.GetSessionResponse{}, errors.TryRefreshTokenError{Msg: "anti-csrf check failed"}
						}
					}
				}
			} else if recipeImplHandshakeInfo.AntiCsrf == antiCSRF_VIA_CUSTOM_HEADER {
				if !containsCustomHeader {
					return models.GetSessionResponse{}, errors.TryRefreshTokenError{Msg: "anti-csrf check failed. Please pass 'rid: \"session\"' header in the request, or set doAntiCsrfCheck to false for this API"}
				}
			}
		}

		if accessTokenInfo != nil &&
			!recipeImplHandshakeInfo.AccessTokenBlacklistingEnabled &&
			accessTokenInfo.parentRefreshTokenHash1 == nil {
			return models.GetSessionResponse{
				Session: models.SessionStruct{
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
		return models.GetSessionResponse{}, err
	}

	status := response["status"]
	if status.(string) == "OK" {
		updateJwtSigningPublicKeyInfo(&recipeImplHandshakeInfo, response["jwtSigningPublicKey"].(string), uint64(response["jwtSigningPublicKeyExpiryTime"].(float64)))
		delete(response, "status")
		delete(response, "jwtSigningPublicKey")
		delete(response, "jwtSigningPublicKeyExpiryTime")
		responseByte, err := json.Marshal(response)
		if err != nil {
			return models.GetSessionResponse{}, err
		}
		var result models.GetSessionResponse
		err = json.Unmarshal(responseByte, &result)
		if err != nil {
			return models.GetSessionResponse{}, err
		}
		return result, nil
	} else if response["status"].(string) == errors.UnauthorizedErrorStr {
		return models.GetSessionResponse{}, errors.UnauthorizedError{Msg: response["message"].(string)}
	} else {
		jwtSigningPublicKey, jwtSigningPublicKeyExist := response["jwtSigningPublicKey"]
		jwtSigningPublicKeyExpiryTime, jwtSigningPublicKeyExpiryTimeExist := response["jwtSigningPublicKeyExpiryTime"]
		if jwtSigningPublicKeyExist && jwtSigningPublicKeyExpiryTimeExist {
			updateJwtSigningPublicKeyInfo(&recipeImplHandshakeInfo, jwtSigningPublicKey.(string), uint64(jwtSigningPublicKeyExpiryTime.(float64)))
		} else {
			// we ignore any errors produced by this function..
			getHandshakeInfo(&recipeImplHandshakeInfo, config, querier, true)
		}
		return models.GetSessionResponse{}, errors.TryRefreshTokenError{Msg: response["message"].(string)}
	}
}

func getSessionInformationHelper(querier supertokens.Querier, sessionHandle string) (models.SessionInformation, error) {
	response, err := querier.SendGetRequest("/recipe/session",
		map[string]interface{}{
			"sessionHandle": sessionHandle,
		})
	if err != nil {
		return models.SessionInformation{}, err
	}
	if response["status"] == "OK" {
		return models.SessionInformation{
			SessionHandle: response["sessionHandle"].(string),
			UserId:        response["userId"].(string),
			SessionData:   response["userDataInDatabase"].(map[string]interface{}),
			Expiry:        uint64(response["expiry"].(float64)),
			TimeCreated:   uint64(response["timeCreated"].(float64)),
			JwtPayload:    response["userDataInJWT"].(map[string]interface{}),
		}, nil
	}
	return models.SessionInformation{}, errors.UnauthorizedError{Msg: response["message"].(string)}
}

func refreshSessionHelper(recipeImplHandshakeInfo *models.HandshakeInfo, config models.TypeNormalisedInput, querier supertokens.Querier, refreshToken string, antiCsrfToken *string, containsCustomHeader bool) (models.CreateOrRefreshAPIResponse, error) {
	err := getHandshakeInfo(&recipeImplHandshakeInfo, config, querier, false)
	if err != nil {
		return models.CreateOrRefreshAPIResponse{}, err
	}

	if recipeImplHandshakeInfo.AntiCsrf == antiCSRF_VIA_CUSTOM_HEADER {
		if !containsCustomHeader {
			clearCookies := false
			return models.CreateOrRefreshAPIResponse{}, errors.UnauthorizedError{
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
		return models.CreateOrRefreshAPIResponse{}, err
	}
	if response["status"] == "OK" {
		delete(response, "status")
		responseByte, err := json.Marshal(response)
		if err != nil {
			return models.CreateOrRefreshAPIResponse{}, err
		}
		var result models.CreateOrRefreshAPIResponse
		err = json.Unmarshal(responseByte, &result)
		if err != nil {
			return models.CreateOrRefreshAPIResponse{}, err
		}
		return result, nil
	} else if response["status"].(string) == errors.UnauthorizedErrorStr {
		return models.CreateOrRefreshAPIResponse{}, errors.UnauthorizedError{Msg: response["message"].(string)}
	} else {
		sessionInfo := errors.TokenTheftDetectedErrorPayload{
			SessionHandle: (response["session"].(map[string]interface{}))["handle"].(string),
			UserID:        (response["session"].(map[string]interface{}))["userId"].(string),
		}
		return models.CreateOrRefreshAPIResponse{}, errors.TokenTheftDetectedError{
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
	response, err := querier.SendGetRequest("/recipe/session/user", map[string]interface{}{
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

func updateSessionDataHelper(querier supertokens.Querier, sessionHandle string, newSessionData interface{}) error {
	if newSessionData == nil {
		newSessionData = map[string]string{}
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

func updateJWTPayloadHelper(querier supertokens.Querier, sessionHandle string, newJWTPayload interface{}) error {
	if newJWTPayload == nil {
		newJWTPayload = map[string]string{}
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
