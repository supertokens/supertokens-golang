package api

import (
	"errors"

	tpm "github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/models"
)

func GetThirdPartyIterfaceImpl(apiImplmentation models.APIInterface) tpm.APIInterface {
	signInUpPOST := apiImplmentation.SignInUpPOST
	if signInUpPOST == nil {
		return tpm.APIInterface{
			AuthorisationUrlGET: apiImplmentation.AuthorisationUrlGET,
			SignInUpPOST:        nil,
		}
	}
	return tpm.APIInterface{

		AuthorisationUrlGET: apiImplmentation.AuthorisationUrlGET,

		SignInUpPOST: func(provider tpm.TypeProvider, code, redirectURI string, options tpm.APIOptions) (tpm.SignInUpPOSTResponse, error) {
			resp, err := signInUpPOST(models.SignInUpAPIInput{
				ThirdPartyInput: &models.ThirdPartyInput{
					Provider:    provider,
					Code:        code,
					RedirectURI: redirectURI,
					Options:     options,
				},
			})
			if err != nil {
				return tpm.SignInUpPOSTResponse{}, err
			}
			result := resp.ThirdPartyOutput
			if result != nil {
				if result.OK != nil {
					return tpm.SignInUpPOSTResponse{
						OK: &struct {
							CreatedNewUser   bool
							User             tpm.User
							AuthCodeResponse interface{}
						}{
							CreatedNewUser: result.OK.CreatedNewUser,
							User: tpm.User{
								ID:         result.OK.User.ID,
								TimeJoined: result.OK.User.TimeJoined,
								Email:      result.OK.User.Email,
								ThirdParty: *result.OK.User.ThirdParty,
							},
						},
					}, nil
				} else if result.NoEmailGivenByProviderError != nil {
					return tpm.SignInUpPOSTResponse{
						NoEmailGivenByProviderError: &struct{}{},
					}, nil
				} else if result.FieldError != nil {
					return tpm.SignInUpPOSTResponse{
						FieldError: &struct{ Error string }{
							Error: result.FieldError.Error,
						},
					}, nil
				}
			}
			return tpm.SignInUpPOSTResponse{}, errors.New("should never come here")
		},
	}
}
