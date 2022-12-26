package api

import (
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/passwordless"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless"
)

func IsValidRecipeId(recipeId string) bool {
	return recipeId == "emailpassword" || recipeId == "thirdparty" || recipeId == "passwordless"
}

/*
This function tries to fetch a user for the given user id and recipe id. The input recipe id
should be one of the primary recipes (emailpassword, thirdparty, passwordless) but the returned
recipe will be the exact recipe that matched for the user (including thirdpartyemailpassword and
thirdpartypasswordless).

When fetching a user we need to check for multiple recipes per input recipe id, for example a user
created using email and password could be present for the EmailPassword recipe and the ThirdPartyEmailPassword
recipe so we need to check for both.

If this function returns an empty user struct, it should be treated as if the user does not exist
*/
func GetUserForRecipeId(userId string, recipeId string) (user dashboardmodels.UserType, recipe string) {
	var userToReturn dashboardmodels.UserType
	var recipeToReturn string

	if recipeId == emailpassword.RECIPE_ID {
		response, error := emailpassword.GetUserByID(userId)

		if error == nil {
			userToReturn.Id = response.ID
			userToReturn.TimeJoined = response.TimeJoined
			userToReturn.FirstName = ""
			userToReturn.LastName = ""
			userToReturn.Email = response.Email

			recipeToReturn = emailpassword.RECIPE_ID
		}

		if userToReturn == (dashboardmodels.UserType{}) {
			tpepResponse, tpepError := thirdpartyemailpassword.GetUserById(userId)

			if tpepError == nil {
				userToReturn.Id = tpepResponse.ID
				userToReturn.TimeJoined = tpepResponse.TimeJoined
				userToReturn.FirstName = ""
				userToReturn.LastName = ""
				userToReturn.Email = tpepResponse.Email

				recipeToReturn = thirdpartyemailpassword.RECIPE_ID
			}
		}
	} else if recipeId == thirdparty.RECIPE_ID {
		response, error := thirdparty.GetUserByID(userId)

		if error == nil {
			userToReturn.Id = response.ID
			userToReturn.TimeJoined = response.TimeJoined
			userToReturn.FirstName = ""
			userToReturn.LastName = ""
			userToReturn.Email = response.Email
			userToReturn.ThirdParty.Id = response.ThirdParty.ID
			userToReturn.ThirdParty.UserId = response.ThirdParty.UserID
		}

		if userToReturn == (dashboardmodels.UserType{}) {
			tpepResponse, tpepError := thirdpartyemailpassword.GetUserById(userId)

			if tpepError == nil {
				userToReturn.Id = tpepResponse.ID
				userToReturn.TimeJoined = tpepResponse.TimeJoined
				userToReturn.FirstName = ""
				userToReturn.LastName = ""
				userToReturn.Email = tpepResponse.Email
				userToReturn.ThirdParty.Id = tpepResponse.ThirdParty.ID
				userToReturn.ThirdParty.UserId = tpepResponse.ThirdParty.UserID
			}
		}
	} else if recipeId == passwordless.RECIPE_ID {
		response, error := passwordless.GetUserByID(userId)

		if error == nil {
			userToReturn.Id = response.ID
			userToReturn.TimeJoined = response.TimeJoined
			userToReturn.FirstName = ""
			userToReturn.LastName = ""

			if response.Email != nil {
				userToReturn.Email = *response.Email
			}

			if response.PhoneNumber != nil {
				userToReturn.Phone = *response.PhoneNumber
			}
		}

		if userToReturn == (dashboardmodels.UserType{}) {
			tppResponse, tppError := thirdpartypasswordless.GetUserByID(userId)

			if tppError == nil {
				userToReturn.Id = tppResponse.ID
				userToReturn.TimeJoined = tppResponse.TimeJoined
				userToReturn.FirstName = ""
				userToReturn.LastName = ""

				if tppResponse.Email != nil {
					userToReturn.Email = *tppResponse.Email
				}

				if tppResponse.PhoneNumber != nil {
					userToReturn.Phone = *tppResponse.PhoneNumber
				}
			}
		}
	}

	return userToReturn, recipeToReturn
}

func IsRecipeInitialised(recipeId string) bool {
	if recipeId == emailpassword.RECIPE_ID {
		instance := emailpassword.GetRecipeInstance()

		return instance != nil
	} else if recipeId == passwordless.RECIPE_ID {
		instance := passwordless.GetRecipeInstance()

		return instance != nil
	} else if recipeId == thirdparty.RECIPE_ID {
		_, err := thirdparty.GetRecipeInstanceOrThrowError()

		return err == nil
	}

	return false
}
