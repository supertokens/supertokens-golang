package session

import (
	"context"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/session/api"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Init(config *sessmodels.TypeInput) supertokens.Recipe {
	return recipeInit(config)
}

func CreateNewSession(res http.ResponseWriter, userID string, jwtPayload map[string]interface{}, sessionData map[string]interface{}) (sessmodels.SessionContainer, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return sessmodels.SessionContainer{}, err
	}
	return instance.RecipeImpl.CreateNewSession(res, userID, jwtPayload, sessionData)
}

func GetSession(req *http.Request, res http.ResponseWriter, options *sessmodels.VerifySessionOptions) (*sessmodels.SessionContainer, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.GetSession(req, res, options)
}

func GetSessionInformation(sessionHandle string) (sessmodels.SessionInformation, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return sessmodels.SessionInformation{}, err
	}
	return instance.RecipeImpl.GetSessionInformation(sessionHandle)
}

func RefreshSession(req *http.Request, res http.ResponseWriter) (sessmodels.SessionContainer, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return sessmodels.SessionContainer{}, err
	}
	return instance.RecipeImpl.RefreshSession(req, res)
}

func RevokeAllSessionsForUser(userID string) ([]string, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.RevokeAllSessionsForUser(userID)
}

func GetAllSessionHandlesForUser(userID string) ([]string, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.GetAllSessionHandlesForUser(userID)
}

func RevokeSession(sessionHandle string) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}
	return instance.RecipeImpl.RevokeSession(sessionHandle)
}

func RevokeMultipleSessions(sessionHandles []string) ([]string, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return instance.RecipeImpl.RevokeMultipleSessions(sessionHandles)
}

func UpdateSessionData(sessionHandle string, newSessionData map[string]interface{}) error {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	return instance.RecipeImpl.UpdateSessionData(sessionHandle, newSessionData)
}

func UpdateJWTPayload(sessionHandle string, newJWTPayload map[string]interface{}) error {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	return instance.RecipeImpl.UpdateJWTPayload(sessionHandle, newJWTPayload)
}

func VerifySession(options *sessmodels.VerifySessionOptions, otherHandler http.HandlerFunc) http.HandlerFunc {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		panic("can't fetch supertokens instance. You should call the supertokens.Init function before using the VerifySession function.")
	}
	return api.VerifySession(*instance, options, otherHandler)
}

func GetSessionFromRequestContext(ctx context.Context) *sessmodels.SessionContainer {
	value := ctx.Value(sessmodels.SessionContext)
	if value == nil {
		return nil
	}
	temp := value.(*sessmodels.SessionContainer)
	return temp
}
