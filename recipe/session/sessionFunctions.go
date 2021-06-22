package session

import (
	"encoding/json"
	"fmt"

	sessionErrors "github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func createNewSessionHelper(recipeImplHandshakeInfo *models.HandshakeInfo, config models.TypeNormalisedInput, querier supertokens.Querier, userID string, JWTPayload, sessionData interface{}) (*models.CreateOrRefreshAPIResponse, error) {
	URL, err := supertokens.NewNormalisedURLPath("/recipe/session")
	if err != nil {
		return nil, err
	}
	// TODO: do a POC with a nil value for sessionData and JWTPayload
	requestBody := map[string]interface{}{
		"userId":             userID,
		"userDataInJWT":      JWTPayload,
		"userDataInDatabase": sessionData,
	}
	handShakeInfo, err := GetHandshakeInfo(recipeImplHandshakeInfo, config, querier)
	if err != nil {
		return nil, err
	}
	requestBody["enableAntiCsrf"] = handShakeInfo.AntiCsrf == antiCSRF_VIA_TOKEN
	response, err := querier.SendPostRequest(*URL, requestBody)
	if err != nil {
		return nil, err
	}
	UpdateJwtSigningPublicKeyInfo(recipeImplHandshakeInfo, response["jwtSigningPublicKey"].(string), response["jwtSigningPublicKeyExpiryTime"].(uint64))

	delete(response, "status")
	delete(response, "jwtSigningPublicKey")
	delete(response, "jwtSigningPublicKeyExpiryTime")

	var resp models.CreateOrRefreshAPIResponse

	bytes := []byte(fmt.Sprintf("%+v", response))
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func getSessionHelper(recipeImplHandshakeInfo *models.HandshakeInfo, config models.TypeNormalisedInput, querier supertokens.Querier, accessToken string, antiCsrfToken *string, doAntiCsrfCheck, containsCustomHeader bool) (*models.GetSessionResponse, error) {
	handShakeInfo, err := GetHandshakeInfo(recipeImplHandshakeInfo, config, querier)
	if err != nil {
		return nil, err
	}
	if handShakeInfo.JWTSigningPublicKeyExpiryTime > getCurrTimeInMS() {
		accessTokenInfo, err := getInfoFromAccessToken(accessToken, handShakeInfo.JWTSigningPublicKey, handShakeInfo.AntiCsrf == antiCSRF_VIA_TOKEN && doAntiCsrfCheck)
		if err != nil {
			if !sessionErrors.IsTryRefreshTokenError(err) {
				return nil, err
			}

			payload, errFromPayload := getPayloadWithoutVerifying(accessToken)

			if errFromPayload != nil {
				// we want to return the original error..
				return nil, err
			}

			expiryTime := uint64(payload["expiryTime"].(float64))
			timeCreated := uint64(payload["timeCreated"].(float64))

			if expiryTime < getCurrTimeInMS() {
				return nil, err
			}

			if handShakeInfo.SigningKeyLastUpdated > timeCreated {
				return nil, err
			}
		}

		if doAntiCsrfCheck {
			if handShakeInfo.AntiCsrf == antiCSRF_VIA_TOKEN {
				if accessTokenInfo != nil {
					if antiCsrfToken == nil || *antiCsrfToken == *accessTokenInfo.antiCsrfToken {
						if antiCsrfToken == nil {
							return nil, sessionErrors.TryRefreshTokenError{Msg: "Provided antiCsrfToken is undefined. If you do not want anti-csrf check for this API, please set doAntiCsrfCheck to false for this API"}
						} else {
							return nil, sessionErrors.TryRefreshTokenError{Msg: "anti-csrf check failed"}
						}
					}
				}
			} else if handShakeInfo.AntiCsrf == antiCSRF_VIA_CUSTOM_HEADER {
				if !containsCustomHeader {
					return nil, sessionErrors.TryRefreshTokenError{Msg: "anti-csrf check failed. Please pass 'rid: \"session\"' header in the request, or set doAntiCsrfCheck to false for this API"}
				}
			}
		}

		if accessTokenInfo != nil &&
			!handShakeInfo.AccessTokenBlacklistingEnabled &&
			accessTokenInfo.parentRefreshTokenHash1 == nil {
			return &models.GetSessionResponse{
				Session: models.SessionStruct{
					Handle:        accessTokenInfo.sessionHandle,
					UserID:        accessTokenInfo.userID,
					UserDataInJWT: accessTokenInfo.userData,
				},
				AccessToken: nil,
			}, nil
		}
	}
	requestBody := map[string]interface{}{
		"accessToken":     accessToken,
		"antiCsrfToken":   antiCsrfToken,
		"doAntiCsrfCheck": doAntiCsrfCheck,
		"enableAntiCsrf":  handShakeInfo.AntiCsrf == antiCSRF_VIA_TOKEN,
	}
	path, err := supertokens.NewNormalisedURLPath("/recipe/session/verify")
	if err != nil {
		return nil, err
	}
	response, err := querier.SendPostRequest(*path, requestBody)
	if err != nil {
		return nil, err
	}
	if response["status"] == "OK" {
		UpdateJwtSigningPublicKeyInfo(recipeImplHandshakeInfo, response["jwtSigningPublicKey"].(string), response["jwtSigningPublicKeyExpiryTime"].(uint64))
		delete(response, "status")
		delete(response, "jwtSigningPublicKey")
		delete(response, "jwtSigningPublicKeyExpiryTime")
		// TODO: For POC: check if pointer is unmarshalled
		var result models.GetSessionResponse
		err := json.Unmarshal([]byte(fmt.Sprintf("%+v", response)), &result)
		if err != nil {
			return nil, err
		}
		return &result, nil
	} else if response["status"] == "UNAUTHORISED" {
		return nil, sessionErrors.MakeUnauthorizedError(response["message"].(string))
	} else {
		// TODO: check using "ok" method instead of checking fir nil directly.
		if response["jwtSigningPublicKey"] != nil && response["jwtSigningPublicKeyExpiryTime"] != nil {
			UpdateJwtSigningPublicKeyInfo(recipeImplHandshakeInfo, response["jwtSigningPublicKey"].(string), response["jwtSigningPublicKeyExpiryTime"].(uint64))
		}
		return nil, sessionErrors.MakeTryRefreshTokenError(response["message"].(string))
	}
}

func refreshSessionHelper(recipeImplHandshakeInfo *models.HandshakeInfo, config models.TypeNormalisedInput, querier supertokens.Querier, refreshToken string, antiCsrfToken *string, containsCustomHeader bool) (*models.CreateOrRefreshAPIResponse, error) {
	handShakeInfo, err := GetHandshakeInfo(recipeImplHandshakeInfo, config, querier)
	if err != nil {
		return nil, err
	}

	if handShakeInfo.AntiCsrf == antiCSRF_VIA_CUSTOM_HEADER {
		if !containsCustomHeader {
			return nil, sessionErrors.MakeUnauthorizedError("anti-csrf check failed. Please pass 'rid: \"session\"' header in the request.")
		}
	}
	path, err := supertokens.NewNormalisedURLPath("/recipe/session/refresh")
	if err != nil {
		return nil, err
	}
	requestBody := map[string]interface{}{
		"refreshToken":   refreshToken,
		"antiCsrfToken":  antiCsrfToken,
		"enableAntiCsrf": handShakeInfo.AntiCsrf == antiCSRF_VIA_TOKEN,
	}
	response, err := querier.SendPostRequest(*path, requestBody)
	if err != nil {
		return nil, err
	}
	if response["status"] == "OK" {
		delete(response, "status")
		var result models.CreateOrRefreshAPIResponse
		err := json.Unmarshal([]byte(fmt.Sprintf("%+v", response)), &result)
		if err != nil {
			return nil, err
		}
		return &result, nil
	} else if response["status"] == "UNAUTHORISED" {
		return nil, sessionErrors.MakeUnauthorizedError(response["message"].(string))
	} else {
		session := response["session"].(sessionErrors.TokenTheftDetectedErrorPayload)
		return nil, sessionErrors.MakeTokenTheftDetectedError(session.SessionHandle, session.UserID, "Token theft detected")
	}
}

func revokeAllSessionsForUserHelper(querier supertokens.Querier, userID string) ([]string, error) {
	path, err := supertokens.NewNormalisedURLPath("/recipe/session/remove")
	if err != nil {
		return nil, err
	}
	response, err := querier.SendPostRequest(*path, map[string]interface{}{
		"userId": userID,
	})
	if err != nil {
		return nil, err
	}
	return response["sessionHandlesRevoked"].([]string), nil
}

func getAllSessionHandlesForUserHelper(querier supertokens.Querier, userID string) ([]string, error) {
	path, err := supertokens.NewNormalisedURLPath("/recipe/session/user")
	if err != nil {
		return nil, err
	}
	response, err := querier.SendPostRequest(*path, map[string]interface{}{
		"userId": userID,
	})
	if err != nil {
		return nil, err
	}
	return response["sessionHandles"].([]string), nil
}

func revokeSessionHelper(querier supertokens.Querier, sessionHandle string) (bool, error) {
	path, err := supertokens.NewNormalisedURLPath("/recipe/session/remove")
	if err != nil {
		return false, err
	}
	response, err := querier.SendPostRequest(*path,
		map[string]interface{}{
			"sessionHandles": [1]string{sessionHandle},
		})
	if err != nil {
		return false, err
	}
	return len(response["sessionHandlesRevoked"].([]interface{})) == 1, nil
}

func revokeMultipleSessionsHelper(querier supertokens.Querier, sessionHandles []string) ([]string, error) {
	path, err := supertokens.NewNormalisedURLPath("/recipe/session/remove")
	if err != nil {
		return nil, err
	}
	response, err := querier.SendPostRequest(*path,
		map[string]interface{}{
			"sessionHandles": sessionHandles,
		})
	if err != nil {
		return nil, err
	}
	return response["sessionHandlesRevoked"].([]string), nil
}

func getSessionDataHelper(querier supertokens.Querier, sessionHandle string) (map[string]interface{}, error) {
	path, err := supertokens.NewNormalisedURLPath("/recipe/session/data")
	if err != nil {
		return nil, err
	}
	response, err := querier.SendGetRequest(*path,
		map[string]interface{}{
			"sessionHandle": sessionHandle,
		})
	if err != nil {
		return nil, err
	}
	if response["status"] == "OK" {
		return response["userDataInDatabase"].(map[string]interface{}), nil
	}
	return nil, sessionErrors.MakeUnauthorizedError(response["message"].(string))
}

func updateSessionDataHelper(querier supertokens.Querier, sessionHandle string, newSessionData interface{}) error {
	if newSessionData == nil {
		newSessionData = map[string]string{}
	}
	path, err := supertokens.NewNormalisedURLPath("/recipe/session/data")
	if err != nil {
		return err
	}
	response, err := querier.SendPutRequest(*path,
		map[string]interface{}{
			"sessionHandle":      sessionHandle,
			"userDataInDatabase": newSessionData,
		})
	if err != nil {
		return err
	}
	if response["status"] == "UNAUTHORISED" {
		return sessionErrors.MakeUnauthorizedError(response["message"].(string))
	}
	return nil
}

func getJWTPayloadHelper(querier supertokens.Querier, sessionHandle string) (interface{}, error) {
	path, err := supertokens.NewNormalisedURLPath("/recipe/jwt/data")
	if err != nil {
		return nil, err
	}
	response, err := querier.SendGetRequest(*path, map[string]interface{}{
		"sessionHandle": sessionHandle,
	})
	if err != nil {
		return nil, err
	}
	if response["status"] == "OK" {
		return response["userDataInJWT"], nil
	}
	return nil, sessionErrors.MakeUnauthorizedError(response["message"].(string))
}

func updateJWTPayloadHelper(querier supertokens.Querier, sessionHandle string, newJWTPayload interface{}) error {
	if newJWTPayload == nil {
		newJWTPayload = map[string]string{}
	}
	path, err := supertokens.NewNormalisedURLPath("/recipe/jwt/data")
	if err != nil {
		return err
	}
	response, err := querier.SendPutRequest(*path, map[string]interface{}{
		"sessionHandle": sessionHandle,
		"userDataInJWT": newJWTPayload,
	})
	if err != nil {
		return err
	}
	if response["status"] == "UNAUTHORISED" {
		return sessionErrors.MakeUnauthorizedError(response["message"].(string))
	}
	return nil
}
