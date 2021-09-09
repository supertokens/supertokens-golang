package recipeimplementation

import (
	epm "github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/models"
)

func MakeEmailPasswordRecipeImplementation(recipeImplementation models.RecipeInterface) epm.RecipeInterface {
	return epm.RecipeInterface{
		SignUp: func(email, password string) (epm.SignUpResponse, error) {
			response, err := recipeImplementation.SignUp(email, password)
			if err != nil {
				return epm.SignUpResponse{}, err
			}
			if response.EmailAlreadyExistsError != nil {
				return epm.SignUpResponse{
					EmailAlreadyExistsError: &struct{}{},
				}, nil
			}
			return epm.SignUpResponse{
				OK: &struct{ User epm.User }{
					User: epm.User{
						ID:         response.OK.User.ID,
						Email:      response.OK.User.Email,
						TimeJoined: response.OK.User.TimeJoined,
					},
				},
			}, nil
		},

		SignIn: func(email, password string) (epm.SignInResponse, error) {
			response, err := recipeImplementation.SignIn(email, password)
			if err != nil {
				return epm.SignInResponse{}, err
			}
			if response.WrongCredentialsError != nil {
				return epm.SignInResponse{
					WrongCredentialsError: &struct{}{},
				}, nil
			}
			return epm.SignInResponse{
				OK: &struct{ User epm.User }{
					User: epm.User{
						ID:         response.OK.User.ID,
						Email:      response.OK.User.Email,
						TimeJoined: response.OK.User.TimeJoined,
					},
				},
			}, nil
		},

		GetUserByID: func(userId string) (*epm.User, error) {
			user, err := recipeImplementation.GetUserByID(userId)
			if err != nil {
				return nil, err
			}
			if user == nil || user.ThirdParty != nil {
				return nil, nil
			}
			return &epm.User{
				ID:         user.ID,
				Email:      user.Email,
				TimeJoined: user.TimeJoined,
			}, nil
		},

		GetUserByEmail: func(email string) (*epm.User, error) {
			users, err := recipeImplementation.GetUsersByEmail(email)
			if err != nil {
				return nil, err
			}

			for _, user := range users {
				if user.ThirdParty == nil {
					return &epm.User{
						ID:         user.ID,
						Email:      user.Email,
						TimeJoined: user.TimeJoined,
					}, nil
				}
			}
			return nil, nil
		},

		CreateResetPasswordToken: func(userID string) (epm.CreateResetPasswordTokenResponse, error) {
			return recipeImplementation.CreateResetPasswordToken(userID)
		},
		ResetPasswordUsingToken: func(token, newPassword string) (epm.ResetPasswordUsingTokenResponse, error) {
			return recipeImplementation.ResetPasswordUsingToken(token, newPassword)
		},
		UpdateEmailOrPassword: func(userId string, email, password *string) (epm.UpdateEmailOrPasswordResponse, error) {
			return recipeImplementation.UpdateEmailOrPassword(userId, email, password)
		},
	}
}
