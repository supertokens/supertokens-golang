package session

import (
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
)

func NewSessionContainerFromSessionContainerInputForTestServer(sessionMap map[string]interface{}) (sessmodels.SessionContainer, error) {
	recipe, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}

	session := &SessionContainerInput{
		sessionHandle:         sessionMap["sessionHandle"].(string),
		userID:                sessionMap["userId"].(string),
		tenantId:              sessionMap["tenantId"].(string),
		userDataInAccessToken: sessionMap["userDataInAccessToken"].(map[string]interface{}),
		accessToken:           sessionMap["accessToken"].(string),
		frontToken:            sessionMap["frontToken"].(string),
		refreshToken:          nil,
		antiCSRFToken:         nil,
		accessTokenUpdated:    sessionMap["accessTokenUpdated"].(bool),
	}

	if refreshToken, ok := sessionMap["refreshToken"].(map[string]interface{}); ok {
		session.refreshToken = &sessmodels.CreateOrRefreshAPIResponseToken{
			Token:       refreshToken["token"].(string),
			Expiry:      uint64(refreshToken["expiry"].(float64)),
			CreatedTime: uint64(refreshToken["createdTime"].(float64)),
		}
	}

	if antiCsrfToken, ok := sessionMap["antiCsrfToken"].(string); ok {
		session.antiCSRFToken = &antiCsrfToken
	}

	session.recipeImpl = recipe.RecipeImpl

	return newSessionContainer(recipe.Config, session), nil
}
