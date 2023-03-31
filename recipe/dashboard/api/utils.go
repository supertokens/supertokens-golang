package api

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/passwordless"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless"
	"github.com/supertokens/supertokens-golang/supertokens"
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
			userToReturn.ThirdParty = &dashboardmodels.ThirdParty{
				Id:     response.ThirdParty.ID,
				UserId: response.ThirdParty.UserID,
			}
		}

		if userToReturn == (dashboardmodels.UserType{}) {
			tpepResponse, tpepError := thirdpartyemailpassword.GetUserById(userId)

			if tpepError == nil {
				userToReturn.Id = tpepResponse.ID
				userToReturn.TimeJoined = tpepResponse.TimeJoined
				userToReturn.FirstName = ""
				userToReturn.LastName = ""
				userToReturn.Email = tpepResponse.Email
				userToReturn.ThirdParty = &dashboardmodels.ThirdParty{
					Id:     tpepResponse.ThirdParty.ID,
					UserId: tpepResponse.ThirdParty.UserID,
				}
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
	isRecipeInitialised := false

	if recipeId == emailpassword.RECIPE_ID {
		_, err := emailpassword.GetRecipeInstanceOrThrowError()

		if err == nil {
			isRecipeInitialised = true
		}

		if !isRecipeInitialised {
			_, err := thirdpartyemailpassword.GetRecipeInstanceOrThrowError()

			if err == nil {
				isRecipeInitialised = true
			}
		}
	} else if recipeId == passwordless.RECIPE_ID {
		_, err := passwordless.GetRecipeInstanceOrThrowError()

		if err == nil {
			isRecipeInitialised = true
		}

		if !isRecipeInitialised {
			_, err := thirdpartypasswordless.GetRecipeInstanceOrThrowError()

			if err == nil {
				isRecipeInitialised = true
			}
		}
	} else if recipeId == thirdparty.RECIPE_ID {
		_, err := thirdparty.GetRecipeInstanceOrThrowError()

		if err == nil {
			isRecipeInitialised = true
		}

		if !isRecipeInitialised {
			_, err := thirdpartyemailpassword.GetRecipeInstanceOrThrowError()

			if err == nil {
				isRecipeInitialised = true
			}
		}

		if !isRecipeInitialised {
			_, err := thirdpartypasswordless.GetRecipeInstanceOrThrowError()

			if err == nil {
				isRecipeInitialised = true
			}
		}
	}

	return isRecipeInitialised
}

// TODO: Add tests
func getUsersWithSearch(timeJoinedOrder string, paginationToken *string, limit *int, includeRecipeIds *[]string, searchParams map[string]string) (supertokens.UserPaginationResult, error) {

	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		return supertokens.UserPaginationResult{}, err
	}

	requestBody := map[string]string{}
	if searchParams != nil {
		requestBody = searchParams
	}
	requestBody["timeJoinedOrder"] = timeJoinedOrder
	if limit != nil {
		requestBody["limit"] = strconv.Itoa(*limit)
	}
	if paginationToken != nil {
		requestBody["paginationToken"] = *paginationToken
	}
	if includeRecipeIds != nil {
		requestBody["includeRecipeIds"] = strings.Join((*includeRecipeIds)[:], ",")
	}

	resp, err := querier.SendGetRequest("/users", requestBody)

	if err != nil {
		return supertokens.UserPaginationResult{}, err
	}

	temporaryVariable, err := json.Marshal(resp)
	if err != nil {
		return supertokens.UserPaginationResult{}, err
	}

	var result = supertokens.UserPaginationResult{}

	err = json.Unmarshal(temporaryVariable, &result)

	if err != nil {
		return supertokens.UserPaginationResult{}, err
	}

	return result, nil
}
