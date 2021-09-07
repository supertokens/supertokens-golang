package recipeimplementation

import (
	tpm "github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/models"
)

func MakeThirdPartyRecipeImplementation(recipeImplementation models.RecipeInterface) tpm.RecipeInterface {
	return tpm.RecipeInterface{

		GetUserByThirdPartyInfo: func(thirdPartyID string, thirdPartyUserID string) (*tpm.User, error) {
			user, err := recipeImplementation.GetUserByThirdPartyInfo(thirdPartyID, thirdPartyUserID)
			if err != nil {
				return nil, err
			}
			if user == nil || user.ThirdParty == nil {
				return nil, nil
			}
			return &tpm.User{
				ID:         user.ID,
				Email:      user.Email,
				TimeJoined: user.TimeJoined,
				ThirdParty: *user.ThirdParty,
			}, nil
		},

		SignInUp: func(thirdPartyID string, thirdPartyUserID string, email tpm.EmailStruct) (tpm.SignInUpResponse, error) {
			result, err := recipeImplementation.SignInUp(thirdPartyID, thirdPartyUserID, models.EmailStruct{
				ID:         email.ID,
				IsVerified: email.IsVerified,
			})
			if err != nil {
				return tpm.SignInUpResponse{}, err
			}
			if result.FieldError != nil {
				return tpm.SignInUpResponse{
					FieldError: &struct{ Error string }{
						Error: result.FieldError.Error,
					},
				}, nil
			}

			return tpm.SignInUpResponse{
				OK: &struct {
					CreatedNewUser bool
					User           tpm.User
				}{
					CreatedNewUser: result.OK.CreatedNewUser,
					User: tpm.User{
						ID:         result.OK.User.ID,
						Email:      result.OK.User.Email,
						TimeJoined: result.OK.User.TimeJoined,
						ThirdParty: *result.OK.User.ThirdParty,
					},
				},
			}, nil
		},

		GetUserByID: func(userID string) (*tpm.User, error) {
			user, err := recipeImplementation.GetUserByID(userID)
			if err != nil {
				return nil, err
			}
			if user == nil || user.ThirdParty == nil {
				return nil, nil
			}
			return &tpm.User{
				ID:         user.ID,
				Email:      user.Email,
				TimeJoined: user.TimeJoined,
				ThirdParty: *user.ThirdParty,
			}, nil
		},
	}
}
