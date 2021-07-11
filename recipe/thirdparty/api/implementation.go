package api

import "github.com/supertokens/supertokens-golang/recipe/thirdparty/models"

// TODO:
func MakeAPIImplementation() models.APIImplementation {
	return models.APIImplementation{
		AuthorisationUrlGET: func(provider models.TypeProvider, options models.APIOptions) models.AuthorisationUrlGETResponse {
			return models.AuthorisationUrlGETResponse{}
		},
		SignInUpPOST: func(provider models.TypeProvider, code, redirectURI string, options models.APIOptions) models.SignInUpPOSTResponse {
			return models.SignInUpPOSTResponse{}
		},
	}
}
