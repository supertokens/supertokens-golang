package emailpassword

import "github.com/supertokens/supertokens-golang/recipe/emailpassword/models"

func SignUp(email string, password string) (models.User, error) {
	_, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return models.User{}, err
	}
	// return instance.RecipeImpl.SignUp(email, password)
	return models.User{}, nil
}
