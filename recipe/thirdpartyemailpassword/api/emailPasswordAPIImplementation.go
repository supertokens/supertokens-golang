package api

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
)

func GetEmailPasswordIterfaceImpl(apiImplmentation tpepmodels.APIInterface) epmodels.APIInterface {
	signInUpPOST := apiImplmentation.SignInUpPOST
	if signInUpPOST == nil {
		return epmodels.APIInterface{
			EmailExistsGET:                 apiImplmentation.EmailExistsGET,
			GeneratePasswordResetTokenPOST: apiImplmentation.GeneratePasswordResetTokenPOST,
			PasswordResetPOST:              apiImplmentation.PasswordResetPOST,
			SignInPOST:                     nil,
			SignUpPOST:                     nil,
		}
	}
	return epmodels.APIInterface{
		EmailExistsGET:                 apiImplmentation.EmailExistsGET,
		GeneratePasswordResetTokenPOST: apiImplmentation.GeneratePasswordResetTokenPOST,
		PasswordResetPOST:              apiImplmentation.PasswordResetPOST,
		SignInPOST: func(formFields []epmodels.TypeFormField, options epmodels.APIOptions) (epmodels.SignInResponse, error) {
			resp, err := signInUpPOST(tpepmodels.SignInUpAPIInput{
				EmailpasswordInput: &tpepmodels.EmailpasswordInput{
					FormFields: formFields,
					Options:    options,
					IsSignIn:   true,
				},
			})
			if err != nil {
				return epmodels.SignInResponse{}, err
			}
			result := resp.EmailpasswordOutput
			if result != nil {
				if result.OK != nil {
					return epmodels.SignInResponse{
						OK: &struct{ User epmodels.User }{
							User: epmodels.User{
								ID:         result.OK.User.ID,
								Email:      result.OK.User.Email,
								TimeJoined: result.OK.User.TimeJoined,
							},
						},
					}, nil
				} else if result.WrongCredentialsError != nil {
					return epmodels.SignInResponse{
						WrongCredentialsError: &struct{}{},
					}, nil
				}
			}
			return epmodels.SignInResponse{}, errors.New("should never come here")
		},
		SignUpPOST: func(formFields []epmodels.TypeFormField, options epmodels.APIOptions) (epmodels.SignUpResponse, error) {
			resp, err := signInUpPOST(tpepmodels.SignInUpAPIInput{
				EmailpasswordInput: &tpepmodels.EmailpasswordInput{
					FormFields: formFields,
					Options:    options,
					IsSignIn:   false,
				},
			})
			if err != nil {
				return epmodels.SignUpResponse{}, err
			}
			result := resp.EmailpasswordOutput
			if result != nil {
				if result.OK != nil {
					return epmodels.SignUpResponse{
						OK: &struct{ User epmodels.User }{
							User: epmodels.User{
								ID:         result.OK.User.ID,
								Email:      result.OK.User.Email,
								TimeJoined: result.OK.User.TimeJoined,
							},
						},
					}, nil
				} else if result.EmailAlreadyExistsError != nil {
					return epmodels.SignUpResponse{
						EmailAlreadyExistsError: &struct{}{},
					}, nil
				}
			}
			return epmodels.SignUpResponse{}, errors.New("should never come here")
		},
	}
}
