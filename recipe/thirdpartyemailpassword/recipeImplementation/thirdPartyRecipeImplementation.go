package recipeimplementation

import (
	"errors"

	tpm "github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/models"
)

func MakeThirdPartyRecipeImplementation(recipeImplementation models.RecipeImplementation) tpm.RecipeImplementation {
	return tpm.RecipeImplementation{
		GetUserByID: func(userID string) *tpm.User {
			user := recipeImplementation.GetUserByID(userID)
			if user == nil || user.ThirdParty == nil {
				return nil
			}
			return &tpm.User{
				ID:         user.ID,
				Email:      user.Email,
				TimeJoined: user.TimeJoined,
				ThirdParty: *user.ThirdParty,
			}
		},
		GetUserByThirdPartyInfo: func(thirdPartyID string, thirdPartyUserID string) *tpm.User {
			user := recipeImplementation.GetUserByThirdPartyInfo(thirdPartyID, thirdPartyUserID)
			if user == nil || user.ThirdParty == nil {
				return nil
			}
			return &tpm.User{
				ID:         user.ID,
				Email:      user.Email,
				TimeJoined: user.TimeJoined,
				ThirdParty: *user.ThirdParty,
			}
		},
		SignInUp: func(thirdPartyID string, thirdPartyUserID string, email tpm.EmailStruct) tpm.SignInUpResponse {
			result := recipeImplementation.SignInUp(thirdPartyID, thirdPartyUserID, email)
			if result.Status == "FIELD_ERROR" {
				return tpm.SignInUpResponse{
					Status:         result.Status,
					CreatedNewUser: result.CreatedNewUser,
					User: tpm.User{
						ID:         result.User.ID,
						Email:      result.User.Email,
						TimeJoined: result.User.TimeJoined,
						ThirdParty: *result.User.ThirdParty,
					},
				}
			}
			if result.User.ThirdParty == nil {
				return tpm.SignInUpResponse{
					Status: "FIELD_ERROR",
					Error:  errors.New("Should never come here"),
				}
			}
			return tpm.SignInUpResponse{
				Status:         "OK",
				CreatedNewUser: result.CreatedNewUser,
				User: tpm.User{
					ID:         result.User.ID,
					Email:      result.User.Email,
					TimeJoined: result.User.TimeJoined,
					ThirdParty: *result.User.ThirdParty,
				},
			}
		},
	}
}
