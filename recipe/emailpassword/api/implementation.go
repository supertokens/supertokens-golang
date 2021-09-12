package api

import (
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
)

func MakeAPIImplementation() epmodels.APIInterface {
	return epmodels.APIInterface{
		EmailExistsGET: func(email string, options epmodels.APIOptions) (epmodels.EmailExistsGETResponse, error) {
			user, err := options.RecipeImplementation.GetUserByEmail(email)
			if err != nil {
				return epmodels.EmailExistsGETResponse{}, err
			}
			return epmodels.EmailExistsGETResponse{
				OK: &struct{ Exists bool }{Exists: user != nil},
			}, nil
		},

		GeneratePasswordResetTokenPOST: func(formFields []epmodels.TypeFormField, options epmodels.APIOptions) (epmodels.GeneratePasswordResetTokenPOSTResponse, error) {
			var email string
			for _, formField := range formFields {
				if formField.ID == "email" {
					email = formField.Value
				}
			}

			user, err := options.RecipeImplementation.GetUserByEmail(email)
			if err != nil {
				return epmodels.GeneratePasswordResetTokenPOSTResponse{}, err
			}

			if user == nil {
				return epmodels.GeneratePasswordResetTokenPOSTResponse{
					OK: &struct{}{},
				}, nil
			}

			response, err := options.RecipeImplementation.CreateResetPasswordToken(user.ID)
			if err != nil {
				return epmodels.GeneratePasswordResetTokenPOSTResponse{}, err
			}
			if response.UnknownUserIdError != nil {
				return epmodels.GeneratePasswordResetTokenPOSTResponse{
					OK: &struct{}{},
				}, nil
			}

			passwordResetLink := options.Config.ResetPasswordUsingTokenFeature.GetResetPasswordURL(*user) + "?token=" + response.OK.Token + "&rid=" + options.RecipeID

			options.Config.ResetPasswordUsingTokenFeature.CreateAndSendCustomEmail(*user, passwordResetLink)

			return epmodels.GeneratePasswordResetTokenPOSTResponse{
				OK: &struct{}{},
			}, nil
		},

		PasswordResetPOST: func(formFields []epmodels.TypeFormField, token string, options epmodels.APIOptions) (epmodels.ResetPasswordUsingTokenResponse, error) {
			var newPassword string
			for _, formField := range formFields {
				if formField.ID == "password" {
					newPassword = formField.Value
				}
			}

			response, err := options.RecipeImplementation.ResetPasswordUsingToken(token, newPassword)
			if err != nil {
				return epmodels.ResetPasswordUsingTokenResponse{}, err
			}

			return response, nil
		},

		SignInPOST: func(formFields []epmodels.TypeFormField, options epmodels.APIOptions) (epmodels.SignInResponse, error) {
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
				return epmodels.SignInResponse{}, err
			}
			if response.WrongCredentialsError != nil {
				return response, nil
			}

			user := response.OK.User
			_, err = session.CreateNewSession(options.Res, user.ID, map[string]interface{}{}, map[string]interface{}{})
			if err != nil {
				return epmodels.SignInResponse{}, err
			}

			return response, nil
		},

		SignUpPOST: func(formFields []epmodels.TypeFormField, options epmodels.APIOptions) (epmodels.SignUpResponse, error) {
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
				return epmodels.SignUpResponse{}, err
			}
			if response.EmailAlreadyExistsError != nil {
				return response, nil
			}

			user := response.OK.User

			_, err = session.CreateNewSession(options.Res, user.ID, map[string]interface{}{}, map[string]interface{}{})
			if err != nil {
				return epmodels.SignUpResponse{}, err
			}

			return response, nil
		},
	}
}
