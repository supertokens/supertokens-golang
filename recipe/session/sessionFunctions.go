package session

import (
	"encoding/json"
	"fmt"

	"github.com/supertokens/supertokens-golang/errors"
	sessionErrors "github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func CreateNewSession(querier supertokens.Querier, userID string, JWTPayload interface{}, sessionData interface{}) (models.CreateOrRefreshAPIResponse, error) {
	URL, err := supertokens.NewNormalisedURLPath("/recipe/session")
	if err != nil {
		return models.CreateOrRefreshAPIResponse{}, err
	}
	requestBody := map[string]interface{}{
		"userId":             userID,
		"userDataInJWT":      JWTPayload,
		"userDataInDatabase": sessionData,
	}
	handShakeInfo, err := GetHandshakeInfo(querier)
	if err != nil {
		return models.CreateOrRefreshAPIResponse{}, err
	}
	requestBody["enableAntiCsrf"] = handShakeInfo.AntiCsrf == AntiCSRF_VIA_TOKEN
	response, err := querier.SendPostRequest(*URL, requestBody)
	if err != nil {
		return models.CreateOrRefreshAPIResponse{}, err
	}
	UpdateJwtSigningPublicKeyInfo(response["jwtSigningPublicKey"].(string), response["jwtSigningPublicKeyExpiryTime"].(uint64))

	delete(response, "status")
	delete(response, "jwtSigningPublicKey")
	delete(response, "jwtSigningPublicKeyExpiryTime")

	var resp models.CreateOrRefreshAPIResponse
	bytes := []byte(fmt.Sprintf("%+v", response))
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return models.CreateOrRefreshAPIResponse{}, err
	}
	return resp, nil
}

