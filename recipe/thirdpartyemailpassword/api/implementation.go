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
		SignInPOST: func(formFields []epm.TypeFormField, options epm.APIOptions) epm.SignInUpResponse {
			return emailPasswordImplementation.SignInPOST(formFields, options)
		},
		SignUpPOST: func(formFields []epm.TypeFormField, options epm.APIOptions) epm.SignInUpResponse {
			return emailPasswordImplementation.SignInPOST(formFields, options)
		},
		AuthorisationUrlGET: func(provider tpm.TypeProvider, options tpm.APIOptions) tpm.AuthorisationUrlGETResponse {
			return thirdPartyImplementation.AuthorisationUrlGET(provider, options)
		},
		SignInUpPOST: func(provider tpm.TypeProvider, code, redirectURI string, options tpm.APIOptions) tpm.SignInUpPOSTResponse {
			return thirdPartyImplementation.SignInUpPOST(provider, code, redirectURI, options)
		},
	}
}
