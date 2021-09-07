package api

import (
	epapi "github.com/supertokens/supertokens-golang/recipe/emailpassword/api"
	epm "github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	tpapi "github.com/supertokens/supertokens-golang/recipe/thirdparty/api"
	tpm "github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/models"
)

func MakeAPIImplementation() models.APIInterface {
	emailPasswordImplementation := epapi.MakeAPIImplementation()
	thirdPartyImplementation := tpapi.MakeAPIImplementation()
	return models.APIInterface{
		EmailExistsGET: func(email string, options epm.APIOptions) (epm.EmailExistsGETResponse, error) {
			return emailPasswordImplementation.EmailExistsGET(email, options)

		},
		GeneratePasswordResetTokenPOST: func(formFields []epm.TypeFormField, options epm.APIOptions) (epm.GeneratePasswordResetTokenPOSTResponse, error) {
			return emailPasswordImplementation.GeneratePasswordResetTokenPOST(formFields, options)
		},

		PasswordResetPOST: func(formFields []epm.TypeFormField, token string, options epm.APIOptions) (epm.ResetPasswordUsingTokenResponse, error) {
			return emailPasswordImplementation.PasswordResetPOST(formFields, token, options)
		},

		SignInUpPOST: func(input models.SignInUpAPIInput) (models.SignInUpAPIOutput, error) {
			if input.EmailpasswordInput != nil {
				if input.EmailpasswordInput.IsSignIn {
					response, err := emailPasswordImplementation.SignInPOST(input.EmailpasswordInput.FormFields, input.EmailpasswordInput.Options)
					if err != nil {
						return models.SignInUpAPIOutput{}, err
					}
					if response.OK != nil {
						return models.SignInUpAPIOutput{
							EmailpasswordOutput: &models.EmailpasswordOutput{
								OK: &struct {
									User           models.User
									CreatedNewUser bool
								}{
									User: models.User{
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
						return models.SignInUpAPIOutput{
							EmailpasswordOutput: &models.EmailpasswordOutput{
								WrongCredentialsError: &struct{}{},
							},
						}, nil
					}
				} else {
					response, err := emailPasswordImplementation.SignUpPOST(input.EmailpasswordInput.FormFields, input.EmailpasswordInput.Options)
					if err != nil {
						return models.SignInUpAPIOutput{}, err
					}
					if response.OK != nil {
						return models.SignInUpAPIOutput{
							EmailpasswordOutput: &models.EmailpasswordOutput{
								OK: &struct {
									User           models.User
									CreatedNewUser bool
								}{
									User: models.User{
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
						return models.SignInUpAPIOutput{
							EmailpasswordOutput: &models.EmailpasswordOutput{
								EmailAlreadyExistsError: &struct{}{},
							},
						}, nil
					}
				}
			} else {
				response, err := thirdPartyImplementation.SignInUpPOST(input.ThirdPartyInput.Provider, input.ThirdPartyInput.Code, input.ThirdPartyInput.RedirectURI, input.ThirdPartyInput.Options)
				if err != nil {
					return models.SignInUpAPIOutput{}, err
				}
				if response.FieldError != nil {
					return models.SignInUpAPIOutput{
						ThirdPartyOutput: &models.ThirdPartyOutput{
							FieldError: &struct{ Error string }{},
						},
					}, nil
				} else if response.NoEmailGivenByProviderError != nil {
					return models.SignInUpAPIOutput{
						ThirdPartyOutput: &models.ThirdPartyOutput{
							NoEmailGivenByProviderError: &struct{}{},
						},
					}, nil
				} else {
					return models.SignInUpAPIOutput{
						ThirdPartyOutput: &models.ThirdPartyOutput{
							OK: &struct {
								CreatedNewUser   bool
								User             models.User
								AuthCodeResponse interface{}
							}{
								CreatedNewUser:   response.OK.CreatedNewUser,
								AuthCodeResponse: response.OK.AuthCodeResponse,
								User: models.User{
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

		AuthorisationUrlGET: func(provider tpm.TypeProvider, options tpm.APIOptions) (tpm.AuthorisationUrlGETResponse, error) {
			return thirdPartyImplementation.AuthorisationUrlGET(provider, options)
		},
	}
}
