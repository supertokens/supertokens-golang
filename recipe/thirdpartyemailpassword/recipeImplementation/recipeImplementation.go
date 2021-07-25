package recipeimplementation

import (
	"reflect"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	epm "github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	tpm "github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeRecipeImplementation(emailPasswordQuerier supertokens.Querier, thirdPartyQuerier *supertokens.Querier) models.RecipeImplementation {
	emailPasswordImplementation := emailpassword.MakeRecipeImplementation(emailPasswordQuerier)
	var thirdPartyImplementation tpm.RecipeImplementation
	if thirdPartyQuerier != nil {
		thirdPartyImplementation = thirdparty.MakeRecipeImplementation(*thirdPartyQuerier)
	}
	return models.RecipeImplementation{
		SignUp: func(email, password string) models.SignInUpResponse {
			response := emailPasswordImplementation.SignUp(email, password)
			return models.SignInUpResponse{
				User: models.User{
					ID:         response.User.ID,
					Email:      response.User.Email,
					TimeJoined: response.User.TimeJoined,
					ThirdParty: nil,
				},
				Status: response.Status,
			}
		},
		SignIn: func(email, password string) models.SignInUpResponse {
			response := emailPasswordImplementation.SignIn(email, password)
			return models.SignInUpResponse{
				User: models.User{
					ID:         response.User.ID,
					Email:      response.User.Email,
					TimeJoined: response.User.TimeJoined,
					ThirdParty: nil,
				},
				Status: response.Status,
			}
		},
		SignInUp: func(thirdPartyID, thirdPartyUserID string, email tpm.EmailStruct) models.SignInUpResponse {
			result := thirdPartyImplementation.SignInUp(thirdPartyID, thirdPartyUserID, email)
			return models.SignInUpResponse{
				Status:         "OK",
				CreatedNewUser: result.CreatedNewUser,
				User: models.User{
					ID:         result.User.ID,
					Email:      result.User.Email,
					TimeJoined: result.User.TimeJoined,
					ThirdParty: &result.User.ThirdParty,
				},
			}
		},
		GetUserByID: func(userID string) *models.User {
			user := emailPasswordImplementation.GetUserByID(userID)
			if user == nil {
				return nil
			}
			if user != nil {
				return &models.User{
					ID:         user.ID,
					Email:      user.Email,
					TimeJoined: user.TimeJoined,
					ThirdParty: nil,
				}
			}
			if reflect.DeepEqual(thirdPartyImplementation, tpm.RecipeImplementation{}) {
				return nil
			}
			userinfo := thirdPartyImplementation.GetUserByID(userID)
			if userinfo == nil {
				return nil
			}
			if userinfo != nil {
				return &models.User{
					ID:         userinfo.ID,
					Email:      userinfo.Email,
					TimeJoined: userinfo.TimeJoined,
					ThirdParty: &userinfo.ThirdParty,
				}
			}
			return nil
		},
		GetUserByThirdPartyInfo: func(thirdPartyID string, thirdPartyUserID string) *models.User {
			if reflect.DeepEqual(thirdPartyImplementation, tpm.RecipeImplementation{}) {
				return nil
			}
			userinfo := thirdPartyImplementation.GetUserByThirdPartyInfo(thirdPartyID, thirdPartyUserID)
			if userinfo == nil {
				return nil
			}
			if userinfo != nil {
				return &models.User{
					ID:         userinfo.ID,
					Email:      userinfo.Email,
					TimeJoined: userinfo.TimeJoined,
					ThirdParty: &userinfo.ThirdParty,
				}
			}
			return nil
		},
		GetUserByEmail: func(email string) *models.User {
			userinfo := emailPasswordImplementation.GetUserByEmail(email)
			if userinfo == nil {
				return nil
			}
			if userinfo != nil {
				return &models.User{
					ID:         userinfo.ID,
					Email:      userinfo.Email,
					TimeJoined: userinfo.TimeJoined,
					ThirdParty: nil,
				}
			}
			return nil
		},
		CreateResetPasswordToken: func(userID string) epm.CreateResetPasswordTokenResponse {
			return emailPasswordImplementation.CreateResetPasswordToken(userID)
		},
		ResetPasswordUsingToken: func(token, newPassword string) epm.ResetPasswordUsingTokenResponse {
			return emailPasswordImplementation.ResetPasswordUsingToken(token, newPassword)
		},
	}
}
