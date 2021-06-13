package session

import (
	"encoding/json"
	"fmt"

	"github.com/supertokens/supertokens-golang/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/schema"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func CreateNewSession(recipeImplementation schema.RecipeImplementation, userID string, JWTPayload interface{}, sessionData interface{}) (schema.CreateOrRefreshAPIResponse, error) {
	URL, err := supertokens.NewNormalisedURLPath("/recipe/session")
	if err != nil {
		return schema.CreateOrRefreshAPIResponse{}, err
	}
	requestBody := map[string]interface{}{
		"userId":             userID,
		"userDataInJWT":      JWTPayload,
		"userDataInDatabase": sessionData,
	}
	handShakeInfo := recipeImplementation.GetHandshakeInfo()
	requestBody["enableAntiCsrf"] = handShakeInfo.AntiCsrf == AntiCSRF_VIA_TOKEN
	response, err := recipeImplementation.Querier.SendPostRequest(*URL, requestBody)
	if err != nil {
		return schema.CreateOrRefreshAPIResponse{}, err
	}
	recipeImplementation.UpdateJwtSigningPublicKeyInfo(response["jwtSigningPublicKey"].(string), response["jwtSigningPublicKeyExpiryTime"].(uint64))

	delete(response, "status")
	delete(response, "jwtSigningPublicKey")
	delete(response, "jwtSigningPublicKeyExpiryTime")

	var resp schema.CreateOrRefreshAPIResponse
	bytes := []byte(fmt.Sprintf("%+v", response))
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return schema.CreateOrRefreshAPIResponse{}, err
	}
	return resp, nil
}

func GetSession(recipeImplementation schema.RecipeImplementation, accessToken string, antiCsrfToken *string, doAntiCsrfCheck bool, scontainsCustomHeader bool) {

}

// RevokeSession function used to revoke a specific session
func revokeSession(recipeImplementation schema.RecipeImplementation, sessionHandle string) (bool, error) {
	path, err := supertokens.NewNormalisedURLPath("/recipe/session/remove")
	if err != nil {
		return false, err
	}
	response, err := recipeImplementation.Querier.SendPostRequest(*path,
		map[string]interface{}{
			"sessionHandles": [1]string{sessionHandle},
		})
	if err != nil {
		return false, err
	}
	return len(response["sessionHandlesRevoked"].([]interface{})) == 1, nil
}

func getSessionData(recipeImplementation schema.RecipeImplementation, sessionHandle string) (map[string]interface{}, error) {
	path, err := supertokens.NewNormalisedURLPath("/recipe/session/data")
	if err != nil {
		return nil, err
	}
	response, err := recipeImplementation.Querier.SendGetRequest(*path,
		map[string]interface{}{
			"sessionHandle": sessionHandle,
		})
	if err != nil {
		return nil, err
	}
	if response["status"] == "OK" {
		return response["userDataInDatabase"].(map[string]interface{}), nil
	}
	return nil, errors.UnauthorizedError{
		Msg: response["message"].(string),
	}
}

func updateSessionData(recipeImplementation schema.RecipeImplementation, sessionHandle string, newSessionData map[string]interface{}) error {
	path, err := supertokens.NewNormalisedURLPath("/recipe/session/data")
	if err != nil {
		return err
	}
	response, err := recipeImplementation.Querier.SendPutRequest(*path,
		map[string]interface{}{
			"sessionHandle":      sessionHandle,
			"userDataInDatabase": newSessionData,
		})
	if response["status"] == "UNAUTHORISED" {
		return errors.UnauthorizedError{
			Msg: response["message"].(string),
		}
	}
	return nil
}
