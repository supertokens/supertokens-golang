package api

import (
	epapi "github.com/supertokens/supertokens-golang/recipe/emailpassword/api"
	epm "github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	tpapi "github.com/supertokens/supertokens-golang/recipe/thirdparty/api"
	tpm "github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/models"
)

func MakeAPIImplementation() models.APIImplementation {
	emailPasswordImplementation := epapi.MakeAPIImplementation()
	thirdPartyImplementation := tpapi.MakeAPIImplementation()
	return models.APIImplementation{
		EmailExistsGET: func(email string, options epm.APIOptions) epm.EmailExistsGETResponse {
			return emailPasswordImplementation.EmailExistsGET(email, options)
		},
		GeneratePasswordResetTokenPOST: func(formFields []epm.TypeFormField, options epm.APIOptions) epm.GeneratePasswordResetTokenPOSTResponse {
			return emailPasswordImplementation.GeneratePasswordResetTokenPOST(formFields, options)
		},
		PasswordResetPOST: func(formFields []epm.TypeFormField, token string, options epm.APIOptions) epm.PasswordResetPOSTResponse {
			return emailPasswordImplementation.PasswordResetPOST(formFields, token, options)
		},
		SignInUpPOST: func(input models.SignInUpAPIInput) models.SignInUpAPIOutput {
			if input.EmailpasswordInput != nil {
				if input.EmailpasswordInput.IsSignIn {
					response := emailPasswordImplementation.SignInPOST(input.EmailpasswordInput.FormFields, input.EmailpasswordInput.Options)
					if response.Status == "OK" {
						return models.SignInUpAPIOutput{
							EmailpasswordOutput: &models.EmailpasswordOutput{
								Status: response.Status,
								User: models.User{
									ID:         response.User.ID,
									Email:      response.User.Email,
									TimeJoined: response.User.TimeJoined,
									ThirdParty: nil,
								},
								CreatedNewUser: false,
							},
						}
					} else {
						return models.SignInUpAPIOutput{
							EmailpasswordOutput: &models.EmailpasswordOutput{
								Status: response.Status,
								CreatedNewUser: false,
							},
						}
					}
				} else {
					response := emailPasswordImplementation.SignUpPOST(input.EmailpasswordInput.FormFields, input.EmailpasswordInput.Options)
					if response.Status == "OK" {
						return models.SignInUpAPIOutput{
							EmailpasswordOutput: &models.EmailpasswordOutput{
								Status: response.Status,
								User: models.User{
									ID:         response.User.ID,
									Email:      response.User.Email,
									TimeJoined: response.User.TimeJoined,
									ThirdParty: nil,
								},
								CreatedNewUser: true,
							},
						}
					} else {
						return models.SignInUpAPIOutput{
							EmailpasswordOutput: &models.EmailpasswordOutput{
								Status: response.Status,
								CreatedNewUser: false,
							},
						}
					}
				}
			}
			response := thirdPartyImplementation.SignInUpPOST(input.ThirdPartyInput.Provider, input.ThirdPartyInput.Code, input.ThirdPartyInput.RedirectURI, input.ThirdPartyInput.Options)
			return models.SignInUpAPIOutput{
				ThirdPartyOutput: &models.ThirdPartyOutput{
					Status:         response.Status,
					CreatedNewUser: response.CreatedNewUser,
					User: models.User{
						ID:         response.User.ID,
						Email:      response.User.Email,
						TimeJoined: response.User.TimeJoined,
						ThirdParty: &response.User.ThirdParty,
					},
					AuthCodeResponse: response.AuthCodeResponse,
					Error:            response.Error,
				},
			}
		},
		AuthorisationUrlGET: func(provider tpm.TypeProvider, options tpm.APIOptions) tpm.AuthorisationUrlGETResponse {
			return thirdPartyImplementation.AuthorisationUrlGET(provider, options)
		},
	}
}
