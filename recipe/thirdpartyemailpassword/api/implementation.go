package api

import (
	epapi "github.com/supertokens/supertokens-golang/recipe/emailpassword/api"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	tpapi "github.com/supertokens/supertokens-golang/recipe/thirdparty/api"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
)

func MakeAPIImplementation() tpepmodels.APIInterface {
	emailPasswordImplementation := epapi.MakeAPIImplementation()
	thirdPartyImplementation := tpapi.MakeAPIImplementation()
	return tpepmodels.APIInterface{
		EmailExistsGET: func(email string, options epmodels.APIOptions) (epmodels.EmailExistsGETResponse, error) {
			return emailPasswordImplementation.EmailExistsGET(email, options)

		},
		GeneratePasswordResetTokenPOST: func(formFields []epmodels.TypeFormField, options epmodels.APIOptions) (epmodels.GeneratePasswordResetTokenPOSTResponse, error) {
			return emailPasswordImplementation.GeneratePasswordResetTokenPOST(formFields, options)
		},

		PasswordResetPOST: func(formFields []epmodels.TypeFormField, token string, options epmodels.APIOptions) (epmodels.ResetPasswordUsingTokenResponse, error) {
			return emailPasswordImplementation.PasswordResetPOST(formFields, token, options)
		},

		SignInUpPOST: func(input tpepmodels.SignInUpAPIInput) (tpepmodels.SignInUpAPIOutput, error) {
			if input.EmailpasswordInput != nil {
				if input.EmailpasswordInput.IsSignIn {
					response, err := emailPasswordImplementation.SignInPOST(input.EmailpasswordInput.FormFields, input.EmailpasswordInput.Options)
					if err != nil {
						return tpepmodels.SignInUpAPIOutput{}, err
					}
					if response.OK != nil {
						return tpepmodels.SignInUpAPIOutput{
							EmailpasswordOutput: &tpepmodels.EmailpasswordOutput{
								OK: &struct {
									User           tpepmodels.User
									CreatedNewUser bool
								}{
									User: tpepmodels.User{
										ID:         response.OK.User.ID,
										Email:      response.OK.User.Email,
										TimeJoined: response.OK.User.TimeJoined,
										ThirdParty: nil,
									},
									CreatedNewUser: false,
								},
							},
						}, nil
					} else {
						return tpepmodels.SignInUpAPIOutput{
							EmailpasswordOutput: &tpepmodels.EmailpasswordOutput{
								WrongCredentialsError: &struct{}{},
							},
						}, nil
					}
				} else {
					response, err := emailPasswordImplementation.SignUpPOST(input.EmailpasswordInput.FormFields, input.EmailpasswordInput.Options)
					if err != nil {
						return tpepmodels.SignInUpAPIOutput{}, err
					}
					if response.OK != nil {
						return tpepmodels.SignInUpAPIOutput{
							EmailpasswordOutput: &tpepmodels.EmailpasswordOutput{
								OK: &struct {
									User           tpepmodels.User
									CreatedNewUser bool
								}{
									User: tpepmodels.User{
										ID:         response.OK.User.ID,
										Email:      response.OK.User.Email,
										TimeJoined: response.OK.User.TimeJoined,
										ThirdParty: nil,
									},
									CreatedNewUser: true,
								},
							},
						}, nil
					} else {
						return tpepmodels.SignInUpAPIOutput{
							EmailpasswordOutput: &tpepmodels.EmailpasswordOutput{
								EmailAlreadyExistsError: &struct{}{},
							},
						}, nil
					}
				}
			} else {
				response, err := thirdPartyImplementation.SignInUpPOST(input.ThirdPartyInput.Provider, input.ThirdPartyInput.Code, input.ThirdPartyInput.RedirectURI, input.ThirdPartyInput.Options)
				if err != nil {
					return tpepmodels.SignInUpAPIOutput{}, err
				}
				if response.FieldError != nil {
					return tpepmodels.SignInUpAPIOutput{
						ThirdPartyOutput: &tpepmodels.ThirdPartyOutput{
							FieldError: &struct{ Error string }{},
						},
					}, nil
				} else if response.NoEmailGivenByProviderError != nil {
					return tpepmodels.SignInUpAPIOutput{
						ThirdPartyOutput: &tpepmodels.ThirdPartyOutput{
							NoEmailGivenByProviderError: &struct{}{},
						},
					}, nil
				} else {
					return tpepmodels.SignInUpAPIOutput{
						ThirdPartyOutput: &tpepmodels.ThirdPartyOutput{
							OK: &struct {
								CreatedNewUser   bool
								User             tpepmodels.User
								AuthCodeResponse interface{}
							}{
								CreatedNewUser:   response.OK.CreatedNewUser,
								AuthCodeResponse: response.OK.AuthCodeResponse,
								User: tpepmodels.User{
									ID:         response.OK.User.ID,
									TimeJoined: response.OK.User.TimeJoined,
									Email:      response.OK.User.Email,
									ThirdParty: &response.OK.User.ThirdParty,
								},
							},
						},
					}, nil
				}
			}
		},

		AuthorisationUrlGET: func(provider tpmodels.TypeProvider, options tpmodels.APIOptions) (tpmodels.AuthorisationUrlGETResponse, error) {
			return thirdPartyImplementation.AuthorisationUrlGET(provider, options)
		},
	}
}
