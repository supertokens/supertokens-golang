package api

import (
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/recipe/session"
)

func MakeAPIImplementation() models.APIInterface {
	return models.APIInterface{
		EmailExistsGET: func(email string, options models.APIOptions) (models.EmailExistsGETResponse, error) {
			user, err := options.RecipeImplementation.GetUserByEmail(email)
			if err != nil {
				return models.EmailExistsGETResponse{}, err
			}
			return models.EmailExistsGETResponse{
				OK: &struct{ Exists bool }{Exists: user != nil},
			}, nil
		},

		GeneratePasswordResetTokenPOST: func(formFields []models.TypeFormField, options models.APIOptions) (models.GeneratePasswordResetTokenPOST, error) {
			var email string
			for _, formField := range formFields {
				if formField.ID == "email" {
					email = formField.Value
				}
			}

			user, err := options.RecipeImplementation.GetUserByEmail(email)
			if err != nil {
				return models.GeneratePasswordResetTokenPOST{}, err
			}

			if user == nil {
				return models.GeneratePasswordResetTokenPOST{
					OK: &struct{}{},
				}, nil
			}

			response, err := options.RecipeImplementation.CreateResetPasswordToken(user.ID)
			if err != nil {
				return models.GeneratePasswordResetTokenPOST{}, err
			}
			if response.UnknownUserIdError != nil {
				return models.GeneratePasswordResetTokenPOST{
					OK: &struct{}{},
				}, nil
			}

			passwordResetLink := options.Config.ResetPasswordUsingTokenFeature.GetResetPasswordURL(*user) + "?token=" + response.OK.Token + "&rid=" + options.RecipeID

			options.Config.ResetPasswordUsingTokenFeature.CreateAndSendCustomEmail(*user, passwordResetLink)

			return models.GeneratePasswordResetTokenPOST{
				OK: &struct{}{},
			}, nil
		},

		PasswordResetPOST: func(formFields []models.TypeFormField, token string, options models.APIOptions) (models.ResetPasswordUsingTokenResponse, error) {
			var newPassword string
			for _, formField := range formFields {
				if formField.ID == "password" {
					newPassword = formField.Value
				}
			}

			response, err := options.RecipeImplementation.ResetPasswordUsingToken(token, newPassword)
			if err != nil {
				return models.ResetPasswordUsingTokenResponse{}, err
			}

			return response, nil
		},

		SignInPOST: func(formFields []models.TypeFormField, options models.APIOptions) (models.SignInResponse, error) {
			var email string
			for _, formField := range formFields {
				if formField.ID == "email" {
					email = formField.Value
				}
			}
			var password string
			for _, formField := range formFields {
				if formField.ID == "password" {
					password = formField.Value
				}
			}

			response, err := options.RecipeImplementation.SignIn(email, password)
			if err != nil {
				return models.SignInResponse{}, err
			}
			if response.WrongCredentialsError != nil {
				return response, nil
			}

			user := response.OK.User
			_, err = session.CreateNewSession(options.Res, user.ID, map[string]interface{}{}, map[string]interface{}{})
			if err != nil {
				return models.SignInResponse{}, err
			}

			return response, nil
		},

		SignUpPOST: func(formFields []models.TypeFormField, options models.APIOptions) (models.SignUpResponse, error) {
			var email string
			for _, formField := range formFields {
				if formField.ID == "email" {
					email = formField.Value
				}
			}
			var password string
			for _, formField := range formFields {
				if formField.ID == "password" {
					password = formField.Value
				}
			}

			response, err := options.RecipeImplementation.SignUp(email, password)
			if err != nil {
				return models.SignUpResponse{}, err
			}
			if response.EmailAlreadyExistsError != nil {
				return response, nil
			}

			user := response.OK.User

			_, err = session.CreateNewSession(options.Res, user.ID, map[string]interface{}{}, map[string]interface{}{})
			if err != nil {
				return models.SignUpResponse{}, err
			}

			return response, nil
		},
	}
}
