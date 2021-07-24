package recipeimplementation

import (
	epm "github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/models"
)

func MakeEmailPasswordRecipeImplementation(recipeImplementation models.RecipeImplementation) epm.RecipeImplementation {
	return epm.RecipeImplementation{
		SignUp: func(email, password string) epm.SignInUpResponse {
			response := recipeImplementation.SignIn(email, password)
			return epm.SignInUpResponse{
				User: epm.User{
					ID:         response.User.ID,
					Email:      response.User.Email,
					TimeJoined: response.User.TimeJoined,
				},
				Status: response.Status,
			}
		},
		SignIn: func(email, password string) epm.SignInUpResponse {
			response := recipeImplementation.SignIn(email, password)
			return epm.SignInUpResponse{
				User: epm.User{
					ID:         response.User.ID,
					Email:      response.User.Email,
					TimeJoined: response.User.TimeJoined,
				},
				Status: response.Status,
			}
		},
		GetUserByID: func(userId string) *epm.User {
			user := recipeImplementation.GetUserByID(userId)
			if user == nil || user.ThirdParty != nil {
				return nil
			}
			return &epm.User{
				ID:         user.ID,
				Email:      user.Email,
				TimeJoined: user.TimeJoined,
			}
		},
		GetUserByEmail: func(email string) *epm.User {
			user := recipeImplementation.GetUserByEmail(email)
			return &epm.User{
				ID:         user.ID,
				Email:      user.Email,
				TimeJoined: user.TimeJoined,
			}
		},
		CreateResetPasswordToken: func(userID string) epm.CreateResetPasswordTokenResponse {
			return recipeImplementation.CreateResetPasswordToken(userID)
		},
		ResetPasswordUsingToken: func(token, newPassword string) epm.ResetPasswordUsingTokenResponse {
			return recipeImplementation.ResetPasswordUsingToken(token, newPassword)
		},
	}
}
