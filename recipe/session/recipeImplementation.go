package session

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/session/schema"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeRecipeImplementation(querier supertokens.Querier, config schema.TypeNormalisedInput) schema.RecipeImplementation {
	return schema.RecipeImplementation{
		CreateNewSession: func(res http.ResponseWriter, userID string, jwtPayload interface{}, sessionData interface{}) schema.SessionContainerInterface {
			return schema.SessionContainerInterface{}
		},
		GetSession: func(req *http.Request, res http.ResponseWriter, options *schema.VerifySessionOptions) *schema.SessionContainerInterface {
			return nil
		},
	}
}
