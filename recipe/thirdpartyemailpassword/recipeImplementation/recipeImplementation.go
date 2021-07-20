package recipeimplementation

import (
	"reflect"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	tpm "github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeRecipeImplementation(emailPasswordQuerier supertokens.Querier, thirdPartyQuerier *supertokens.Querier) models.RecipeImplementation {
	emailPasswordImplementation := emailpassword.MakeRecipeImplementation(emailPasswordQuerier)
	var thirdPartyImplementation tpm.RecipeImplementation
	if thirdPartyQuerier != nil {
		thirdPartyImplementation = thirdparty.MakeRecipeImplementation(*thirdPartyQuerier)
	}
	return models.RecipeImplementation{
		GetUserByID: func(userID string) *models.User {
			user := emailPasswordImplementation.GetUserByID(userID)
			if user != nil {
				return &models.User{}
			}
			if reflect.DeepEqual(thirdPartyImplementation, tpm.RecipeImplementation{}) {
				return nil
			}
			userinfo := thirdPartyImplementation.GetUserByID(userID)
			if userinfo != nil {
				return &models.User{}
			}
			return nil
		},
	}
}
