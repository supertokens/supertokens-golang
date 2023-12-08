package supertokens

import "errors"

type RecipeUserID struct {
	recipeUserID string
}

func NewRecipeUserID(recipeUserID string) (*RecipeUserID, error) {
	if recipeUserID == "" {
		return nil, errors.New("recipeUserID cannot be empty")
	}
	return &RecipeUserID{recipeUserID: recipeUserID}, nil
}

func (r *RecipeUserID) GetAsString() string {
	return r.recipeUserID
}
