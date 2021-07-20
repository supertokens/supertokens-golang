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

type TypeContextEmailPasswordSignUp struct {
	LoginType  string
	FormFields []epm.TypeFormField
}

type TypeContextEmailPasswordSessionDataAndJWT struct {
	LoginType  string
	FormFields []epm.TypeFormField
}

type TypeContextEmailPasswordSignIn struct {
	LoginType string
}

type TypeContextThirdParty struct {
	LoginType                  string
	ThirdPartyAuthCodeResponse interface{}
}

type TypeInputSetJwtPayloadForSession func(user User, context interface{}, action string) map[string]interface{}
type TypeInputSetSessionDataForSession func(user User, context interface{}, action string) map[string]interface{}

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
	ResetPasswordUsingTokenFeature epm.TypeInputResetPasswordUsingTokenFeature
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
	SignInUp                func(thirdPartyID string, thirdPartyUserID string, email tpm.EmailStruct) tpm.SignInUpResponse

	SignUp                   func(email string, password string) epm.SignInUpResponse
	SignIn                   func(email string, password string) epm.SignInUpResponse
	GetUserByEmail           func(email string) *User
	CreateResetPasswordToken func(userId string) epm.CreateResetPasswordTokenResponse
	ResetPasswordUsingToken  func(token string, newPassword string) epm.ResetPasswordUsingTokenResponse
}

type EmailPasswordAPIOptions epm.APIOptions

type ThirdPartyAPIOptions tpm.APIOptions

// type SignInAPIInput struct {}
// type SignInUpAPIOutput struct{}

type APIImplementation struct {
	AuthorisationUrlGET func(provider tpm.TypeProvider, options tpm.APIOptions) tpm.AuthorisationUrlGETResponse
	SignInUpPOST        func(provider tpm.TypeProvider, code string, redirectURI string, options tpm.APIOptions) tpm.SignInUpPOSTResponse

	EmailExistsGET                 func(email string, options epm.APIOptions) epm.EmailExistsGETResponse
	GeneratePasswordResetTokenPOST func(formFields []epm.TypeFormField, options epm.APIOptions) epm.GeneratePasswordResetTokenPOSTResponse
	PasswordResetPOST              func(formFields []epm.TypeFormField, token string, options epm.APIOptions) epm.PasswordResetPOSTResponse
	SignInPOST                     func(formFields []epm.TypeFormField, options epm.APIOptions) epm.SignInUpResponse
	SignUpPOST                     func(formFields []epm.TypeFormField, options epm.APIOptions) epm.SignInUpResponse
}
