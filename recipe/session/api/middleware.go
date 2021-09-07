package api

import (
	"context"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/session/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func VerifySession(recipeInstance models.SessionRecipe, options *models.VerifySessionOptions, otherHandler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := recipeInstance.APIImpl.VerifySession(options, models.APIOptions{
			Config:               recipeInstance.Config,
			OtherHandler:         otherHandler,
			Req:                  r,
			Res:                  w,
			RecipeID:             recipeInstance.RecipeModule.GetRecipeID(),
			RecipeImplementation: recipeInstance.RecipeImpl,
		})
		if err != nil {
			supertokens.ErrorHandler(err, r, w)
			return
		}
		if session != nil {
			ctx := context.WithValue(r.Context(), models.SessionContext, session)
			otherHandler(w, r.WithContext(ctx))
		} else {
			otherHandler(w, r)
		}
	})
}