func GetSession(querier supertokens.Querier, accessToken string, antiCsrfToken *string, doAntiCsrfCheck bool, containsCustomHeader bool) (models.GetSessionResponse, error) {
	handShakeInfo, err := GetHandshakeInfo(querier)
	if err != nil {
		return models.GetSessionResponse{}, err
	}
	if handShakeInfo.JWTSigningPublicKeyExpiryTime > getCurrTimeInMS() {
		accessTokenInfo, err := getInfoFromAccessToken(accessToken, handShakeInfo.JWTSigningPublicKey, handShakeInfo.AntiCsrf == AntiCSRF_VIA_TOKEN && doAntiCsrfCheck)
		if err != nil {
			return models.GetSessionResponse{}, err
		}
		if handShakeInfo.AntiCsrf == AntiCSRF_VIA_TOKEN && doAntiCsrfCheck {
			if antiCsrfToken == nil || *antiCsrfToken == *accessTokenInfo.antiCsrfToken {
				if antiCsrfToken == nil {
					return models.GetSessionResponse{}, errors.BadInputError{Msg: "Provided antiCsrfToken is undefined. If you do not want anti-csrf check for this API, please set doAntiCsrfCheck to false for this API"}
				} else {
					return models.GetSessionResponse{}, errors.BadInputError{Msg: "anti-csrf check failed"}
				}
			}
		} else if handShakeInfo.AntiCsrf == AntiCSRF_VIA_CUSTOM_HEADER && doAntiCsrfCheck {
			if !containsCustomHeader {
				return models.GetSessionResponse{}, errors.BadInputError{Msg: "anti-csrf check failed. Please pass 'rid: \"session\"' header in the request, or set doAntiCsrfCheck to false for this API"}
			}
		}

		if !handShakeInfo.AccessTokenBlacklistingEnabled && accessTokenInfo.parentRefreshTokenHash1 == nil {
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
		"antiCsrfToken":   antiCsrfToken,
		"doAntiCsrfCheck": doAntiCsrfCheck,
		"enableAntiCsrf":  handShakeInfo.AntiCsrf == AntiCSRF_VIA_TOKEN,
	}
	path, err := supertokens.NewNormalisedURLPath("/recipe/session/verify")
	if err != nil {
		return models.GetSessionResponse{}, err
	}
	response, err := querier.SendPostRequest(*path, requestBody)
	if err != nil {
		return models.GetSessionResponse{}, err
	}
	if response["status"] == "OK" {
		UpdateJwtSigningPublicKeyInfo(response["jwtSigningPublicKey"].(string), response["jwtSigningPublicKeyExpiryTime"].(uint64))
		delete(response, "status")
		delete(response, "jwtSigningPublicKey")
		delete(response, "jwtSigningPublicKeyExpiryTime")
		var result models.GetSessionResponse
		err := json.Unmarshal([]byte(fmt.Sprintf("%+v", response)), &result)
		if err != nil {
			return models.GetSessionResponse{}, err
		}
		return result, nil
	} else if response["status"] == "UNAUTHORISED" {
		return models.GetSessionResponse{}, sessionErrors.MakeUnauthorizedError(response["message"].(string))
	} else {
		return models.GetSessionResponse{}, sessionErrors.MakeTryRefreshTokenError(response["message"].(string))
	}
}

func refreshSession(querier supertokens.Querier, refreshToken string, antiCsrfToken *string, containsCustomHeader bool) (models.CreateOrRefreshAPIResponse, error) {
	handShakeInfo, err := GetHandshakeInfo(querier)
	if err != nil {
		return models.CreateOrRefreshAPIResponse{}, err
	}

	if handShakeInfo.AntiCsrf == AntiCSRF_VIA_CUSTOM_HEADER {
		if !containsCustomHeader {
			return models.CreateOrRefreshAPIResponse{}, sessionErrors.MakeUnauthorizedError("anti-csrf check failed. Please pass 'rid: \"session\"' header in the request.")
		}
	}
	path, err := supertokens.NewNormalisedURLPath("/recipe/session/refresh")
	if err != nil {
		return models.CreateOrRefreshAPIResponse{}, err
	}
	requestBody := map[string]interface{}{
		"refreshToken":   refreshToken,
		"antiCsrfToken":  antiCsrfToken,
		"enableAntiCsrf": handShakeInfo.AntiCsrf == AntiCSRF_VIA_TOKEN,
	}
	response, err := querier.SendPostRequest(*path, requestBody)
	if response["status"] == "OK" {
		delete(response, "status")
		var result models.CreateOrRefreshAPIResponse
		err := json.Unmarshal([]byte(fmt.Sprintf("%+v", response)), &result)
		if err != nil {
			return models.CreateOrRefreshAPIResponse{}, err
		}
		return result, nil
	} else if response["status"] == "UNAUTHORISED" {
		return models.CreateOrRefreshAPIResponse{}, sessionErrors.MakeUnauthorizedError(response["message"].(string))
	} else {
		session := response["session"].(sessionErrors.TokenTheftDetectedErrorPayload)
		return models.CreateOrRefreshAPIResponse{}, sessionErrors.MakeTokenTheftDetectedError(session.SessionHandle, session.UserID, "Token theft detected")
	}
}

func revokeAllSessionsForUser(querier supertokens.Querier, userID string) ([]string, error) {
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

func getAllSessionHandlesForUser(querier supertokens.Querier, userID string) ([]string, error) {
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

// revokeSession function used to revoke a specific session
func revokeSession(querier supertokens.Querier, sessionHandle string) (bool, error) {
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

func revokeMultipleSessions(querier supertokens.Querier, sessionHandles []string) ([]string, error) {
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

func getSessionData(querier supertokens.Querier, sessionHandle string) (map[string]interface{}, error) {
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

func updateSessionData(querier supertokens.Querier, sessionHandle string, newSessionData interface{}) error {
	path, err := supertokens.NewNormalisedURLPath("/recipe/session/data")
	if err != nil {
		return err
	}
	response, err := querier.SendPutRequest(*path,
		map[string]interface{}{
			"sessionHandle":      sessionHandle,
			"userDataInDatabase": newSessionData,
		})
	if response["status"] == "UNAUTHORISED" {
		return sessionErrors.MakeUnauthorizedError(response["message"].(string))
	}
	return nil
}

func getJWTPayload(querier supertokens.Querier, sessionHandle string) (interface{}, error) {
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

func updateJWTPayload(querier supertokens.Querier, sessionHandle string, newJWTPayload interface{}) error {
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
	return sessionErrors.MakeUnauthorizedError(response["message"].(string))
}
