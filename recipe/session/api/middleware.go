package api

import (
	"context"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/session/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

// TODO: we need to expose this?
const sessionContext contextKey = iota

func VerifySession(recipeInstance models.SessionRecipe, options *models.VerifySessionOptions) func(w http.ResponseWriter, r *http.Request, otherHandler http.HandlerFunc) error {
	return func(w http.ResponseWriter, r *http.Request, otherHandler http.HandlerFunc) error {
		session, err := recipeInstance.APIImpl.VerifySession(options, models.APIOptions{
			Config:               recipeInstance.Config,
			OtherHandler:         otherHandler,
			Req:                  r,
			Res:                  w,
			RecipeID:             recipeInstance.RecipeModule.GetRecipeID(),
			RecipeImplementation: recipeInstance.RecipeImpl,
		})
		if err != nil {
			return supertokens.ErrorHandler(err, r, w)
		}
		if session != nil {
			var ctx context.Context
			ctx = context.WithValue(r.Context(), sessionContext, session)
			otherHandler(w, r.WithContext(ctx))
		} else {
			otherHandler(w, r)
		}
		return nil
	}
}
