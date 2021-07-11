package thirdparty

import "errors"

func getEmailForUserId(userID string) (string, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return "", err
	}
	userInfo := instance.RecipeImpl.GetUserById(userID)
	if userInfo == nil {
		return "", errors.New("Unknown User ID provided")
	}
	return userInfo.Email, nil
}
