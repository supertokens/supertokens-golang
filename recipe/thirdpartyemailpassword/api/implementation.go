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
			return epm.EmailExistsGETResponse{}
		},
		GeneratePasswordResetTokenPOST: func(formFields []epm.TypeFormField, options epm.APIOptions) epm.GeneratePasswordResetTokenPOSTResponse {
			return epm.GeneratePasswordResetTokenPOSTResponse{}
		},
		PasswordResetPOST: func(formFields []epm.TypeFormField, token string, options epm.APIOptions) epm.PasswordResetPOSTResponse {
			return epm.PasswordResetPOSTResponse{}
		},
		SignInPOST: func(formFields []epm.TypeFormField, options epm.APIOptions) epm.SignInUpResponse {
			return epm.SignInUpResponse{}
		},
		SignUpPOST: func(formFields []epm.TypeFormField, options epm.APIOptions) epm.SignInUpResponse {
			return epm.SignInUpResponse{}
		},
		AuthorisationUrlGET: func(provider tpm.TypeProvider, options tpm.APIOptions) tpm.AuthorisationUrlGETResponse {
			return tpm.AuthorisationUrlGETResponse{}
		},
		SignInUpPOST: func(provider tpm.TypeProvider, code, redirectURI string, options tpm.APIOptions) tpm.SignInUpPOSTResponse {
			return tpm.SignInUpPOSTResponse{}
		},
	}
}
