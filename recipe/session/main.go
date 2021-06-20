package session

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/session/models"
)

// functions that can be used by the user...

func CreateNewSession(res http.ResponseWriter, userID string, jwtPayload interface{}, sessionData interface{}) (*models.SessionContainer, error) {
	instance, err := GetInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.CreateNewSession(res, userID, jwtPayload, sessionData)
}

func GetSession(req *http.Request, res http.ResponseWriter, options *models.VerifySessionOptions) (*models.SessionContainer, error) {
	instance, err := GetInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.GetSession(req, res, options)
}

// TODO: Add all the functions - these will be used by the end user
