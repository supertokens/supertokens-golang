package api

import (
	epm "github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/models"
)

func GetEmailPasswordIterfaceImpl(apiImplmentation models.APIImplementation) epm.APIImplementation {
	signInUpPOST := apiImplmentation.SignInUpPOST
	if signInUpPOST == nil {
		return epm.APIImplementation{
			EmailExistsGET:                 apiImplmentation.EmailExistsGET,
			GeneratePasswordResetTokenPOST: apiImplmentation.GeneratePasswordResetTokenPOST,
			PasswordResetPOST:              apiImplmentation.PasswordResetPOST,
			SignInPOST:                     nil,
			SignUpPOST:                     nil,
		}
	}
	return epm.APIImplementation{
		EmailExistsGET:                 apiImplmentation.EmailExistsGET,
		GeneratePasswordResetTokenPOST: apiImplmentation.GeneratePasswordResetTokenPOST,
		PasswordResetPOST:              apiImplmentation.PasswordResetPOST,
		SignInPOST: func(formFields []epm.TypeFormField, options epm.APIOptions) epm.SignInUpResponse {
			resp := signInUpPOST(models.SignInUpAPIInput{
				EmailpasswordInput: &models.EmailpasswordInput{
					FormFields: formFields,
					Options:    options,
					IsSignIn:   true,
				},
			})
			result := resp.EmailpasswordOutput
			if result != nil {
				if result.Status == "OK" {
					return epm.SignInUpResponse{
						Status: result.Status,
						User: epm.User{
							ID:         result.User.ID,
							Email:      result.User.Email,
							TimeJoined: result.User.TimeJoined,
						},
					}
				} else if result.Status == "WRONG_CREDENTIALS_ERROR" {
					return epm.SignInUpResponse{
						Status: "WRONG_CREDENTIALS_ERROR",
					}
				}
			}
			return epm.SignInUpResponse{}
		},
		SignUpPOST: func(formFields []epm.TypeFormField, options epm.APIOptions) epm.SignInUpResponse {
			resp := signInUpPOST(models.SignInUpAPIInput{
				EmailpasswordInput: &models.EmailpasswordInput{
					FormFields: formFields,
					Options:    options,
					IsSignIn:   false,
				},
			})
			result := resp.EmailpasswordOutput
			if result != nil {
				if result.Status == "OK" {
					return epm.SignInUpResponse{
						Status: result.Status,
						User: epm.User{
							ID:         result.User.ID,
							Email:      result.User.Email,
							TimeJoined: result.User.TimeJoined,
						},
					}
				} else if result.Status == "EMAIL_ALREADY_EXISTS_ERROR" {
					return epm.SignInUpResponse{
						Status: "EMAIL_ALREADY_EXISTS_ERROR",
					}
				}
			}
			return epm.SignInUpResponse{}
		},
	}
}
