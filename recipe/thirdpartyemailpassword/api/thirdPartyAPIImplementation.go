package api

import (
	tpm "github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/models"
)

func GetThirdPartyIterfaceImpl(apiImplmentation models.APIImplementation) tpm.APIImplementation {
	signInUpPOST := apiImplmentation.SignInUpPOST
	if signInUpPOST == nil {
		return tpm.APIImplementation{
			AuthorisationUrlGET: apiImplmentation.AuthorisationUrlGET,
			SignInUpPOST:        nil,
		}
	}
	return tpm.APIImplementation{
		AuthorisationUrlGET: apiImplmentation.AuthorisationUrlGET,
		SignInUpPOST: func(provider tpm.TypeProvider, code, redirectURI string, options tpm.APIOptions) tpm.SignInUpPOSTResponse {
			resp := signInUpPOST(models.SignInUpAPIInput{
				ThirdPartyInput: &models.ThirdPartyInput{
					Provider:    provider,
					Code:        code,
					RedirectURI: redirectURI,
					Options:     options,
				},
			})
			result := resp.ThirdPartyOutput
			if result != nil {
				if result.Status == "OK" {
					return tpm.SignInUpPOSTResponse{
						Status:         result.Status,
						CreatedNewUser: result.CreatedNewUser,
						User: tpm.User{
							ID:         result.User.ID,
							Email:      result.User.Email,
							TimeJoined: result.User.TimeJoined,
							ThirdParty: *result.User.ThirdParty,
						},
						AuthCodeResponse: result.AuthCodeResponse,
					}
				} else if result.Status == "NO_EMAIL_GIVEN_BY_PROVIDER" {
					return tpm.SignInUpPOSTResponse{
						Status: "NO_EMAIL_GIVEN_BY_PROVIDER",
					}
				} else if result.Status == "FIELD_ERROR" {
					return tpm.SignInUpPOSTResponse{
						Status: "FIELD_ERROR",
						Error:  result.Error,
					}
				}
			}
			return tpm.SignInUpPOSTResponse{}
		},
	}
}
