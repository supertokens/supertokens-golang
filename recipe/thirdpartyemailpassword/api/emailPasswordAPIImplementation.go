package api

import (
	"errors"

	epm "github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/models"
)

func GetEmailPasswordIterfaceImpl(apiImplmentation models.APIInterface) epm.APIInterface {
	signInUpPOST := apiImplmentation.SignInUpPOST
	if signInUpPOST == nil {
		return epm.APIInterface{
			EmailExistsGET:                 apiImplmentation.EmailExistsGET,
			GeneratePasswordResetTokenPOST: apiImplmentation.GeneratePasswordResetTokenPOST,
			PasswordResetPOST:              apiImplmentation.PasswordResetPOST,
			SignInPOST:                     nil,
			SignUpPOST:                     nil,
		}
	}
	return epm.APIInterface{
		EmailExistsGET:                 apiImplmentation.EmailExistsGET,
		GeneratePasswordResetTokenPOST: apiImplmentation.GeneratePasswordResetTokenPOST,
		PasswordResetPOST:              apiImplmentation.PasswordResetPOST,
		SignInPOST: func(formFields []epm.TypeFormField, options epm.APIOptions) (epm.SignInResponse, error) {
			resp, err := signInUpPOST(models.SignInUpAPIInput{
				EmailpasswordInput: &models.EmailpasswordInput{
					FormFields: formFields,
					Options:    options,
					IsSignIn:   true,
				},
			})
			if err != nil {
				return epm.SignInResponse{}, err
			}
			result := resp.EmailpasswordOutput
			if result != nil {
				if result.OK != nil {
					return epm.SignInResponse{
						OK: &struct{ User epm.User }{
							User: epm.User{
								ID:         result.OK.User.ID,
								Email:      result.OK.User.Email,
								TimeJoined: result.OK.User.TimeJoined,
							},
						},
					}, nil
				} else if result.WrongCredentialsError != nil {
					return epm.SignInResponse{
						WrongCredentialsError: &struct{}{},
					}, nil
				}
			}
			return epm.SignInResponse{}, errors.New("should never come here")
		},
		SignUpPOST: func(formFields []epm.TypeFormField, options epm.APIOptions) (epm.SignUpResponse, error) {
			resp, err := signInUpPOST(models.SignInUpAPIInput{
				EmailpasswordInput: &models.EmailpasswordInput{
					FormFields: formFields,
					Options:    options,
					IsSignIn:   false,
				},
			})
			if err != nil {
				return epm.SignUpResponse{}, err
			}
			result := resp.EmailpasswordOutput
			if result != nil {
				if result.OK != nil {
					return epm.SignUpResponse{
						OK: &struct{ User epm.User }{
							User: epm.User{
								ID:         result.OK.User.ID,
								Email:      result.OK.User.Email,
								TimeJoined: result.OK.User.TimeJoined,
							},
						},
					}, nil
				} else if result.EmailAlreadyExistsError != nil {
					return epm.SignUpResponse{
						EmailAlreadyExistsError: &struct{}{},
					}, nil
				}
			}
			return epm.SignUpResponse{}, errors.New("should never come here")
		},
	}
}
