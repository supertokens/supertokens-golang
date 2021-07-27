package models

import (
	epm "github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	evm "github.com/supertokens/supertokens-golang/recipe/emailverification/models"
	tpm "github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
)

type User struct {
	ID         string
	TimeJoined uint64
	Email      string
	ThirdParty *struct {
		ID     string
		UserID string
	}
}

type TypeContext struct {
	FormFields                 []epm.TypeFormField
	ThirdPartyAuthCodeResponse interface{}
}

type TypeInputSetJwtPayloadForSession func(user User, context TypeContext, action string) map[string]interface{}
type TypeInputSetSessionDataForSession func(user User, context TypeContext, action string) map[string]interface{}

type TypeNormalisedInputSessionFeature struct {
	SetJwtPayload  TypeInputSetJwtPayloadForSession
	SetSessionData TypeInputSetSessionDataForSession
}

type TypeInputSignUp struct {
	FormFields []epm.TypeInputFormField
}

type TypeNormalisedInputSignUp struct {
	FormFields []epm.NormalisedFormField
}

type TypeInputEmailVerificationFeature struct {
	GetEmailVerificationURL  func(user User) (string, error)
	CreateAndSendCustomEmail func(user User, emailVerificationURLWithToken string) error
}

type TypeInput struct {
	SessionFeature                 *TypeNormalisedInputSessionFeature
	SignUpFeature                  *TypeInputSignUp
	Providers                      []tpm.TypeProvider
	ResetPasswordUsingTokenFeature *epm.TypeInputResetPasswordUsingTokenFeature
	EmailVerificationFeature       *TypeInputEmailVerificationFeature
	Override                       *struct {
		Functions                func(originalImplementation RecipeImplementation) RecipeImplementation
		APIs                     func(originalImplementation APIImplementation) APIImplementation
		EmailVerificationFeature *struct {
			Functions func(originalImplementation evm.RecipeImplementation) evm.RecipeImplementation
			APIs      func(originalImplementation evm.APIImplementation) evm.APIImplementation
		}
	}
}

type TypeNormalisedInput struct {
	SessionFeature                 TypeNormalisedInputSessionFeature
	SignUpFeature                  TypeNormalisedInputSignUp
	Providers                      []tpm.TypeProvider
	ResetPasswordUsingTokenFeature *epm.TypeInputResetPasswordUsingTokenFeature
	EmailVerificationFeature       evm.TypeInput
	Override                       struct {
		Functions                func(originalImplementation RecipeImplementation) RecipeImplementation
		APIs                     func(originalImplementation APIImplementation) APIImplementation
		EmailVerificationFeature *struct {
			Functions func(originalImplementation evm.RecipeImplementation) evm.RecipeImplementation
			APIs      func(originalImplementation evm.APIImplementation) evm.APIImplementation
		}
	}
}

type RecipeImplementation struct {
	GetUserByID             func(userID string) *User
	GetUserByThirdPartyInfo func(thirdPartyID string, thirdPartyUserID string) *User
	SignInUp                func(thirdPartyID string, thirdPartyUserID string, email tpm.EmailStruct) SignInUpResponse

	SignUp                   func(email string, password string) SignInUpResponse
	SignIn                   func(email string, password string) SignInUpResponse
	GetUserByEmail           func(email string) *User
	CreateResetPasswordToken func(userID string) epm.CreateResetPasswordTokenResponse
	ResetPasswordUsingToken  func(token string, newPassword string) epm.ResetPasswordUsingTokenResponse
}

type SignInUpResponse struct {
	Status         string `json:"status"`
	CreatedNewUser bool   `json:"createdNewUser"`
	User           User   `json:"user"`
	Error          error  `json:"error"`
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
	Provider    tpm.TypeProvider
	Code        string
	RedirectURI string
	Options     tpm.APIOptions
}

type SignInUpAPIOutput struct {
	EmailpasswordOutput *EmailpasswordOutput
	ThirdPartyOutput    *ThirdPartyOutput
}

type EmailpasswordOutput struct {
	Status         string `json:"status"`
	User           User   `json:"user"`
	CreatedNewUser bool   `json:"createdNewUser"`
}

type ThirdPartyOutput struct {
	Status           string      `json:"status"`
	CreatedNewUser   bool        `json:"createdNewUser"`
	User             User        `json:"user"`
	AuthCodeResponse interface{} `json:"authCodeResponse"`
	Error            error       `json:"error"`
}

type APIImplementation struct {
	AuthorisationUrlGET func(provider tpm.TypeProvider, options tpm.APIOptions) tpm.AuthorisationUrlGETResponse

	EmailExistsGET                 func(email string, options epm.APIOptions) epm.EmailExistsGETResponse
	GeneratePasswordResetTokenPOST func(formFields []epm.TypeFormField, options epm.APIOptions) epm.GeneratePasswordResetTokenPOSTResponse
	PasswordResetPOST              func(formFields []epm.TypeFormField, token string, options epm.APIOptions) epm.PasswordResetPOSTResponse

	SignInUpPOST func(input SignInUpAPIInput) SignInUpAPIOutput
}
