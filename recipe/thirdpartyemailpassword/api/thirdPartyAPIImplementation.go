package api

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
)

func GetThirdPartyIterfaceImpl(apiImplmentation tpepmodels.APIInterface) tpmodels.APIInterface {
	signInUpPOST := apiImplmentation.SignInUpPOST
	if signInUpPOST == nil {
		return tpmodels.APIInterface{
			AuthorisationUrlGET: apiImplmentation.AuthorisationUrlGET,
			SignInUpPOST:        nil,
		}
	}
	return tpmodels.APIInterface{

		AuthorisationUrlGET: apiImplmentation.AuthorisationUrlGET,

		SignInUpPOST: func(provider tpmodels.TypeProvider, code, redirectURI string, options tpmodels.APIOptions) (tpmodels.SignInUpPOSTResponse, error) {
			resp, err := signInUpPOST(tpepmodels.SignInUpAPIInput{
				ThirdPartyInput: &tpepmodels.ThirdPartyInput{
					Provider:    provider,
					Code:        code,
					RedirectURI: redirectURI,
					Options:     options,
				},
			})
			if err != nil {
				return tpmodels.SignInUpPOSTResponse{}, err
			}
			result := resp.ThirdPartyOutput
			if result != nil {
				if result.OK != nil {
					return tpmodels.SignInUpPOSTResponse{
						OK: &struct {
							CreatedNewUser   bool
							User             tpmodels.User
							AuthCodeResponse interface{}
						}{
							CreatedNewUser: result.OK.CreatedNewUser,
							User: tpmodels.User{
								ID:         result.OK.User.ID,
								TimeJoined: result.OK.User.TimeJoined,
								Email:      result.OK.User.Email,
								ThirdParty: *result.OK.User.ThirdParty,
							},
						},
					}, nil
				} else if result.NoEmailGivenByProviderError != nil {
					return tpmodels.SignInUpPOSTResponse{
						NoEmailGivenByProviderError: &struct{}{},
					}, nil
				} else if result.FieldError != nil {
					return tpmodels.SignInUpPOSTResponse{
						FieldError: &struct{ Error string }{
							Error: result.FieldError.Error,
						},
					}, nil
				}
			}
			return tpmodels.SignInUpPOSTResponse{}, errors.New("should never come here")
		},
	}
}
