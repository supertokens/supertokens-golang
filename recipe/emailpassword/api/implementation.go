package api

import "github.com/supertokens/supertokens-golang/recipe/emailpassword/models"

// TODO:
func MakeAPIImplementation() models.APIImplementation {
	return models.APIImplementation{
		EmailExistsGET: func(email string, options models.APIOptions) models.EmailExistsGETResponse {
			return models.EmailExistsGETResponse{}
		},

		GeneratePasswordResetTokenPOST: func(formFields []models.FormFieldValue, options models.APIOptions) models.GeneratePasswordResetTokenPOSTResponse {
			return models.GeneratePasswordResetTokenPOSTResponse{}
		},

		PasswordResetPOST: func(formFields []models.FormFieldValue, token string, options models.APIOptions) models.PasswordResetPOSTResponse {
			return models.PasswordResetPOSTResponse{}
		},

		SignInPOST: func(formFields []models.FormFieldValue, options models.APIOptions) models.SignInUpResponse {
			return models.SignInUpResponse{}
		},

		SignUpPOST: func(formFields []models.FormFieldValue, options models.APIOptions) models.SignInUpResponse {
			return models.SignInUpResponse{}
		},
	}
}
