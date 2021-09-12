package models

import (
	epm "github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
)

type APIInterface struct {
	AuthorisationUrlGET            func(provider tpmodels.TypeProvider, options tpmodels.APIOptions) (tpmodels.AuthorisationUrlGETResponse, error)
	EmailExistsGET                 func(email string, options epm.APIOptions) (epm.EmailExistsGETResponse, error)
	GeneratePasswordResetTokenPOST func(formFields []epm.TypeFormField, options epm.APIOptions) (epm.GeneratePasswordResetTokenPOSTResponse, error)
	PasswordResetPOST              func(formFields []epm.TypeFormField, token string, options epm.APIOptions) (epm.ResetPasswordUsingTokenResponse, error)
	SignInUpPOST                   func(input SignInUpAPIInput) (SignInUpAPIOutput, error)
}

type SignInUpAPIInput struct {
	EmailpasswordInput *EmailpasswordInput
	ThirdPartyInput    *ThirdPartyInput
}

type EmailpasswordInput struct {
	IsSignIn   bool
	FormFields []epm.TypeFormField
	Options    epm.APIOptions
}

type ThirdPartyInput struct {
	Provider    tpmodels.TypeProvider
	Code        string
	RedirectURI string
	Options     tpmodels.APIOptions
}

type SignInUpAPIOutput struct {
	EmailpasswordOutput *EmailpasswordOutput
	ThirdPartyOutput    *ThirdPartyOutput
}

type EmailpasswordOutput struct {
	OK *struct {
		User           User
		CreatedNewUser bool
	}
	EmailAlreadyExistsError *struct{}
	WrongCredentialsError   *struct{}
}

type ThirdPartyOutput struct {
	OK *struct {
		CreatedNewUser   bool
		User             User
		AuthCodeResponse interface{}
	}
	NoEmailGivenByProviderError *struct{}
	FieldError                  *struct{ Error string }
}
