package tpepmodels

import (
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
)

type APIInterface struct {
	AuthorisationUrlGET            func(provider tpmodels.TypeProvider, options tpmodels.APIOptions) (tpmodels.AuthorisationUrlGETResponse, error)
	EmailExistsGET                 func(email string, options epmodels.APIOptions) (epmodels.EmailExistsGETResponse, error)
	GeneratePasswordResetTokenPOST func(formFields []epmodels.TypeFormField, options epmodels.APIOptions) (epmodels.GeneratePasswordResetTokenPOSTResponse, error)
	PasswordResetPOST              func(formFields []epmodels.TypeFormField, token string, options epmodels.APIOptions) (epmodels.ResetPasswordUsingTokenResponse, error)
	SignInUpPOST                   func(input SignInUpAPIInput) (SignInUpAPIOutput, error)
}

type SignInUpAPIInput struct {
	EmailpasswordInput *EmailpasswordInput
	ThirdPartyInput    *ThirdPartyInput
}

type EmailpasswordInput struct {
	IsSignIn   bool
	FormFields []epmodels.TypeFormField
	Options    epmodels.APIOptions
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
