package emailpassword

import (
	defaultErrors "errors"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const RECIPE_ID = "emailpassword"

type Recipe struct {
	RecipeModule supertokens.RecipeModule
	Config       models.TypeNormalisedInput
	RecipeImpl   models.RecipeImplementation
	APIImpl      models.APIImplementation
}

var r *Recipe = nil

func GetRecipeInstanceOrThrowError() (*Recipe, error) {
	if r != nil {
		return r, nil
	}
	return nil, defaultErrors.New("Initialisation not done. Did you forget to call the init function?")
}
