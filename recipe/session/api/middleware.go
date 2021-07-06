package api

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/session/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func VerifySession(recipeInstance models.SessionRecipe, options *models.VerifySessionOptions) func(w http.ResponseWriter, r *http.Request, otherHandler http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, otherHandler http.HandlerFunc) {
		recipeInstance.APIImpl.VerifySession(options, models.APIOptions{
			Config:               recipeInstance.Config,
			OtherHandler:         otherHandler,
			Req:                  r,
			Res:                  w,
			RecipeID:             recipeInstance.RecipeModule.GetRecipeID(),
			RecipeImplementation: recipeInstance.RecipeImpl,
		})
	}
}

func SendTryRefreshTokenResponse(recipeInstance models.SessionRecipe, _ string, _ *http.Request, response http.ResponseWriter) error {
	return supertokens.SendNon200Response(response, "try refresh token", recipeInstance.Config.SessionExpiredStatusCode)
}

func SendUnauthorisedResponse(recipeInstance models.SessionRecipe, _ string, _ *http.Request, response http.ResponseWriter) error {
	return supertokens.SendNon200Response(response, "unauthorised", recipeInstance.Config.SessionExpiredStatusCode)
}

func SendTokenTheftDetectedResponse(recipeInstance models.SessionRecipe, sessionHandle string, _ string, _ *http.Request, response http.ResponseWriter) error {
	_, err := recipeInstance.RecipeImpl.RevokeSession(sessionHandle)
	if err != nil {
		return err
	}
	return supertokens.SendNon200Response(response, "token theft detected", recipeInstance.Config.SessionExpiredStatusCode)
}
