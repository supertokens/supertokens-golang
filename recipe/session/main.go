package session

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/session/models"
)

func CreateNewSession(res http.ResponseWriter, userID string, jwtPayload interface{}, sessionData interface{}) (*models.SessionContainer, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.CreateNewSession(res, userID, jwtPayload, sessionData)
}

func GetSession(req *http.Request, res http.ResponseWriter, options *models.VerifySessionOptions) (*models.SessionContainer, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.GetSession(req, res, options)
}

func RefreshSession(req *http.Request, res http.ResponseWriter) (*models.SessionContainer, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.RefreshSession(req, res)
}

func RevokeAllSessionsForUser(userID string) ([]string, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.RevokeAllSessionsForUser(userID)
}

func GetAllSessionHandlesForUser(userID string) ([]string, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.GetAllSessionHandlesForUser(userID)
}

func RevokeSession(sessionHandle string) (bool, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}
	return instance.RecipeImpl.RevokeSession(sessionHandle)
}

func RevokeMultipleSessions(sessionHandles []string) ([]string, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.RevokeMultipleSessions(sessionHandles)
}

func GetSessionData(sessionHandle string) (interface{}, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.GetSessionData(sessionHandle)
}

func UpdateSessionData(sessionHandle string, newSessionData interface{}) error {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	return instance.RecipeImpl.UpdateSessionData(sessionHandle, newSessionData)
}

func GetJWTPayload(sessionHandle string) (interface{}, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.GetJWTPayload(sessionHandle)
}

func UpdateJWTPayload(sessionHandle string, newJWTPayload interface{}) error {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	return instance.RecipeImpl.UpdateJWTPayload(sessionHandle, newJWTPayload)
}

func GetAccessTokenLifeTimeMS() (uint64, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return 0, err
	}
	return instance.RecipeImpl.GetAccessTokenLifeTimeMS()
}

func GetRefreshTokenLifeTimeMS() (uint64, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return 0, err
	}
	return instance.RecipeImpl.GetRefreshTokenLifeTimeMS()
}
