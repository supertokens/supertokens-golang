package recipeimplementation

import (
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
)

func MakeThirdPartyRecipeImplementation(recipeImplementation tpepmodels.RecipeInterface) tpmodels.RecipeInterface {
	return tpmodels.RecipeInterface{

		GetUserByThirdPartyInfo: func(thirdPartyID string, thirdPartyUserID string) (*tpmodels.User, error) {
			user, err := recipeImplementation.GetUserByThirdPartyInfo(thirdPartyID, thirdPartyUserID)
			if err != nil {
				return nil, err
			}
			if user == nil || user.ThirdParty == nil {
				return nil, nil
			}
			return &tpmodels.User{
				ID:         user.ID,
				Email:      user.Email,
				TimeJoined: user.TimeJoined,
				ThirdParty: *user.ThirdParty,
			}, nil
		},

		SignInUp: func(thirdPartyID string, thirdPartyUserID string, email tpmodels.EmailStruct) (tpmodels.SignInUpResponse, error) {
			result, err := recipeImplementation.SignInUp(thirdPartyID, thirdPartyUserID, tpepmodels.EmailStruct{
				ID:         email.ID,
				IsVerified: email.IsVerified,
			})
			if err != nil {
				return tpmodels.SignInUpResponse{}, err
			}
			if result.FieldError != nil {
				return tpmodels.SignInUpResponse{
					FieldError: &struct{ Error string }{
						Error: result.FieldError.Error,
					},
				}, nil
			}

			return tpmodels.SignInUpResponse{
				OK: &struct {
					CreatedNewUser bool
					User           tpmodels.User
				}{
					CreatedNewUser: result.OK.CreatedNewUser,
					User: tpmodels.User{
						ID:         result.OK.User.ID,
						Email:      result.OK.User.Email,
						TimeJoined: result.OK.User.TimeJoined,
						ThirdParty: *result.OK.User.ThirdParty,
					},
				},
			}, nil
		},

		GetUserByID: func(userID string) (*tpmodels.User, error) {
			user, err := recipeImplementation.GetUserByID(userID)
			if err != nil {
				return nil, err
			}
			if user == nil || user.ThirdParty == nil {
				return nil, nil
			}
			return &tpmodels.User{
				ID:         user.ID,
				Email:      user.Email,
				TimeJoined: user.TimeJoined,
				ThirdParty: *user.ThirdParty,
			}, nil
		},
	}
}
