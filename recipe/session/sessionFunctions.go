package session

import (
	"encoding/json"
	"fmt"

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
	handShakeInfo := GetHandshakeInfo()
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

func GetSession(recipeImplementation models.RecipeImplementation, accessToken string, antiCsrfToken *string, doAntiCsrfCheck bool, scontainsCustomHeader bool) {

}

// RevokeSession function used to revoke a specific session
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
	return nil, UnauthorizedError{
		Msg: response["message"].(string),
	}
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
		return UnauthorizedError{
			Msg: response["message"].(string),
		}
	}
	return nil
}
