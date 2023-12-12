package supertokens

import (
	"encoding/json"
	"errors"
)

type RecipeUserID struct {
	recipeUserID string
}

func NewRecipeUserID(recipeUserID string) (RecipeUserID, error) {
	if recipeUserID == "" {
		return RecipeUserID{}, errors.New("recipeUserID cannot be empty")
	}
	return RecipeUserID{recipeUserID: recipeUserID}, nil
}

func (r *RecipeUserID) GetAsString() string {
	return r.recipeUserID
}

func (r *RecipeUserID) MarshalJSON() ([]byte, error) {
	// convert r.recipeUserId to string and return that
	return json.Marshal(r.recipeUserID)
}
