package emailverification

import "github.com/supertokens/supertokens-golang/supertokens"

const RECIPE_ID = "emailverification"

type Recipe struct {
	instance            *Recipe
	RecipeID            string
	config              TypeNormalisedInput
	recipeInterfaceImpl RecipeInterface
	apiImpl             APIInterface
}

func NewRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, isInServerlessEnv bool, config TypeInput) *Recipe {
	return &Recipe{
		// config: config,
		// recipeInterfaceImpl: config.override.functions(NewRecipeImplementation()),
	}
}
