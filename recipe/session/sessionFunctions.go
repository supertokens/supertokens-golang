package session

import (
	"encoding/json"
	"fmt"

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
